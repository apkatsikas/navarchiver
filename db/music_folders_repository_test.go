package db_test

import (
	"fmt"
	"time"

	"github.com/apkatsikas/archiver/db"
	testutils "github.com/apkatsikas/archiver/tests/test-utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/mattn/go-sqlite3"
)

const crazyRhythms = "music/Crazy Rhythms"
const guyIncognito = "music/Guy Incognito - Lovedrug"

var lastQueryEnded = time.Date(2024, time.January, 12, 13, 19, 51, 0, time.UTC)

type mediaFolderTestData struct {
	createdAtTimeRelativeToCutoffTime int
	updatedAtTimeRelativeToCutoffTime int
	folderName                        string
}

var _ = Describe("NewMediaFilesSinceDate", func() {
	var testDbFullPath = ""
	var err error
	var sqliteHandler *db.SQLiteHandler

	BeforeEach(func() {
		By("Resetting and connecting to DB")

		testDbFullPath, err = testutils.SetupTestDb("fakenavidrome")
		Expect(err).To(BeNil(), "Error trying to setup DB")
		sqliteHandler = &db.SQLiteHandler{}
		Expect(sqliteHandler.ConnectSQLite(testDbFullPath)).To(BeNil(), "Failed to connect to sqlite")
	})

	var mediaFolderScenarios = []TableEntry{
		Entry("when all media files are created after cutoff time", []mediaFolderTestData{
			{
				createdAtTimeRelativeToCutoffTime: 10,
				folderName:                        crazyRhythms,
			},
			{
				createdAtTimeRelativeToCutoffTime: 20,
				folderName:                        guyIncognito,
			},
		}, []db.MediaFile{
			{
				Id:        "5c214deb5b2dba739e0d6af56f61d1c7",
				Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 09 - crazy rhythms.mp3",
				CreatedAt: timeParse("2024-01-12T13:29:51"),
				LibraryId: 1,
			},
			{
				Id:        "37141ae2932c8e06cc3716c3b9c55a48",
				Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 15 - i wanna sleep in your arms (live).mp3",
				CreatedAt: timeParse("2024-01-12T13:29:51"),
				LibraryId: 1,
			},
			{
				Id:        "54c5999927b56e2887c3a5cfd21bdfbf",
				Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 13 - moscow nights (demo).mp3",
				CreatedAt: timeParse("2024-01-12T13:29:51"),
				LibraryId: 1,
			},
			{
				Id:        "1fde8382304840139358a101b081db9c",
				Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 02 Can We Please Go Back To How It Was Before I messed Up-.mp3",
				CreatedAt: timeParse("2024-01-12T13:39:51"),
				LibraryId: 1,
			},
			{
				Id:        "07bb5aec148d087b0192c538721d0627",
				Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 03 Seven.mp3",
				CreatedAt: timeParse("2024-01-12T13:39:51"),
				LibraryId: 1,
			},
			{
				Id:        "6ea5a2baa32842109925f67b3151fb80",
				Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 01 Lovedrug.mp3",
				CreatedAt: timeParse("2024-01-12T13:39:51"),
				LibraryId: 1,
			},
		}),
		Entry("when some media files are created after cutoff time", []mediaFolderTestData{
			{
				createdAtTimeRelativeToCutoffTime: -10,
				folderName:                        crazyRhythms,
			},
			{
				createdAtTimeRelativeToCutoffTime: 5,
				folderName:                        guyIncognito,
			},
		}, []db.MediaFile{
			{
				Id:        "1fde8382304840139358a101b081db9c",
				Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 02 Can We Please Go Back To How It Was Before I messed Up-.mp3",
				CreatedAt: timeParse("2024-01-12T13:24:51"),
				LibraryId: 1,
			},
			{
				Id:        "07bb5aec148d087b0192c538721d0627",
				Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 03 Seven.mp3",
				CreatedAt: timeParse("2024-01-12T13:24:51"),
				LibraryId: 1,
			},
			{
				Id:        "6ea5a2baa32842109925f67b3151fb80",
				Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 01 Lovedrug.mp3",
				CreatedAt: timeParse("2024-01-12T13:24:51"),
				LibraryId: 1,
			},
		}),
		Entry("when no media files are created after cutoff time", []mediaFolderTestData{
			{
				createdAtTimeRelativeToCutoffTime: -10,
				folderName:                        crazyRhythms,
			},
			{
				createdAtTimeRelativeToCutoffTime: -5,
				folderName:                        guyIncognito,
			},
		}, nil),
	}

	DescribeTableSubtree("media folder scenarios",
		func(testDataRecords []mediaFolderTestData, expectedMediaFiles []db.MediaFile) {
			var musicFoldersRepository *db.MusicFoldersRepository

			BeforeEach(func() {
				for _, td := range testDataRecords {
					Expect(
						testutils.UpdateMediaFileRecordCreatedAt(testDbFullPath, td.createdAtTimeRelativeToCutoffTime, td.folderName, lastQueryEnded)).
						To(BeNil(), "Failed to update record %v", td.folderName)
				}
				musicFoldersRepository = &db.MusicFoldersRepository{SqliteHandler: sqliteHandler}
			})

			It("should return the expected media files", func() {
				Expect(musicFoldersRepository.NewMediaFilesSinceDate(lastQueryEnded)).To(ConsistOf(expectedMediaFiles))
			})

		}, mediaFolderScenarios)
})

var _ = Describe("UpdatedMediaFilesSinceDate", func() {
	var testDbFullPath = ""

	var err error
	var sqliteHandler *db.SQLiteHandler

	BeforeEach(func() {
		By("Resetting and connecting to DB")

		testDbFullPath, err = testutils.SetupTestDb("fakenavidrome")
		Expect(err).To(BeNil(), "Error trying to setup DB")
		sqliteHandler = &db.SQLiteHandler{}
		Expect(sqliteHandler.ConnectSQLite(testDbFullPath)).To(BeNil(), "Failed to connect to sqlite")
	})

	var mediaFolderScenarios = []TableEntry{
		Entry("when all media files are updated after cutoff time and created earlier than updated", []mediaFolderTestData{
			{
				createdAtTimeRelativeToCutoffTime: -10,
				updatedAtTimeRelativeToCutoffTime: 10,
				folderName:                        crazyRhythms,
			},
			{
				createdAtTimeRelativeToCutoffTime: -10,
				updatedAtTimeRelativeToCutoffTime: 20,
				folderName:                        guyIncognito,
			},
		}, []db.MediaFile{
			{
				Id:        "5c214deb5b2dba739e0d6af56f61d1c7",
				Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 09 - crazy rhythms.mp3",
				CreatedAt: timeParse("2024-01-12T13:09:51"),
				UpdatedAt: timeParse("2024-01-12T13:29:51"),
				LibraryId: 1,
			},
			{
				Id:        "37141ae2932c8e06cc3716c3b9c55a48",
				Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 15 - i wanna sleep in your arms (live).mp3",
				CreatedAt: timeParse("2024-01-12T13:09:51"),
				UpdatedAt: timeParse("2024-01-12T13:29:51"),
				LibraryId: 1,
			},
			{
				Id:        "54c5999927b56e2887c3a5cfd21bdfbf",
				Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 13 - moscow nights (demo).mp3",
				CreatedAt: timeParse("2024-01-12T13:09:51"),
				UpdatedAt: timeParse("2024-01-12T13:29:51"),
				LibraryId: 1,
			},
			{
				Id:        "1fde8382304840139358a101b081db9c",
				Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 02 Can We Please Go Back To How It Was Before I messed Up-.mp3",
				CreatedAt: timeParse("2024-01-12T13:09:51"),
				UpdatedAt: timeParse("2024-01-12T13:39:51"),
				LibraryId: 1,
			},
			{
				Id:        "07bb5aec148d087b0192c538721d0627",
				Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 03 Seven.mp3",
				CreatedAt: timeParse("2024-01-12T13:09:51"),
				UpdatedAt: timeParse("2024-01-12T13:39:51"),
				LibraryId: 1,
			},
			{
				Id:        "6ea5a2baa32842109925f67b3151fb80",
				Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 01 Lovedrug.mp3",
				CreatedAt: timeParse("2024-01-12T13:09:51"),
				UpdatedAt: timeParse("2024-01-12T13:39:51"),
				LibraryId: 1,
			},
		}),
		Entry("when some media files are updated after cutoff time and all are created earlier than updated", []mediaFolderTestData{
			{
				createdAtTimeRelativeToCutoffTime: -10,
				updatedAtTimeRelativeToCutoffTime: 10,
				folderName:                        crazyRhythms,
			},
			{
				createdAtTimeRelativeToCutoffTime: -20,
				updatedAtTimeRelativeToCutoffTime: -10,
				folderName:                        guyIncognito,
			},
		}, []db.MediaFile{
			{
				Id:        "5c214deb5b2dba739e0d6af56f61d1c7",
				Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 09 - crazy rhythms.mp3",
				CreatedAt: timeParse("2024-01-12T13:09:51"),
				UpdatedAt: timeParse("2024-01-12T13:29:51"),
				LibraryId: 1,
			},
			{
				Id:        "37141ae2932c8e06cc3716c3b9c55a48",
				Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 15 - i wanna sleep in your arms (live).mp3",
				CreatedAt: timeParse("2024-01-12T13:09:51"),
				UpdatedAt: timeParse("2024-01-12T13:29:51"),
				LibraryId: 1,
			},
			{
				Id:        "54c5999927b56e2887c3a5cfd21bdfbf",
				Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 13 - moscow nights (demo).mp3",
				CreatedAt: timeParse("2024-01-12T13:09:51"),
				UpdatedAt: timeParse("2024-01-12T13:29:51"),
				LibraryId: 1,
			},
		}),
		Entry("when no media files are updated after cutoff time and all are created earlier than updated", []mediaFolderTestData{
			{
				createdAtTimeRelativeToCutoffTime: -20,
				updatedAtTimeRelativeToCutoffTime: -10,
				folderName:                        crazyRhythms,
			},
			{
				createdAtTimeRelativeToCutoffTime: -20,
				updatedAtTimeRelativeToCutoffTime: -10,
				folderName:                        guyIncognito,
			},
		}, nil),

		Entry("when all media files are updated after cutoff time and all are created after updated", []mediaFolderTestData{
			{
				createdAtTimeRelativeToCutoffTime: 20,
				updatedAtTimeRelativeToCutoffTime: 10,
				folderName:                        crazyRhythms,
			},
			{
				createdAtTimeRelativeToCutoffTime: 20,
				updatedAtTimeRelativeToCutoffTime: 10,
				folderName:                        guyIncognito,
			},
		}, nil),
	}

	DescribeTableSubtree("media folder scenarios",
		func(testDataRecords []mediaFolderTestData, expectedMediaFiles []db.MediaFile) {
			var musicFoldersRepository *db.MusicFoldersRepository

			BeforeEach(func() {
				for _, td := range testDataRecords {
					Expect(
						testutils.UpdateMediaFileRecordCreatedAndUpdatedAt(
							testDbFullPath, td.createdAtTimeRelativeToCutoffTime, td.updatedAtTimeRelativeToCutoffTime, td.folderName, lastQueryEnded)).
						To(BeNil(), "Failed to update record %v", td.folderName)
				}
				musicFoldersRepository = &db.MusicFoldersRepository{SqliteHandler: sqliteHandler}
			})

			It("should return the expected media files", func() {
				Expect(musicFoldersRepository.UpdatedMediaFilesSinceDate(lastQueryEnded)).To(ConsistOf(expectedMediaFiles))
			})

		}, mediaFolderScenarios)
})

var _ = Describe("AllMediaFiles", func() {
	var musicFoldersRepository *db.MusicFoldersRepository

	var testDbFullPath = ""

	var err error
	var sqliteHandler *db.SQLiteHandler

	var expectedMediaFiles = []db.MediaFile{
		{
			Id:        "5c214deb5b2dba739e0d6af56f61d1c7",
			Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 09 - crazy rhythms.mp3",
			CreatedAt: timeParse("2024-01-12T13:09:51"),
			UpdatedAt: timeParse("2024-01-12T13:29:51"),
			LibraryId: 1,
		},
		{
			Id:        "37141ae2932c8e06cc3716c3b9c55a48",
			Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 15 - i wanna sleep in your arms (live).mp3",
			CreatedAt: timeParse("2024-01-12T13:09:51"),
			UpdatedAt: timeParse("2024-01-12T13:29:51"),
			LibraryId: 1,
		},
		{
			Id:        "54c5999927b56e2887c3a5cfd21bdfbf",
			Path:      "music/Crazy Rhythms/feelies, the - crazy rhythms - 13 - moscow nights (demo).mp3",
			CreatedAt: timeParse("2024-01-12T13:09:51"),
			UpdatedAt: timeParse("2024-01-12T13:29:51"),
			LibraryId: 1,
		},
		{
			Id:        "1fde8382304840139358a101b081db9c",
			Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 02 Can We Please Go Back To How It Was Before I messed Up-.mp3",
			CreatedAt: timeParse("2024-01-12T13:09:51"),
			UpdatedAt: timeParse("2024-01-12T13:39:51"),
			LibraryId: 1,
		},
		{
			Id:        "07bb5aec148d087b0192c538721d0627",
			Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 03 Seven.mp3",
			CreatedAt: timeParse("2024-01-12T13:09:51"),
			UpdatedAt: timeParse("2024-01-12T13:39:51"),
			LibraryId: 1,
		},
		{
			Id:        "6ea5a2baa32842109925f67b3151fb80",
			Path:      "music/Guy Incognito - Lovedrug/Guy Incognito - Lovedrug - 01 Lovedrug.mp3",
			CreatedAt: timeParse("2024-01-12T13:09:51"),
			UpdatedAt: timeParse("2024-01-12T13:39:51"),
			LibraryId: 1,
		},
	}

	BeforeEach(func() {
		By("Resetting and connecting to DB")

		testDbFullPath, err = testutils.SetupTestDb("fakenavidrome")
		Expect(err).To(BeNil(), "Error trying to setup DB")
		sqliteHandler = &db.SQLiteHandler{}
		Expect(sqliteHandler.ConnectSQLite(testDbFullPath)).To(BeNil(), "Failed to connect to sqlite")
		musicFoldersRepository = &db.MusicFoldersRepository{SqliteHandler: sqliteHandler}

		testDataRecords := []mediaFolderTestData{
			{
				createdAtTimeRelativeToCutoffTime: -10,
				updatedAtTimeRelativeToCutoffTime: 10,
				folderName:                        crazyRhythms,
			},
			{
				createdAtTimeRelativeToCutoffTime: -10,
				updatedAtTimeRelativeToCutoffTime: 20,
				folderName:                        guyIncognito,
			},
		}

		for _, td := range testDataRecords {
			Expect(
				testutils.UpdateMediaFileRecordCreatedAndUpdatedAt(
					testDbFullPath, td.createdAtTimeRelativeToCutoffTime,
					td.updatedAtTimeRelativeToCutoffTime, td.folderName, lastQueryEnded)).
				To(BeNil(), "Failed to update record %v", td.folderName)
		}
	})

	It("should return the expected media files", func() {
		Expect(musicFoldersRepository.AllMediaFiles()).To(ConsistOf(expectedMediaFiles))
	})
})

func timeParse(dateString string) time.Time {
	parsedTime, err := time.Parse("2006-01-02T15:04:05", dateString)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse time %v", dateString))
	}
	return parsedTime
}
