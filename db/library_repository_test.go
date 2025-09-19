package db_test

import (
	"github.com/apkatsikas/archiver/db"

	testutils "github.com/apkatsikas/archiver/tests/test-utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/mattn/go-sqlite3"
)

var _ = Describe("LastRun", func() {
	var testDbFullPath = ""
	var err error
	var sqliteHandler *db.SQLiteHandler

	var libraryRepository *db.LibraryRepository

	BeforeEach(func() {
		By("Resetting and connecting to DB")

		testDbFullPath, err = testutils.SetupTestDb("fakenavidrome")
		Expect(err).To(BeNil(), "Error trying to setup DB")

		sqliteHandler = &db.SQLiteHandler{}
		Expect(sqliteHandler.ConnectSQLite(testDbFullPath)).To(BeNil(), "Failed to connect to sqlite")
		libraryRepository = &db.LibraryRepository{SqliteHandler: sqliteHandler}
	})

	It("Returns a library record", func() {
		Expect(libraryRepository.LibraryById(1)).To(Equal(&db.Library{
			Id:   1,
			Path: "/lib/path",
		}))
	})
})
