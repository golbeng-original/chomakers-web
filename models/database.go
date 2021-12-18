package models

import (
	"database/sql"
	"fmt"

	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type DBConnection struct {
	db *sql.DB
}

func (connect *DBConnection) Open(dataSourceName string) error {

	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	connect.db = db

	return nil
}

func (connect *DBConnection) GetDB() (*sql.DB, error) {

	if connect.db == nil {
		return nil, errors.New("db is null")
	}

	return connect.db, nil
}

func (connect *DBConnection) Close() {

	if connect.db == nil {
		return
	}

	connect.db.Close()
}

func CloseTranstion(tx *sql.Tx, completed *bool) {

	var err error
	if !*completed {
		err = tx.Rollback()
	} else {
		err = tx.Commit()
	}

	if err != nil {
		fmt.Println(err)
	}
}