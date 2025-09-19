package db

import (
	"database/sql"
)

type SQLiteHandler struct {
	db *sql.DB
}

func (handler *SQLiteHandler) Db() *sql.DB {
	return handler.db
}

func (handler *SQLiteHandler) ConnectSQLite(dsn string) error {
	database, err := sql.Open("sqlite3", dsn)

	if err != nil {
		return err
	}
	handler.db = database
	return nil
}
