package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/thoas/go-funk"
)

type AboutModel struct {
	ProfileImage     *string `json:"profile_image"`
	ProfileName      *string `json:"profile_name"`
	Contact          *string `json:"contact"`
	IntroduceContent *string `json:"introduce_content"`
}

func (aboutModel *AboutModel) GetProfileImagePath(saveDir string, prefixUri string) string {

	if aboutModel.ProfileImage == nil {
		return ""
	}

	return strings.Replace(*aboutModel.ProfileImage, prefixUri, saveDir, 1)
}

type AboutHistoryModel struct {
	Id       int64  `json:"id"`
	Category string `json:"category"`
	Duration string `json:"duration"`
	Content  string `json:"content"`
}

type AboutHistoryContent struct {
	Category string
	Duration string
	Content  string
}

type AboutHistoryIdContent struct {
	Id int64
	AboutHistoryContent
}

type AboutRepository struct {
	DBConnect *DBConnection
}

func (repo *AboutRepository) getAboutId() (int64, error) {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return -1, err
	}

	countRow := db.QueryRow("SELECT id FROM about")
	if countRow == nil {
		return -1, fmt.Errorf("about count query error")
	}

	if countRow.Err() != nil {
		return -1, countRow.Err()
	}

	var aboutId int64 = -1
	countRow.Scan(&aboutId)

	return aboutId, nil
}

func (repo *AboutRepository) createAbout() (int64, error) {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return -1, err
	}

	aboutCreateQuery := "INSERT INTO about (profileImage, profileName, Contact, IntroduceContent) VALUES (NULL, NULL, NULL, NULL)"

	result, err := db.Exec(aboutCreateQuery)
	if err != nil {
		return -1, err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return insertId, nil
}

func (repo *AboutRepository) getProfileImage() (*string, error) {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	aboutId, err := repo.getAboutId()
	if err != nil {
		return nil, err
	}

	if aboutId == -1 {
		return nil, nil
	}

	profileIamgeRow := db.QueryRow("SELECT profileImage FROM about WHERE id = $1", aboutId)
	if profileIamgeRow == nil {
		return nil, nil
	}

	if profileIamgeRow.Err() != nil {
		return nil, profileIamgeRow.Err()
	}

	var profileImage *string
	profileIamgeRow.Scan(&profileImage)

	return profileImage, nil
}

func (repo *AboutRepository) CreateTable() error {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	createAboutTableQuery := `
		CREATE TABLE IF NOT EXISTS "about"
		(
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"profileImage" TEXT,
			"profileName" TEXT,
			"contact" TEXT,
			"introduceContent" TEXT
		)`

	_, err = db.Exec(createAboutTableQuery)
	if err != nil {
		log.Printf("[error] create table potofolio [%v]\n", err)
		return err
	}

	createAboutHistoryTableQuery := `
		CREATE TABLE IF NOT EXISTS "about_history"
		(
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"category" TEXT,
			"duration" TEXT,
			"content" TEXT
		)`

	_, err = db.Exec(createAboutHistoryTableQuery)
	if err != nil {
		log.Printf("[error] create table potofolio [%v]\n", err)
		return err
	}

	return nil
}

func (repo *AboutRepository) GetAbout() (*AboutModel, error) {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	aboutRow, err := db.Query("SELECT profileImage, profileName, Contact, IntroduceContent FROM about LIMIT 1")
	if err != nil {
		log.Printf("[error] about query [%v]\n", err)
		return nil, err
	}

	defer aboutRow.Close()

	if !aboutRow.Next() {
		return &AboutModel{}, nil
	}

	about := AboutModel{}

	err = aboutRow.Scan(&about.ProfileImage, &about.ProfileName, &about.Contact, &about.IntroduceContent)
	if err != nil {
		log.Printf("[error] about scan [%v]\n", err)
		return nil, err
	}

	return &about, nil
}

func (repo *AboutRepository) GetHistory() ([]AboutHistoryModel, error) {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	aboutHistoryRows, err := db.Query("SELECT id, category, duration, content FROM about_history ORDER BY id")
	if err != nil {
		return nil, err
	}

	defer aboutHistoryRows.Close()

	aboutHistoryModels := make([]AboutHistoryModel, 0)

	for aboutHistoryRows.Next() {

		aboutHistoryModel := AboutHistoryModel{}
		err = aboutHistoryRows.Scan(&aboutHistoryModel.Id, &aboutHistoryModel.Category, &aboutHistoryModel.Duration, &aboutHistoryModel.Content)
		if err != nil {
			log.Printf("[error] aboutHistory scan [%v]\n", err)
			continue
		}

		aboutHistoryModels = append(aboutHistoryModels, aboutHistoryModel)
	}

	return aboutHistoryModels, nil
}

func (repo *AboutRepository) UpdateAbout(profileImage *string, profileName *string, contact *string, introduceContent *string) (*string, error) {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	complete := false

	transaction, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer CloseTranstion(transaction, &complete)

	aboutId, err := repo.getAboutId()
	if err != nil {
		return nil, err
	}

	if aboutId == -1 {
		aboutId, err = repo.createAbout()
		if err != nil {
			return nil, err
		}
	}

	var removePorfileIamge *string
	if profileImage != nil {

		removePorfileIamge, _ = repo.getProfileImage()

		profileImageUpdateQuery := fmt.Sprintf("UPDATE about SET profileImage = \"%v\" WHERE id = $1", *profileImage)

		_, err = transaction.Exec(profileImageUpdateQuery, aboutId)
		if err != nil {
			return nil, err
		}
	}

	if profileName != nil {
		profileNameUpdateQuery := fmt.Sprintf("UPDATE about SET profileName = \"%v\" WHERE id = $1", *profileName)

		_, err = transaction.Exec(profileNameUpdateQuery, aboutId)
		if err != nil {
			return nil, err
		}
	}

	if contact != nil {
		profileContactUpdateQuery := fmt.Sprintf("UPDATE about SET contact = \"%v\" WHERE id = $1", *contact)

		_, err = transaction.Exec(profileContactUpdateQuery, aboutId)
		if err != nil {
			return nil, err
		}
	}

	if introduceContent != nil {
		profileContactUpdateQuery := fmt.Sprintf("UPDATE about SET introduceContent = \"%v\" WHERE id = $1", *introduceContent)

		_, err = transaction.Exec(profileContactUpdateQuery, aboutId)
		if err != nil {
			return nil, err
		}
	}

	complete = true
	return removePorfileIamge, nil
}

func (repo *AboutRepository) UpdateAboutHistory(removeId []int64, updateHistoryInfos []AboutHistoryIdContent, addHistoryInfos []AboutHistoryContent) error {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	complete := false
	defer CloseTranstion(transaction, &complete)

	if len(removeId) > 0 {
		removeIdStrs := funk.Map(removeId, func(e int64) string {
			return fmt.Sprintf("%d", e)
		})

		removeIdJoined := strings.Join(removeIdStrs.([]string), ",")
		removeQuery := fmt.Sprintf("DELETE FROM about_history WHERE id in (%s)", removeIdJoined)

		_, err := transaction.Exec(removeQuery)
		if err != nil {
			return err
		}
	}

	if len(updateHistoryInfos) > 0 {

		for _, updateHistoryInfo := range updateHistoryInfos {
			update_history_query := fmt.Sprintf("UPDATE about_history SET category = '%s', duration = '%s', content = '%s' WHERE id = $1",
				updateHistoryInfo.Category, updateHistoryInfo.Duration, updateHistoryInfo.Content)

			_, err = transaction.Exec(update_history_query, updateHistoryInfo.Id)
			if err != nil {
				return err
			}
		}
	}

	if len(addHistoryInfos) > 0 {

		stmt, err := transaction.Prepare("INSERT INTO about_history (category, duration, content) VALUES ($1, $2, $3)")
		if err != nil {
			return err
		}

		defer stmt.Close()

		for _, addHistoryInfo := range addHistoryInfos {
			_, err = stmt.Exec(addHistoryInfo.Category, addHistoryInfo.Duration, addHistoryInfo.Content)
			if err != nil {
				return err
			}
		}
	}

	complete = true

	return nil
}
