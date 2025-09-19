package runner_test

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/apkatsikas/archiver/db"
	"github.com/apkatsikas/archiver/fileutil"
	fsoMocks "github.com/apkatsikas/archiver/fileutil/mocks"
	"github.com/apkatsikas/archiver/filter"
	"github.com/apkatsikas/archiver/runner"
	storageMocks "github.com/apkatsikas/archiver/storage-client/mocks"

	testutils "github.com/apkatsikas/archiver/tests/test-utils"
	"github.com/apkatsikas/archiver/zipper"
)

const (
	fakeNavidromeDb  = "fakenavidromerunner"
	fakeArchiveRunDb = "fakearchiverun"
	navidromeBackup  = "navidrome-backup.sqlite"
)

type runTestData struct {
	runTypeTest
	hueyTimeDiff timeDiff
	mc5TimeDiff  timeDiff
	priorRun     bool
}

type timeDiff struct {
	createdDiff int
	updatedDiff int
}

type runTypeTest int64

type artistPathZips struct {
	hueyPathZip string
	mc5PathZip  string
}

const (
	Replace runTypeTest = 0
	Upload  runTypeTest = 1
	NoOp    runTypeTest = 2
	Both    runTypeTest = 3
)

var _ = DescribeTableSubtree("Runner when there is a prior run",
	func(testData runTestData) {
		var runner = &runner.Runner{}
		var err error
		var artistPathZips *artistPathZips

		BeforeEach(func() {
			artistPathZips = setup(runner, testData)
			err = runner.RunScheduled()
		})

		It("Runs without error", func() {
			Expect(err).To(BeNil())
		})

		It("Does not leave behind a zip file", func() {
			Expect(testutils.FileExists(artistPathZips.hueyPathZip)).To(
				BeFalse(), "Did not expect file to exist %v", artistPathZips.hueyPathZip)
			Expect(testutils.FileExists(artistPathZips.mc5PathZip)).To(
				BeFalse(), "Did not expect file to exist %v", artistPathZips.mc5PathZip)
		})
	},
	Entry("Upload only", runTestData{
		hueyTimeDiff: timeDiff{createdDiff: 10, updatedDiff: 10},
		mc5TimeDiff:  timeDiff{createdDiff: 10, updatedDiff: 10},
		runTypeTest:  Upload,
		priorRun:     true,
	}),
	Entry("Replace only", runTestData{
		hueyTimeDiff: timeDiff{createdDiff: -10, updatedDiff: 10},
		mc5TimeDiff:  timeDiff{createdDiff: -10, updatedDiff: 10},
		runTypeTest:  Replace,
		priorRun:     true,
	}),
	Entry("Upload and Replace", runTestData{
		hueyTimeDiff: timeDiff{createdDiff: 10, updatedDiff: 10},
		mc5TimeDiff:  timeDiff{createdDiff: -10, updatedDiff: 10},
		runTypeTest:  Both,
		priorRun:     true,
	}),
	Entry("NoOp", runTestData{
		hueyTimeDiff: timeDiff{createdDiff: -10, updatedDiff: -10},
		mc5TimeDiff:  timeDiff{createdDiff: -10, updatedDiff: -10},
		runTypeTest:  NoOp,
		priorRun:     true,
	}),
)

var _ = Describe("Runner when there is not a prior run", func() {
	var runner = &runner.Runner{}
	var err error
	var artistPathZips *artistPathZips

	BeforeEach(func() {
		artistPathZips = setup(runner, runTestData{
			hueyTimeDiff: timeDiff{createdDiff: 10, updatedDiff: 10},
			mc5TimeDiff:  timeDiff{createdDiff: 10, updatedDiff: 10},
			runTypeTest:  NoOp,
			priorRun:     false,
		})
		err = runner.RunScheduled()
	})

	It("Runs without error", func() {
		Expect(err).To(BeNil())
	})

	It("Does not create any files", func() {
		Expect(testutils.FileExists(artistPathZips.hueyPathZip)).To(
			BeFalse(), "Did not expect file to exist %v", artistPathZips.hueyPathZip)
		Expect(testutils.FileExists(artistPathZips.mc5PathZip)).To(
			BeFalse(), "Did not expect file to exist %v", artistPathZips.mc5PathZip)
	})
})

func setup(runner *runner.Runner, testData runTestData) *artistPathZips {
	var lastRun = time.Date(2024, time.January, 12, 13, 19, 51, 0, time.UTC)
	const hueyPath = "tests/fixtures/huey lewis - sports"
	const mc5Path = "tests/fixtures/mc5 - back in the usa"

	GinkgoHelper()
	gt := GinkgoT()

	By("Setting up FileSystemOperator")
	fso := &fileutil.FileSystemOperator{}

	runner.FileSystemOperator = fso
	By("Get working directory")
	testDir, err := os.Getwd()
	Expect(err).To(BeNil(), "Got an error getting working directory")
	archiverDir := filepath.Join(testDir, "..")

	By("Deleting prior artifacts")
	hueyPathZip := deletePriorArtifacts(hueyPath, archiverDir)
	mc5PathZip := deletePriorArtifacts(mc5Path, archiverDir)

	By("Delete prior backup")
	err = testutils.RemoveFileIfExists(navidromeBackup)
	Expect(err).To(BeNil(), "Got an error trying to remove backup file")

	By("Setting up MusicFoldersRepository")
	fakeNavidromeDbFullPath := setupNavidromeRepositories(runner)

	By("Setting up ArchiveRunRepository")
	var fakeArchiveRunDbFullPath = ""
	if testData.priorRun {
		fakeArchiveRunDbFullPath, err = testutils.SetupTestDb(fakeArchiveRunDb)
		Expect(err).To(BeNil(), "Error trying to setup fakenavidrome DB")

	} else {
		dir, err := os.Getwd()
		Expect(err).To(BeNil(), "Got an error getting working directory")
		// If we want to parallelize these, we can use a random string when creating the DB
		// and delete it after the test run. As of now, they operate on the same file - could collide on
		// a parallel run
		fixturesPath := "tests/fixtures"
		testDbPath := fmt.Sprintf("%v/%v.db", fixturesPath, fakeArchiveRunDb)
		testDbFullPath := filepath.Join(dir, "..", testDbPath)

		Expect(testutils.RemoveFileIfExists(testDbFullPath)).To(BeNil(), "Got an error trying to remove %v", testDbFullPath)
	}
	runner.ArchiveRunRepository = &db.ArchiveRunRepository{SqliteHandler: &db.SQLiteHandler{}}
	Expect(runner.ArchiveRunRepository.SqliteHandler.ConnectSQLite(fakeArchiveRunDbFullPath)).To(
		BeNil(), "Failed to connect to sqlite for ArchiveRunRepository")

	By("Setting up Zipper")
	runner.Zipper = &zipper.Zipper{FileSystemOperator: fso}

	By("Setting up FilterService")
	runner.FilterService = &filter.FilterService{}

	By("Setting up StorageClient")
	mockStorageClient := storageMocks.NewIStorageClient(gt)

	if testData.priorRun {
		switch td := testData.runTypeTest; td {
		case Upload:
			By("Expecting to upload 2 new files")
			mockStorageClient.EXPECT().ReplaceFile(navidromeBackup, navidromeBackup).Return(nil).Once()
			mockStorageClient.EXPECT().UploadNewFile(
				hueyPathZip, "huey lewis - sports5c214deb5b2dba739e0d6af56f61d1c7.zip").Return(nil).Once()
			mockStorageClient.EXPECT().UploadNewFile(
				mc5PathZip, "mc5 - back in the usa6ea5a2baa32842109925f67b3151fb80.zip").Return(nil).Once()
		case Replace:
			By("Expecting to replace 2 files")
			mockStorageClient.EXPECT().ReplaceFile(navidromeBackup, navidromeBackup).Return(nil).Once()
			mockStorageClient.EXPECT().ReplaceFile(
				hueyPathZip, "huey lewis - sports5c214deb5b2dba739e0d6af56f61d1c7.zip").Return(nil).Once()
			mockStorageClient.EXPECT().ReplaceFile(
				mc5PathZip, "mc5 - back in the usa6ea5a2baa32842109925f67b3151fb80.zip").Return(nil).Once()
		case Both:
			By("Expecting to upload 1 new file and replace 1 file")
			mockStorageClient.EXPECT().ReplaceFile(navidromeBackup, navidromeBackup).Return(nil).Once()
			mockStorageClient.EXPECT().UploadNewFile(
				hueyPathZip, "huey lewis - sports5c214deb5b2dba739e0d6af56f61d1c7.zip").Return(nil).Once()
			mockStorageClient.EXPECT().ReplaceFile(
				mc5PathZip, "mc5 - back in the usa6ea5a2baa32842109925f67b3151fb80.zip").Return(nil).Once()
		case NoOp:
			By("Expecting to do nothing with storage")
		}
		runner.StorageClient = mockStorageClient
	}

	By("Setting up media file record create and update dates")
	hueySongPath := "tests/fixtures/huey lewis - sports/hue lou.mp3"
	Expect(
		testutils.UpdateMediaFileRecordCreatedAndUpdatedAt(
			fakeNavidromeDbFullPath, testData.hueyTimeDiff.createdDiff, testData.hueyTimeDiff.updatedDiff, hueySongPath, lastRun)).
		To(BeNil(), "Failed to update record %v", hueyPath)
	mc5SongPath := "tests/fixtures/mc5 - back in the usa/tutti fruitti.mp3"
	Expect(
		testutils.UpdateMediaFileRecordCreatedAndUpdatedAt(
			fakeNavidromeDbFullPath, testData.mc5TimeDiff.createdDiff, testData.mc5TimeDiff.updatedDiff, mc5SongPath, lastRun)).
		To(BeNil(), "Failed to update record %v", mc5Path)

	By("Setting up last run date")
	if testData.priorRun {
		Expect(testutils.UpdateLastRunRecord(fakeArchiveRunDbFullPath, lastRun)).To(BeNil(), "Failed to update archive_run")
	}
	By("Setting up media file record paths")
	testutils.UpdateMediaFilePathPrefix(fakeNavidromeDbFullPath, archiverDir, hueySongPath)
	testutils.UpdateMediaFilePathPrefix(fakeNavidromeDbFullPath, archiverDir, mc5SongPath)
	return &artistPathZips{hueyPathZip, mc5PathZip}
}

var _ = Describe("BuildLedger", func() {
	var runn *runner.Runner
	const filePath = "/path/to/file.json"
	var expectedBytes = []byte{
		0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x22, 0x2f, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x73, 0x2f, 0x66, 0x69, 0x78, 0x74, 0x75, 0x72, 0x65, 0x73, 0x2f, 0x68, 0x75, 0x65, 0x79, 0x20, 0x6c, 0x65, 0x77, 0x69, 0x73, 0x20, 0x2d, 0x20, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x22, 0x3a, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x30, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x49, 0x64, 0x22, 0x3a, 0x20, 0x22, 0x35, 0x63, 0x32, 0x31, 0x34, 0x64, 0x65, 0x62, 0x35, 0x62, 0x32, 0x64, 0x62, 0x61, 0x37, 0x33, 0x39, 0x65, 0x30, 0x64, 0x36, 0x61, 0x66, 0x35, 0x36, 0x66, 0x36, 0x31, 0x64, 0x31, 0x63, 0x37, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x42, 0x61, 0x73, 0x65, 0x50, 0x61, 0x74, 0x68, 0x22, 0x3a, 0x20, 0x22, 0x68, 0x75, 0x65, 0x79, 0x20, 0x6c, 0x65, 0x77, 0x69, 0x73, 0x20, 0x2d, 0x20, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x22, 0x2f, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x73, 0x2f, 0x66, 0x69, 0x78, 0x74, 0x75, 0x72, 0x65, 0x73, 0x2f, 0x6d, 0x63, 0x35, 0x20, 0x2d, 0x20, 0x62, 0x61, 0x63, 0x6b, 0x20, 0x69, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x75, 0x73, 0x61, 0x22, 0x3a, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x30, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x49, 0x64, 0x22, 0x3a, 0x20, 0x22, 0x36, 0x65, 0x61, 0x35, 0x61, 0x32, 0x62, 0x61, 0x61, 0x33, 0x32, 0x38, 0x34, 0x32, 0x31, 0x30, 0x39, 0x39, 0x32, 0x35, 0x66, 0x36, 0x37, 0x62, 0x33, 0x31, 0x35, 0x31, 0x66, 0x62, 0x38, 0x30, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x42, 0x61, 0x73, 0x65, 0x50, 0x61, 0x74, 0x68, 0x22, 0x3a, 0x20, 0x22, 0x6d, 0x63, 0x35, 0x20, 0x2d, 0x20, 0x62, 0x61, 0x63, 0x6b, 0x20, 0x69, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x75, 0x73, 0x61, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x7d, 0xa, 0x7d,
	}

	BeforeEach(func() {
		gt := GinkgoT()
		By("Setting up FileSystemOperator")
		mockFileSystemOperator := fsoMocks.NewIFileSystemOperator(gt)
		mockFileSystemOperator.EXPECT().WriteNewFile(filePath, expectedBytes).Return(nil).Once()

		By("Setting up Runner")
		runn = &runner.Runner{FileSystemOperator: mockFileSystemOperator}

		By("Setting up MusicFoldersRepository")
		setupNavidromeRepositories(runn)
	})

	It("Runs without error", func() {
		Expect(runn.BuildLedger(filePath)).To(BeNil())
	})
})

var _ = Describe("RunBatch", func() {
	const hueyPath = "tests/fixtures/huey lewis - sports"
	const mc5Path = "tests/fixtures/mc5 - back in the usa"

	var runn *runner.Runner

	var archiveDir = ""

	BeforeEach(func() {
		gt := GinkgoT()
		By("Setting up FileSystemOperator")
		fso := &fileutil.FileSystemOperator{}

		By("Setting up Runner")
		runn = &runner.Runner{FileSystemOperator: fso, Zipper: &zipper.Zipper{FileSystemOperator: fso}}

		By("Get working directory")
		testDir, err := os.Getwd()
		Expect(err).To(BeNil(), "Got an error getting working directory")
		archiveDir = filepath.Join(testDir, "..")

		By("Deleting prior artifacts")
		deletePriorArtifacts(hueyPath, archiveDir)
		deletePriorArtifacts(mc5Path, archiveDir)

		By("Setting up Storage")
		storage := storageMocks.NewIStorageClient(gt)
		storage.EXPECT().UploadNewFile(
			"../tests/fixtures/huey lewis - sports.zip", "huey lewis - sports37141ae2932c8e06cc3716c3b9c55a48.zip").Return(nil).Once()
		storage.EXPECT().UploadNewFile(
			"../tests/fixtures/mc5 - back in the usa.zip", "mc5 - back in the usa07bb5aec148d087b0192c538721d0627.zip").Return(nil).Once()
		runn.StorageClient = storage
	})

	It("Runs without error", func() {
		Expect(runn.RunBatch(filepath.Join(archiveDir, "tests/fixtures/testledger.json"))).To(BeNil())
	})
})

func setupNavidromeRepositories(runner *runner.Runner) string {
	GinkgoHelper()
	fakeNavidromeDbFullPath, err := testutils.SetupTestDb(fakeNavidromeDb)
	Expect(err).To(BeNil(), "Error trying to setup fakenavidrome DB")
	sqlLiteHandler := &db.SQLiteHandler{}
	runner.MusicFoldersRepository = &db.MusicFoldersRepository{SqliteHandler: sqlLiteHandler}
	runner.AdminRepository = &db.AdminRepository{SqliteHandler: sqlLiteHandler}
	runner.LibraryRepository = &db.LibraryRepository{SqliteHandler: sqlLiteHandler}
	Expect(runner.MusicFoldersRepository.SqliteHandler.ConnectSQLite(fakeNavidromeDbFullPath)).To(
		BeNil(), "Failed to connect to sqlite for MusicFoldersRepository")
	return fakeNavidromeDbFullPath
}

func deletePriorArtifacts(path string, archiveDir string) string {
	GinkgoHelper()
	pathZip := filepath.Join(archiveDir, fmt.Sprintf("%v.zip", path))
	err := testutils.RemoveFileIfExists(pathZip)
	Expect(err).To(BeNil(), "Got an error trying to remove test file")
	Expect(testutils.FileExists(pathZip)).To(BeFalse(), "Found unexpected test file")
	return pathZip
}
