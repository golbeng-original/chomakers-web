package models

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type PotofolioModel struct {
	Id     int64        `json:"id"`
	Title  string       `json:"title"`
	Images []ImageModel `json:"images"`
}

type PotofolioRepository struct {
	DBConnect *DBConnection
	ImageRepo *ImageRepository
}

func (repo *PotofolioRepository) CreateTable() error {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	createPotofolioTableQuery := `
		CREATE TABLE IF NOT EXISTS "potofolio" 
		(
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"title" TEXT
		)
		`
	_, err = db.Exec(createPotofolioTableQuery)
	if err != nil {
		log.Printf("[error] create table potofolio [%v]\n", err)
		return err
	}

	return nil
}

func (repo *PotofolioRepository) GetPotofolioList() ([]PotofolioModel, error) {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	potofolioRows, err := db.Query("SELECT id, title FROM potofolio")
	if err != nil {
		log.Printf("[error] potofolio query [%v]\n", err)
		return nil, err
	}

	defer potofolioRows.Close()

	potofolios := make([]PotofolioModel, 0)
	for potofolioRows.Next() {

		var potofolioId int64
		var title string
		err := potofolioRows.Scan(&potofolioId, &title)
		if err != nil {
			log.Printf("[error] potofolio scan [%v]\n", err)
			continue
		}

		potofolio := PotofolioModel{Id: potofolioId, Title: title}
		potofolio.Images, err = repo.ImageRepo.GetImages(PotofolioType, potofolioId)
		if err != nil {
			log.Printf("[error] potofolio Image Query [%v]\n", err)
		}

		potofolios = append(potofolios, potofolio)

	}

	return potofolios, nil
}

func (repo *PotofolioRepository) FindPotofolio(potofolioId int64) (*PotofolioModel, error) {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	potofolioRows, err := db.Query("SELECT title FROM potofolio WHERE id = $1", potofolioId)
	if err != nil {
		log.Printf("[error] potofolio query [%v]\n", err)
		return nil, err
	}

	defer potofolioRows.Close()

	if !potofolioRows.Next() {
		return nil, fmt.Errorf("potofolio not found [id: %v]", potofolioId)
	}

	var title string
	err = potofolioRows.Scan(&title)
	if err != nil {
		return nil, err
	}

	potofolioModel := PotofolioModel{Id: potofolioId, Title: title}

	potofolioModel.Images, err = repo.ImageRepo.GetImages(PotofolioType, potofolioId)

	if err != nil {
		log.Printf("[error] potofolio Image Query [%v]\n", err)
	}

	return &potofolioModel, nil
}

func (repo *PotofolioRepository) AddPotofolio(title string, images []string) (int64, error) {

	completed := false

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return 0, err
	}

	transaction, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer CloseTranstion(transaction, &completed)

	potofolioInsertQuery := "INSERT INTO potofolio (title) VALUES ($1)"

	result, err := transaction.Exec(potofolioInsertQuery, title)
	if err != nil {
		return 0, err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = repo.ImageRepo.AddImgesTransaction(transaction, PotofolioType, insertId, images)
	if err != nil {
		return 0, err
	}

	completed = true

	return insertId, nil
}

func (repo *PotofolioRepository) UpdatePotofolio(potofolioId int64, title *string, removeImageIds []int64, addImages []string) ([]string, error) {

	completed := false
	removeImagePaths := make([]string, 0)

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	transaction, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer CloseTranstion(transaction, &completed)

	if title != nil {
		potofolioUpdateQuery := fmt.Sprintf("UPDATE potofolio SET title = \"%v\" WHERE id = $1", *title)

		_, err = transaction.Exec(potofolioUpdateQuery, potofolioId)
		if err != nil {
			return nil, err
		}
	}

	// 지워질 이미지 찾기
	if removeImageIds != nil {

		images, err := repo.ImageRepo.GetImages(PotofolioType, potofolioId)
		if err != nil {
			return nil, err
		}

		for _, image := range images {

			exists := false
			for _, removeId := range removeImageIds {
				if removeId == image.Id {
					exists = true
					break
				}
			}

			if !exists {
				continue
			}

			err = repo.ImageRepo.RemoveImageIdTransaction(transaction, PotofolioType, potofolioId, image.Id)
			if err != nil {
				return nil, err
			}

			removeImagePaths = append(removeImagePaths, image.Path)
		}
	}

	if addImages != nil {

		err = repo.ImageRepo.AddImgesTransaction(transaction, PotofolioType, potofolioId, addImages)
		if err != nil {
			return nil, err
		}
	}

	if removeImageIds != nil || addImages != nil {

		err = repo.ImageRepo.SortImageOrderTransation(transaction, PotofolioType, potofolioId)

		if err != nil {
			return nil, err
		}

	}

	completed = true

	return removeImagePaths, nil
}

func (repo *PotofolioRepository) RemovePotofolio(potofolioId int64) error {

	completed := false

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	defer CloseTranstion(transaction, &completed)

	potofolioDeleteQuery := "DELETE FROM potofolio WHERE id = $1"
	_, err = transaction.Exec(potofolioDeleteQuery, potofolioId)
	if err != nil {
		return err
	}

	err = repo.ImageRepo.RemoveImagesTransaction(transaction, PotofolioType, potofolioId)
	if err != nil {
		return err
	}

	completed = true

	return nil
}
