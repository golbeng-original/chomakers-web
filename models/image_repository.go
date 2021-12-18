package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
)

type ImageModel struct {
	Id   int64  `json:"id"`
	Path string `json:"path"`
}

type ImageRepository struct {
	DBConnect *DBConnection
}

func (repo *ImageRepository) CreateTable() error {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	createImageTableQuery := `
		CREATE TABLE IF NOT EXISTS "images"
		(
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"dependencyId" INTEGER,
			"dependencyType" INTEGER,
			"imagePath" TEXT,
			"imageOrder" INTEGER
		)
	`

	_, err = db.Exec(createImageTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ImageRepository) GetImages(dependencyType RepositoryType, dependencyId int64) ([]ImageModel, error) {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	selectQuery := `
		SELECT id, imagePath 
		FROM images 
		WHERE dependencyId = $1 AND 
		dependencyType = $2 
		ORDER BY imageOrder
	`

	imageRows, err := db.Query(selectQuery, dependencyId, dependencyType)
	if err != nil {
		return nil, err
	}
	defer imageRows.Close()

	images := make([]ImageModel, 0)

	for imageRows.Next() {
		var imageId int64
		var imagePath string
		err = imageRows.Scan(&imageId, &imagePath)
		if err != nil {
			return nil, err
		}

		images = append(images, ImageModel{Id: imageId, Path: imagePath})
	}

	return images, nil
}

func (repo *ImageRepository) FindImages(dependencyType RepositoryType, dependencyId int64, imageIds []int64) ([]ImageModel, error) {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	selectQuery := `
		SELECT id, imagePath 
		FROM images 
		WHERE dependencyId = $1 AND 
		dependencyType = $2 AND
		id in (%v)
		ORDER BY imageOrder
	`

	imagesStrIds := funk.Map(imageIds, func(id int64) string {
		return fmt.Sprintf("%v", id)
	})
	imageIdsJoined := strings.Join(imagesStrIds.([]string), ",")
	selectQuery = fmt.Sprintf(selectQuery, imageIdsJoined)

	imageRows, err := db.Query(selectQuery, dependencyId, dependencyType)
	if err != nil {
		return nil, err
	}
	defer imageRows.Close()

	images := make([]ImageModel, 0)

	for imageRows.Next() {
		var imageId int64
		var imagePath string
		err = imageRows.Scan(&imageId, &imagePath)
		if err != nil {
			return nil, err
		}

		images = append(images, ImageModel{Id: imageId, Path: imagePath})
	}

	return images, nil
}

func (repo *ImageRepository) FindImageFromPath(dependencyType RepositoryType, dependencyId int64, imagePath string) (*ImageModel, error) {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	selectQuery := `
		SELECT id, imagePath 
		FROM images 
		WHERE dependencyId = $1 AND 
		dependencyType = $2 AND 
		imagePath = $3 
		ORDER BY imageOrder
	`

	rows, err := db.Query(selectQuery, dependencyId, dependencyType, imagePath)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var id int64
	var findImagePath string
	err = rows.Scan(&id, &findImagePath)
	if err != nil {
		return nil, err
	}

	return &ImageModel{Id: id, Path: findImagePath}, nil
}

func (repo *ImageRepository) AddImges(dependencyType RepositoryType, dependencyId int64, images []string) error {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	imageInsertQuery := "INSERT INTO images (dependencyId, dependencyType, imagePath) VALUES ($1, $2, $3)"

	stmt, err := db.Prepare(imageInsertQuery)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, value := range images {

		_, err := stmt.Exec(dependencyId, dependencyType, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *ImageRepository) AddImgesTransaction(tx *sql.Tx, dependencyType RepositoryType, dependencyId int64, images []string) error {

	imageInsertQuery := "INSERT INTO images (dependencyId, dependencyType, imagePath) VALUES ($1, $2, $3)"

	stmt, err := tx.Prepare(imageInsertQuery)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, value := range images {

		_, err := stmt.Exec(dependencyId, dependencyType, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *ImageRepository) RemoveImages(dependencyType RepositoryType, dependencyId int64) error {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	removeImges := "DELETE FROM images WHERE dependencyId = $1 AND dependencyType = $2"
	_, err = db.Exec(removeImges, dependencyId, dependencyType)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ImageRepository) RemoveImagesTransaction(tx *sql.Tx, dependencyType RepositoryType, dependencyId int64) error {

	removeImges := "DELETE FROM images WHERE dependencyId = $1 AND dependencyType = $2"
	_, err := tx.Exec(removeImges, dependencyId, dependencyType)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ImageRepository) RemoveImageFromPath(dependencyType RepositoryType, dependencyId int64, imagePath string) error {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	removeImges := "DELETE FROM images WHERE dependencyId = $1 AND dependencyType = $2 AND imagePath = $3"
	_, err = db.Exec(removeImges, dependencyId, dependencyType, imagePath)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ImageRepository) RemoveImagePathTransaction(tx *sql.Tx, dependencyType RepositoryType, dependencyId int64, imagePath string) error {

	removeImges := "DELETE FROM images WHERE dependencyId = $1 AND dependencyType = $2 AND imagePath = $3"
	_, err := tx.Exec(removeImges, dependencyId, dependencyType, imagePath)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ImageRepository) RemoveImageFromId(dependencyType RepositoryType, dependencyId int64, imageId int64) error {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	removeImges := "DELETE FROM images WHERE dependencyId = $1 AND dependencyType = $2 AND id = $3"
	_, err = db.Exec(removeImges, dependencyId, dependencyType, imageId)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ImageRepository) RemoveImageIdTransaction(tx *sql.Tx, dependencyType RepositoryType, dependencyId int64, imageId int64) error {

	removeImges := "DELETE FROM images WHERE dependencyId = $1 AND dependencyType = $2 AND id = $3"
	_, err := tx.Exec(removeImges, dependencyId, dependencyType, imageId)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ImageRepository) SortImageOrder(dependencyType RepositoryType, dependencyId int64) error {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	selectQuery := "SELECT id FROM images WHERE dependencyId = $1 AND dependencyType = $2 ORDER BY id"
	rows, err := db.Query(selectQuery, dependencyId, dependencyType)
	if err != nil {
		return err
	}

	defer rows.Close()

	orderIndex := 0
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return err
		}

		updateQuery := "UPDATE images SET imageOrder = $4 WHERE dependencyId = $1 AND dependencyType = $2 AND id = $3"
		_, err = db.Exec(updateQuery, dependencyId, dependencyType, id, orderIndex)
		if err != nil {
			return err
		}

		orderIndex++
	}

	return nil
}

func (repo *ImageRepository) SortImageOrderTransation(tx *sql.Tx, dependencyType RepositoryType, dependencyId int64) error {
	selectQuery := "SELECT id FROM images WHERE dependencyId = $1 AND dependencyType = $2 ORDER BY id"
	rows, err := tx.Query(selectQuery, dependencyId, dependencyType)
	if err != nil {
		return err
	}

	defer rows.Close()

	orderIndex := 0
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return err
		}

		updateQuery := "UPDATE images SET imageOrder = $4 WHERE dependencyId = $1 AND dependencyType = $2 AND id = $3"
		_, err = tx.Exec(updateQuery, dependencyId, dependencyType, id, orderIndex)
		if err != nil {
			return err
		}

		orderIndex++
	}

	return nil
}
