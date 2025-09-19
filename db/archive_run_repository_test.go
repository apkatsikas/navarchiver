package db_test

import (
	"os"
	"path/filepath"
	"time"

	"github.com/apkatsikas/archiver/db"
	testutils "github.com/apkatsikas/archiver/tests/test-utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/mattn/go-sqlite3"
)

const fakedb = "fakearchiverun"

var lastRun = time.Date(2024, time.January, 12, 13, 19, 51, 0, time.UTC)

var _ = Describe("LastRun", func() {
	var testDbFullPath = ""
	var err error
	var sqliteHandler *db.SQLiteHandler

	var archiveRunRepository *db.ArchiveRunRepository

	BeforeEach(func() {
		By("Resetting and connecting to DB")

		testDbFullPath, err = testutils.SetupTestDb(fakedb)
		Expect(err).To(BeNil(), "Error trying to setup DB")

		sqliteHandler = &db.SQLiteHandler{}
		Expect(sqliteHandler.ConnectSQLite(testDbFullPath)).To(BeNil(), "Failed to connect to sqlite")
		archiveRunRepository = &db.ArchiveRunRepository{SqliteHandler: sqliteHandler}
	})

	Context("When there is data in the repository", func() {
		BeforeEach(func() {
			Expect(testutils.UpdateLastRunRecord(testDbFullPath, lastRun)).To(BeNil(), "Failed to update archive_run")
		})

		It("Returns the expected archive run", func() {
			Expect(archiveRunRepository.LastRun()).To(Equal(&db.ArchiveRun{LastRun: lastRun, Id: 1}))
		})
	})

	Context("When there is no data in the repository", func() {
		It("Returns nil", func() {
			Expect(archiveRunRepository.LastRun()).To(BeNil())
		})
	})
})

var _ = Describe("UpdateLastRun", func() {
	var testDbFullPath = ""
	var err error
	var sqliteHandler *db.SQLiteHandler

	var archiveRunRepository *db.ArchiveRunRepository
	var newLastRun = lastRun.Add(time.Duration(48) * time.Hour)

	BeforeEach(func() {
		By("Resetting and connecting to DB")

		testDbFullPath, err = testutils.SetupTestDb(fakedb)
		Expect(err).To(BeNil(), "Error trying to setup DB")

		sqliteHandler = &db.SQLiteHandler{}
		Expect(sqliteHandler.ConnectSQLite(testDbFullPath)).To(BeNil(), "Failed to connect to sqlite")
		archiveRunRepository = &db.ArchiveRunRepository{SqliteHandler: sqliteHandler}
	})

	Context("When there is data in the repository", func() {
		BeforeEach(func() {
			Expect(testutils.UpdateLastRunRecord(testDbFullPath, lastRun)).To(BeNil(), "Failed to update archive_run")
			Expect(archiveRunRepository.UpdateLastRun(newLastRun)).To(BeNil(), "Failed to run UpdateLastRun")
		})

		It("Correctly sets the expected data", func() {
			Expect(testutils.LastRunRecord(testDbFullPath)).To(Equal(newLastRun))
		})
	})

	Context("When there is no data in the repository", func() {
		BeforeEach(func() {
			Expect(archiveRunRepository.UpdateLastRun(newLastRun)).To(BeNil(), "Failed to run UpdateLastRun")
		})

		It("Correctly sets the expected data", func() {
			Expect(testutils.LastRunRecord(testDbFullPath)).To(Equal(newLastRun))
		})
	})
})

var _ = Describe("CreateTable", func() {
	var testDbFullPath = ""
	var testDir = ""
	var err error
	var sqliteHandler *db.SQLiteHandler

	var archiveRunRepository *db.ArchiveRunRepository

	Context("When there is no Archive Run table", func() {
		BeforeEach(func() {
			By("Get working directory")
			testDir, err = os.Getwd()
			Expect(err).To(BeNil(), "Got an error getting working directory")

			By("Resetting and connecting to DB")
			testDbFullPath = filepath.Join(testDir, "..", "tests/fixtures", "missing_fakenavidrome.db")
			err = testutils.RemoveFileIfExists(testDbFullPath)
			Expect(err).To(BeNil(), "Got an error removing missing_fakenavidrome")
			Expect(testutils.CheckForArchiveRun(testDbFullPath)).To(BeEmpty(), "Archive run table existed after setting up")

			By("Setting up and running repo")
			sqliteHandler = &db.SQLiteHandler{}
			Expect(sqliteHandler.ConnectSQLite(testDbFullPath)).To(BeNil(), "Failed to connect to sqlite")
			archiveRunRepository = &db.ArchiveRunRepository{SqliteHandler: sqliteHandler}
			err = archiveRunRepository.CreateTable()
		})

		It("Returns nil", func() {
			Expect(err).To(BeNil())
		})

		It("Creates the Archive Run table", func() {
			Expect(testutils.CheckForArchiveRun(testDbFullPath)).To(Equal("archive_run"))
		})

		It("Returns no Last Run record", func() {
			Expect(testutils.LastRunRecordString(testDbFullPath)).To(BeEmpty())
		})
	})

	Context("When there is an Archive Run table with a record", func() {
		BeforeEach(func() {
			By("Get working directory")
			testDir, err = os.Getwd()
			Expect(err).To(BeNil(), "Got an error getting working directory")

			By("Resetting and connecting to DB")
			testDbFullPath, err = testutils.SetupTestDb(fakedb)
			Expect(testutils.UpdateLastRunRecord(testDbFullPath, lastRun)).To(BeNil(), "Failed to update archive_run")

			By("Setting up and running repo")
			sqliteHandler = &db.SQLiteHandler{}
			Expect(sqliteHandler.ConnectSQLite(testDbFullPath)).To(BeNil(), "Failed to connect to sqlite")
			archiveRunRepository = &db.ArchiveRunRepository{SqliteHandler: sqliteHandler}
			err = archiveRunRepository.CreateTable()
		})

		It("Returns nil", func() {
			Expect(err).To(BeNil())
		})

		It("Returns the expected Last Run record", func() {
			Expect(testutils.LastRunRecord(testDbFullPath)).To(Equal(lastRun))
		})
	})
})
