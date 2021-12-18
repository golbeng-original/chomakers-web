package models

import (
	"fmt"
	"strings"
)

type EssayThumnailModel struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	ThumbnailImage string `json:"thumbnail"`
}

type EssayModel struct {
	Id             int64        `json:"id"`
	Title          string       `json:"title"`
	ThumbnailImage string       `json:"thumbnail"`
	Images         []ImageModel `json:"images"`
	EssayContent   string       `json:"essayContent"`
}

func (essayModel *EssayModel) GetThumbnailImagePath(saveDir string, prefixUri string) string {
	return strings.Replace(essayModel.ThumbnailImage, prefixUri, saveDir, 1)
}

type EssayRepository struct {
	DBConnect *DBConnection
	ImageRepo *ImageRepository
}

func (repo *EssayRepository) CreateTable() error {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	createEssayTableQuery := `
		CREATE TABLE IF NOT EXISTS "essay"
		(
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"title" TEXT,
			"thumbImage" TEXT,
			"essayContent" TEXT
		)`

	_, err = db.Exec(createEssayTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func (repo *EssayRepository) GetEssayList() ([]EssayThumnailModel, error) {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	essayRows, err := db.Query("SELECT id, title, thumbImage FROM essay")
	if err != nil {
		return nil, err
	}

	defer essayRows.Close()

	essaies := make([]EssayThumnailModel, 0)

	for essayRows.Next() {
		var id int64
		var title string
		var thumbnailImage string
		err = essayRows.Scan(&id, &title, &thumbnailImage)
		if err != nil {
			return nil, err
		}

		essay := EssayThumnailModel{Id: id, Title: title, ThumbnailImage: thumbnailImage}
		essaies = append(essaies, essay)

	}

	return essaies, nil
}

func (repo *EssayRepository) FindEssay(essayId int64) (*EssayModel, error) {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	findQuery := "SELECT id, title, thumbImage, essayContent FROM essay WHERE id = $1"
	essayRow, err := db.Query(findQuery, essayId)
	if err != nil {
		return nil, err
	}

	defer essayRow.Close()
	if !essayRow.Next() {
		return nil, fmt.Errorf("essay not found [id:%v]", essayId)
	}

	var id int64
	var title string
	var thumbnail string
	var essayContent string
	err = essayRow.Scan(&id, &title, &thumbnail, &essayContent)
	if err != nil {
		return nil, err
	}

	essayModel := EssayModel{Id: id, Title: title, ThumbnailImage: thumbnail, EssayContent: essayContent}
	images, err := repo.ImageRepo.GetImages(EssayType, id)
	if err != nil {
		return nil, err
	}

	essayModel.Images = images

	return &essayModel, nil
}

func (repo *EssayRepository) AddEssay(title string, thumbnailPath string, essayContent string, images []string) (int64, error) {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return 0, err
	}

	transaction, err := db.Begin()
	if err != nil {
		return 0, err
	}

	completed := false
	defer CloseTranstion(transaction, &completed)

	insertQuery := "INSERT INTO essay (title, thumbImage, essayContent) VALUES ($1, $2, $3)"
	insertResult, err := transaction.Exec(insertQuery, title, thumbnailPath, essayContent)
	if err != nil {
		return 0, err
	}

	insertId, err := insertResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = repo.ImageRepo.AddImgesTransaction(transaction, EssayType, insertId, images)
	if err != nil {
		return 0, err
	}

	completed = true

	return insertId, nil
}

func (repo *EssayRepository) UpdateEssay(essayId int64, title *string, thumbnailPath *string, essayContent *string, removeImageIds []int64, addIamge []string) ([]string, error) {

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

	updateColumns := ""
	if title != nil {
		updateColumns = fmt.Sprintf("title = \"%v\"", *title)
	}

	if thumbnailPath != nil {
		if len(updateColumns) > 0 {
			updateColumns += ", "
		}
		updateColumns += fmt.Sprintf("thumbImage = \"%v\"", *thumbnailPath)
	}

	if essayContent != nil {
		if len(updateColumns) > 0 {
			updateColumns += ", "
		}

		updateColumns += fmt.Sprintf("essayContent = \"%v\"", *essayContent)
	}

	if len(updateColumns) > 0 {
		updateColumns = fmt.Sprintf("UPDATE essay SET %v WHERE id = $1", updateColumns)

		_, err = db.Exec(updateColumns, essayId)
		if err != nil {
			return nil, err
		}
	}

	if removeImageIds != nil {
		essayImages, err := repo.ImageRepo.GetImages(EssayType, essayId)
		if err != nil {
			return nil, err
		}

		for _, essayImage := range essayImages {
			exists := false
			for _, removeId := range removeImageIds {
				if removeId == essayImage.Id {
					exists = true
					break
				}
			}

			if !exists {
				continue
			}

			err = repo.ImageRepo.RemoveImageIdTransaction(transaction, EssayType, essayId, essayImage.Id)
			if err != nil {
				return nil, err
			}

			removeImagePaths = append(removeImagePaths, essayImage.Path)
		}
	}

	if addIamge != nil {
		err = repo.ImageRepo.AddImgesTransaction(transaction, EssayType, essayId, addIamge)
		if err != nil {
			return nil, err
		}
	}

	if removeImageIds != nil || addIamge != nil {
		err = repo.ImageRepo.SortImageOrderTransation(transaction, EssayType, essayId)

		if err != nil {
			return nil, err
		}
	}

	completed = true

	return removeImagePaths, nil
}

func (repo *EssayRepository) RemoveEssay(essayId int64) error {

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

	essayDeleteQuery := "DELETE FROM essay WHERE id = $1"
	_, err = transaction.Exec(essayDeleteQuery, essayId)
	if err != nil {
		return err
	}

	err = repo.ImageRepo.RemoveImagesTransaction(transaction, EssayType, essayId)
	if err != nil {
		return err
	}

	completed = true

	return nil
}
