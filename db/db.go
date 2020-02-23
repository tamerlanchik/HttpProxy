package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbname = "db/proxy.db"
)

var (
	DB *sql.DB
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}