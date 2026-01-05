package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/apkatsikas/archiver/db"
	"github.com/apkatsikas/archiver/fileutil"
	"github.com/apkatsikas/archiver/filter"
	"github.com/apkatsikas/archiver/runner"
	storageclient "github.com/apkatsikas/archiver/storage-client"
	flagutil "github.com/apkatsikas/archiver/util/flag"
	"github.com/apkatsikas/archiver/zipper"

	"github.com/apkatsikas/discord-alert/alert"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fu := flagutil.Get()
	fu.Setup()

	switch fu.RunMode {
	case flagutil.RunModeLedger:
		if err := runLedger(); err != nil {
			panic(err)
		}
		return
	case flagutil.RunModeBatch:
		if err := runBatch(); err != nil {
			panic(err)
		}
		return
	}

	err := performScheduledArchive()
	if err != nil {
		errMsg := fmt.Sprintf("\nARCHIVER: Failed to perform scheduled archive - %v", err.Error())
		alertErr := alert.SendAlert(errMsg)
		if alertErr != nil {
			log.Printf("Encountered issue trying to send alert %v", alertErr)
		}
		panic(err)
	}
}

func runLedger() error {
	arguments := flag.Args()

	if len(arguments) < 2 {
		return fmt.Errorf("ledger mode requires arguments for sqlite DB file and output file")
	}

	sqliteDbFile := arguments[0]
	ledgerOutputFile := arguments[1]

	sqliteHandler := &db.SQLiteHandler{}
	sqliteHandler.ConnectSQLite(sqliteDbFile)
	mfr := &db.MusicFoldersRepository{SqliteHandler: sqliteHandler}
	libraryRepository := &db.LibraryRepository{SqliteHandler: sqliteHandler}
	runn := &runner.Runner{MusicFoldersRepository: mfr,
		LibraryRepository: libraryRepository, FileSystemOperator: &fileutil.FileSystemOperator{}}

	if err := runn.BuildLedger(ledgerOutputFile); err != nil {
		return err
	}
	return nil
}

func runBatch() error {
	arguments := flag.Args()

	if len(arguments) < 1 {
		return fmt.Errorf("batch mode requires arguments for ledger file")
	}

	ledgerFile := arguments[0]

	fso := &fileutil.FileSystemOperator{}
	runn := &runner.Runner{
		FileSystemOperator: fso,
		Zipper: &zipper.Zipper{
			FileSystemOperator: fso,
		},
		StorageClient: storageclient.New(),
	}

	if err := runn.RunBatch(ledgerFile); err != nil {
		return err
	}
	return nil
}

func performScheduledArchive() error {
	arguments := flag.Args()

	if len(arguments) < 2 {
		return fmt.Errorf("scheduled mode requires arguments for navidrome DB path and archive DB path")
	}

	navidromeDbPath := arguments[0]
	archiveDbPath := arguments[1]

	sqliteHandlerArchiveRun := &db.SQLiteHandler{}
	if err := sqliteHandlerArchiveRun.ConnectSQLite(archiveDbPath); err != nil {
		return err
	}

	sqliteNavidrome := &db.SQLiteHandler{}
	if err := sqliteNavidrome.ConnectSQLite(navidromeDbPath); err != nil {
		return err
	}

	fso := &fileutil.FileSystemOperator{}

	runn := &runner.Runner{
		FilterService:          &filter.FilterService{},
		StorageClient:          storageclient.New(),
		MusicFoldersRepository: &db.MusicFoldersRepository{SqliteHandler: sqliteNavidrome},
		ArchiveRunRepository:   &db.ArchiveRunRepository{SqliteHandler: sqliteHandlerArchiveRun},
		AdminRepository:        &db.AdminRepository{SqliteHandler: sqliteNavidrome},
		LibraryRepository:      &db.LibraryRepository{SqliteHandler: sqliteNavidrome},
		Zipper: &zipper.Zipper{
			FileSystemOperator: fso,
		},
		FileSystemOperator: fso,
	}

	if err := runn.RunScheduled(); err != nil {
		return err
	}
	return nil
}
