package db

import (
	"database/sql"
	"time"
)

type MediaFile struct {
	Id        string
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
	LibraryId int
}

type MusicFoldersRepository struct {
	SqliteHandler *SQLiteHandler
}

func (pr *MusicFoldersRepository) NewMediaFilesSinceDate(cutoffTime time.Time) ([]MediaFile, error) {
	statement, err := pr.SqliteHandler.Db().Prepare(
		"SELECT id, path, created_at, library_id FROM media_file WHERE created_at > ?")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(cutoffTime.Format(timeFormat))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []MediaFile
	for rows.Next() {
		var mediaFile MediaFile
		if err := rows.Scan(
			&mediaFile.Id, &mediaFile.Path, &mediaFile.CreatedAt, &mediaFile.LibraryId); err != nil {
			return nil, err
		}
		all = append(all, mediaFile)
	}
	return all, nil
}

func (mfr *MusicFoldersRepository) UpdatedMediaFilesSinceDate(cutoffTime time.Time) ([]MediaFile, error) {
	statement, err := mfr.SqliteHandler.Db().Prepare(
		"SELECT id, path, created_at, updated_at, library_id " +
			"FROM media_file WHERE updated_at > created_at AND updated_at > ?")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(cutoffTime.Format(timeFormat))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	all, err := mfr.mediaFilesFromRows(rows)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (mfr *MusicFoldersRepository) AllMediaFiles() ([]MediaFile, error) {
	rows, err := mfr.SqliteHandler.Db().Query("SELECT id, path, created_at, updated_at, library_id FROM media_file")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	all, err := mfr.mediaFilesFromRows(rows)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (mfr *MusicFoldersRepository) mediaFilesFromRows(rows *sql.Rows) ([]MediaFile, error) {
	var all []MediaFile
	for rows.Next() {
		var mediaFile MediaFile
		if err := rows.Scan(
			&mediaFile.Id, &mediaFile.Path, &mediaFile.CreatedAt,
			&mediaFile.UpdatedAt, &mediaFile.LibraryId); err != nil {
			return nil, err
		}
		all = append(all, mediaFile)
	}
	return all, nil
}
