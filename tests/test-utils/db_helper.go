package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const sqliteBinary = "sqlite3"
const dateFormat = "2006-01-02T15:04:05"

func SetupTestDb(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %v", err)
	}

	// If we want to parallelize these, we can use a random string when creating the DB
	// and delete it after the test run. As of now, they operate on the same file - could collide on
	// a parallel run
	fixturesPath := "tests/fixtures"
	testDbScript := fmt.Sprintf("%v/%v.sql", fixturesPath, name)
	testDbScriptFullPath := filepath.Join(dir, "..", testDbScript)
	testDbScriptReadCommand := fmt.Sprintf(".read %v", testDbScriptFullPath)
	testDbPath := fmt.Sprintf("%v/%v.db", fixturesPath, name)
	testDbFullPath := filepath.Join(dir, "..", testDbPath)

	_, err = os.Stat(testDbFullPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("got unexpected error checking if DB exists %v", err)
		}
	} else {
		err = os.Remove(testDbFullPath)
		if err != nil {
			return "", fmt.Errorf("got an error trying to delete DB file: %v", err)
		}
	}

	cmd := exec.Command(sqliteBinary, testDbFullPath, testDbScriptReadCommand)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("got an error trying to setup sqlite DB: %v - command was %v", err, cmd)
	}

	return testDbFullPath, nil
}

func UpdateMediaFileRecordCreatedAt(dbPath string, timeDiffMinutes int, path string, cutoff time.Time) error {
	createdAtTime := cutoff.Add(time.Duration(timeDiffMinutes) * time.Minute).Format(dateFormat)
	statement := fmt.Sprintf("UPDATE media_file SET created_at = '%v' WHERE path LIKE '%%%v%%';", createdAtTime, path)
	cmd := exec.Command(sqliteBinary, dbPath, statement)
	err := cmd.Run()
	return err
}

func UpdateMediaFileRecordCreatedAndUpdatedAt(
	dbPath string, timeDiffMinutesCreatedAt int, timeDiffMinutesUpdatedAt int, path string, cutoff time.Time) error {
	createdAtTime := cutoff.Add(time.Duration(timeDiffMinutesCreatedAt) * time.Minute).Format(dateFormat)
	updatedAtTime := cutoff.Add(time.Duration(timeDiffMinutesUpdatedAt) * time.Minute).Format(dateFormat)

	statement := fmt.Sprintf(
		"UPDATE media_file SET created_at = '%v', updated_at = '%v' WHERE path LIKE '%%%v%%';", createdAtTime, updatedAtTime, path)
	cmd := exec.Command(sqliteBinary, dbPath, statement)
	err := cmd.Run()
	return err
}

func UpdateMediaFilePathPrefix(
	dbPath string, pathPrefix string, path string) error {
	statement := fmt.Sprintf(
		"UPDATE media_file SET path = '%v/%v' WHERE path LIKE '%%%v%%';", pathPrefix, path, path)
	cmd := exec.Command(sqliteBinary, dbPath, statement)
	err := cmd.Run()
	return err
}

func UpdateLastRunRecord(dbPath string, date time.Time) error {
	statement := fmt.Sprintf(
		"INSERT INTO archive_run (\"last_run\") VALUES ('%v');",
		date.Format(dateFormat))
	cmd := exec.Command(sqliteBinary, dbPath, statement)
	err := cmd.Run()
	return err
}

func DeleteLastRunRecord(dbPath string) error {
	statement := "DELETE FROM archive_run;"

	cmd := exec.Command(sqliteBinary, dbPath, statement)
	err := cmd.Run()
	return err
}

func LastRunRecord(dbPath string) (time.Time, error) {
	cmd := exec.Command(sqliteBinary, dbPath, "SELECT last_run FROM archive_run LIMIT 1;")
	result, err := cmd.Output()
	if err != nil {
		return time.Now(), err
	}
	date, err := time.Parse(dateFormat, strings.TrimSpace(string(result)))
	if err != nil {
		return time.Now(), err
	}
	return date, nil
}

func LastRunRecordString(dbPath string) (string, error) {
	cmd := exec.Command(sqliteBinary, dbPath, "SELECT last_run FROM archive_run LIMIT 1;")
	result, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func CheckForArchiveRun(dbPath string) (string, error) {
	cmd := exec.Command(sqliteBinary, dbPath, "SELECT name FROM sqlite_master WHERE type='table' AND name='archive_run';")
	result, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), nil
}
