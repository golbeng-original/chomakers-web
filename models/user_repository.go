package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

type UserNotExistError struct{}

func (e *UserNotExistError) Error() string {
	return "user not exist"
}

type UserModel struct {
	Id           int64
	UserName     string
	Password     string
	RefreshToken *string
}

type UserRespository struct {
	DBConnect *DBConnection
}

func (repo *UserRespository) CreateTable() error {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	createUserTableQuery := `
		CREATE TABLE IF NOT EXISTS "user"
		(
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"username" TEXT,
			"password" TEXT,
			"refreshtoken" TEXT
		)`

	_, err = db.Exec(createUserTableQuery)
	if err != nil {
		log.Printf("[error] create table user [%v]\n", err)
		return err
	}

	return nil
}

func (repo *UserRespository) IsExist(username string) (bool, error) {

	userModel, err := repo.GetUserModelFromUserName(username)
	if err != nil {
		return false, err
	}

	if userModel == nil {
		return false, nil
	}

	return true, nil
}

func (repo *UserRespository) AddUser(username, password string) error {

	hash := md5.New()
	hash.Write([]byte(password))

	digest := hash.Sum(nil)
	md5Password := hex.EncodeToString(digest)

	return repo.AddUserMd5(username, md5Password)
}

func (repo *UserRespository) AddUserMd5(username, password string) error {

	_, err := repo.IsExist(username)

	if err != nil {
		if !errors.Is(err, &UserNotExistError{}) {
			return err
		}
	}

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	userInsertQuery := "INSERT INTO user (username, password) VALUES ($1, $2)"

	_, err = db.Exec(userInsertQuery, username, password)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRespository) GetUserModelFromUserName(username string) (*UserModel, error) {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT id, username, password, refreshtoken FROM user WHERE username = $1", username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, &UserNotExistError{}
	}

	userModel := UserModel{}
	rows.Scan(&userModel.Id, &userModel.UserName, &userModel.Password, &userModel.RefreshToken)

	return &userModel, nil
}

func (repo *UserRespository) GetUserModel(userId int64) (*UserModel, error) {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM user WHERE id = $1", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, &UserNotExistError{}
	}

	userModel := UserModel{}
	rows.Scan(&userModel.Id, &userModel.UserName, &userModel.Password, &userModel.RefreshToken)

	return &userModel, nil
}

func (repo *UserRespository) SetRefreshToken(userId int64, refreshToken string) error {

	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	tokenUpdateQuery := fmt.Sprintf("UPDATE user SET refreshtoken = '%s' WHERE id = $1", refreshToken)

	_, err = db.Exec(tokenUpdateQuery, userId)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRespository) ClearRefreshToken(userId int64) error {
	db, err := repo.DBConnect.GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE user SET refreshtoken = '' WHERE id = $1", userId)
	if err != nil {
		return err
	}

	return nil
}
