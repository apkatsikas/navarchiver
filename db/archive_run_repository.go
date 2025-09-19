package db

import (
	"database/sql"
	"errors"
	"time"
)

type ArchiveRun struct {
	LastRun time.Time
	Id      int
}

type ArchiveRunRepository struct {
	SqliteHandler *SQLiteHandler
}

func (arr *ArchiveRunRepository) LastRun() (*ArchiveRun, error) {
	var archiveRun ArchiveRun
	err := arr.SqliteHandler.Db().QueryRow(
		"SELECT id, last_run FROM archive_run").Scan(&archiveRun.Id, &archiveRun.LastRun)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &archiveRun, nil
}

func (arr *ArchiveRunRepository) UpdateLastRun(lastRun time.Time) error {
	statement, err := arr.SqliteHandler.Db().Prepare(
		"INSERT OR REPLACE INTO archive_run (id, last_run) VALUES (1, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(lastRun.Format(timeFormat))
	if err != nil {
		return err
	}
	return nil
}

func (arr *ArchiveRunRepository) CreateTable() error {
	_, err := arr.SqliteHandler.Db().Exec(
		"CREATE TABLE IF NOT EXISTS archive_run (id INTEGER PRIMARY KEY AUTOINCREMENT,last_run DATE NOT NULL,CONSTRAINT id_unique UNIQUE (id));")
	if err != nil {
		return err
	}
	return nil
}
