package zipper_test

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/apkatsikas/archiver/fileutil"
	"github.com/apkatsikas/archiver/fileutil/mocks"
	testutils "github.com/apkatsikas/archiver/tests/test-utils"
	"github.com/apkatsikas/archiver/zipper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("zipper when there are only files - integrated", func() {
	const folderName = "huey lewis - sports"

	var zipPath string
	var zipError error
	var expectedFilePath string

	BeforeEach(func() {
		dir, err := os.Getwd()
		Expect(err).To(BeNil(), "Got an error getting working directory")
		parentDir := filepath.Join(dir, "..")

		fixturesDir := filepath.Join(parentDir, "tests", "fixtures")
		folderPath := filepath.Join(fixturesDir, folderName)

		fileSystemOperator := fileutil.FileSystemOperator{}

		zipp := zipper.Zipper{FileSystemOperator: &fileSystemOperator}

		expectedFilePath = fmt.Sprintf("%v.zip", folderPath)

		err = testutils.RemoveFileIfExists(expectedFilePath)
		Expect(err).To(BeNil(), "Got an error trying to remove test file")
		Expect(testutils.FileExists(expectedFilePath)).To(BeFalse(), "Found unexpected test file")

		// If we want to parallelize these, we can use a random string when creating the zip
		// and delete it after the test run. As of now, they operate on the same file - could collide on
		// a parallel run
		zipPath, zipError = zipp.ZipFilesInFolder(folderPath)
	})

	It("should zip without error", func() {
		Expect(zipError).To(BeNil())
	})

	It("should return the expected file path", func() {
		Expect(zipPath).To(Equal(expectedFilePath))
	})

	It("should create a zip", func() {
		Expect(testutils.FileExists(expectedFilePath)).To(BeTrue())
	})

	var _ = Describe("zip contents", func() {
		var fileNames []string

		BeforeEach(func() {
			files, err := openZip(expectedFilePath)
			Expect(err).To(BeNil(), "Got an error trying to open zip file")
			defer files.Close()

			for _, file := range files.File {
				fileNames = append(fileNames, file.Name)
			}
		})

		It("should create a zip with the expected contents", func() {
			Expect(fileNames).To(ConsistOf([]string{
				filepath.Join(folderName, "hue lou.mp3"),
				filepath.Join(folderName, "cover.jpg")}))
		})
	})
})

var _ = Describe("zipper when there are only files - integrated - when user overrides file size and exceeds limit", func() {
	const folderName = "huey lewis - sports"

	var zipPath string
	var zipError error
	var expectedFilePath string

	BeforeEach(func() {
		dir, err := os.Getwd()
		Expect(err).To(BeNil(), "Got an error getting working directory")
		parentDir := filepath.Join(dir, "..")

		fixturesDir := filepath.Join(parentDir, "tests", "fixtures")
		folderPath := filepath.Join(fixturesDir, folderName)

		fileSystemOperator := fileutil.FileSystemOperator{}

		zipp := zipper.Zipper{FileSystemOperator: &fileSystemOperator}
		zipp.SetFileLimits(1, 0)

		expectedFilePath = fmt.Sprintf("%v.zip", folderPath)

		err = testutils.RemoveFileIfExists(expectedFilePath)
		Expect(err).To(BeNil(), "Got an error trying to remove test file")
		Expect(testutils.FileExists(expectedFilePath)).To(BeFalse(), "Found unexpected test file")

		// If we want to parallelize these, we can use a random string when creating the zip
		// and delete it after the test run. As of now, they operate on the same file - could collide on
		// a parallel run
		zipPath, zipError = zipp.ZipFilesInFolder(folderPath)
	})

	It("should return an error", func() {
		Expect(zipError).To(Not(BeNil()))
	})

	It("should return an empty file path", func() {
		Expect(zipPath).To(BeEmpty())
	})
})

var _ = Describe("zipper when there are only files - integrated - when user overrides file count and exceeds limit", func() {
	const folderName = "huey lewis - sports"

	var zipPath string
	var zipError error
	var expectedFilePath string

	BeforeEach(func() {
		dir, err := os.Getwd()
		Expect(err).To(BeNil(), "Got an error getting working directory")
		parentDir := filepath.Join(dir, "..")

		fixturesDir := filepath.Join(parentDir, "tests", "fixtures")
		folderPath := filepath.Join(fixturesDir, folderName)

		fileSystemOperator := fileutil.FileSystemOperator{}

		zipp := zipper.Zipper{FileSystemOperator: &fileSystemOperator}
		zipp.SetFileLimits(0, 1)

		expectedFilePath = fmt.Sprintf("%v.zip", folderPath)

		err = testutils.RemoveFileIfExists(expectedFilePath)
		Expect(err).To(BeNil(), "Got an error trying to remove test file")
		Expect(testutils.FileExists(expectedFilePath)).To(BeFalse(), "Found unexpected test file")

		// If we want to parallelize these, we can use a random string when creating the zip
		// and delete it after the test run. As of now, they operate on the same file - could collide on
		// a parallel run
		zipPath, zipError = zipp.ZipFilesInFolder(folderPath)
	})

	It("should return an error", func() {
		Expect(zipError).To(Not(BeNil()))
	})

	It("should return an empty file path", func() {
		Expect(zipPath).To(BeEmpty())
	})
})

var _ = Describe("zipper when there are a mix of folders and files", func() {
	const folderPath = "/path/to/music/album"
	const hueySong = "hue lou.mp3"
	const randomFolder = "randomfolder"
	const hueyCover = "cover.jpg"

	var expectedZipOutput = fmt.Sprintf("%v.zip", folderPath)

	var zipp *zipper.Zipper

	BeforeEach(func() {
		gt := GinkgoT()
		mockFileSystemOperator := mocks.NewIFileSystemOperator(gt)

		setupFileNamesFromPath(mockFileSystemOperator, folderPath, hueySong, hueyCover, randomFolder)
		setupZipFile(gt, mockFileSystemOperator)
		setupMockFile(gt, folderPath, hueySong, mockFileSystemOperator)
		setupMockFile(gt, folderPath, hueyCover, mockFileSystemOperator)
		setupMockDirFile(gt, folderPath, randomFolder, mockFileSystemOperator)

		zipp = &zipper.Zipper{FileSystemOperator: mockFileSystemOperator}
	})

	It("should successfully create a zip", func() {
		Expect(zipp.ZipFilesInFolder(folderPath)).To(Equal(expectedZipOutput))
	})
})

var _ = Describe("zipper when there are no files or folders", func() {
	const folderPath = "/path/to/music/album"

	var zipp *zipper.Zipper
	var zipError error
	var zipFullPath string

	BeforeEach(func() {
		gt := GinkgoT()
		mockFileSystemOperator := mocks.NewIFileSystemOperator(gt)

		setupFileNamesFromPath(mockFileSystemOperator, folderPath)

		zipp = &zipper.Zipper{FileSystemOperator: mockFileSystemOperator}
		zipFullPath, zipError = zipp.ZipFilesInFolder(folderPath)
	})

	It("should return an error", func() {
		Expect(zipError).To(Not(BeNil()))
	})

	It("should return an empty path", func() {
		Expect(zipFullPath).To(BeEmpty())
	})
})

var _ = Describe("zipper when there are only folders", func() {
	const folderPath = "/path/to/music/album"
	const innerFolder = "folder"
	const anotherInnerFolder = "anotherFolder"

	var zipp *zipper.Zipper

	BeforeEach(func() {
		gt := GinkgoT()
		mockFileSystemOperator := mocks.NewIFileSystemOperator(gt)

		setupFileNamesFromPath(mockFileSystemOperator, folderPath, innerFolder, anotherInnerFolder)
		setupMockDirFile(gt, folderPath, innerFolder, mockFileSystemOperator)
		setupMockDirFile(gt, folderPath, anotherInnerFolder, mockFileSystemOperator)

		zipp = &zipper.Zipper{FileSystemOperator: mockFileSystemOperator}
	})

	It("should return an empty path", func() {
		Expect(zipp.ZipFilesInFolder(folderPath)).To(BeEmpty())
	})
})

var _ = Describe("zipper when a file in the folder is too large", func() {
	const folderPath = "/path/to/music/album"
	const hueySong = "hue lou.mp3"
	const hueyCover = "cover.jpg"

	var zipError error
	var zipFullPath string

	BeforeEach(func() {
		gt := GinkgoT()
		mockFileSystemOperator := mocks.NewIFileSystemOperator(gt)

		setupFileNamesFromPath(mockFileSystemOperator, folderPath, hueySong, hueyCover)
		setupMockFileLarge(gt, folderPath, hueySong, mockFileSystemOperator)

		zipp := &zipper.Zipper{FileSystemOperator: mockFileSystemOperator}

		zipFullPath, zipError = zipp.ZipFilesInFolder(folderPath)
	})

	It("should return an error", func() {
		Expect(zipError).To(Not(BeNil()))
	})

	It("should return an empty path", func() {
		Expect(zipFullPath).To(BeEmpty())
	})
})

var _ = Describe("zipper when there are too many files", func() {
	const folderPath = "/path/to/music/album"

	var zipError error
	var zipFullPath string

	BeforeEach(func() {
		gt := GinkgoT()
		files := generateRandomStringSlice(151)
		mockFileSystemOperator := mocks.NewIFileSystemOperator(gt)

		setupFileNamesFromPath(mockFileSystemOperator, folderPath, files...)

		zipp := &zipper.Zipper{FileSystemOperator: mockFileSystemOperator}
		zipFullPath, zipError = zipp.ZipFilesInFolder(folderPath)
	})

	It("should return an error", func() {
		Expect(zipError).To(Not(BeNil()))
	})

	It("should return an empty path", func() {
		Expect(zipFullPath).To(BeEmpty())
	})
})

var _ = Describe("zipper when it fails to open a file", func() {
	const folderPath = "/path/to/music/album"
	const hueySong = "hue lou.mp3"
	const hueyCover = "cover.jpg"

	var expectedZipToBeDeleted = fmt.Sprintf("%v.zip", folderPath)

	var openFileError = fmt.Errorf("oh no! failed to open file!")
	var zipError error
	var zipFullPath string

	BeforeEach(func() {
		gt := GinkgoT()
		mockFileSystemOperator := mocks.NewIFileSystemOperator(gt)

		setupFileNamesFromPath(mockFileSystemOperator, folderPath, hueySong, hueyCover)
		setupZipFile(gt, mockFileSystemOperator)

		fileInfo := mocks.NewIArchiveFileInfo(gt)
		fileInfo.EXPECT().IsDir().Return(false).Once()
		filePath := path.Join(folderPath, hueySong)

		setupFileInfoExpectations(fileInfo, hueySong)

		mockFileSystemOperator.EXPECT().GetInfo(filePath).Return(fileInfo, nil).Once()
		mockFileSystemOperator.EXPECT().OpenFile(filePath).Return(nil, openFileError).Once()

		mockFileSystemOperator.EXPECT().DeleteFile(expectedZipToBeDeleted).Return(nil).Once()

		zipp := &zipper.Zipper{FileSystemOperator: mockFileSystemOperator}
		zipFullPath, zipError = zipp.ZipFilesInFolder(folderPath)
	})

	It("should return an error", func() {
		Expect(zipError).To(Not(BeNil()))
	})

	It("should return an empty path", func() {
		Expect(zipFullPath).To(BeEmpty())
	})
})

var _ = Describe("zipper when there is a zip in the folder", func() {
	const folderPath = "/path/to/music/album"
	const hueyZip = "huey.zip"

	var zipp *zipper.Zipper
	var zipError error
	var zipFullPath string

	BeforeEach(func() {
		gt := GinkgoT()
		mockFileSystemOperator := mocks.NewIFileSystemOperator(gt)

		mockFileSystemOperator.EXPECT().DeleteFile(mock.AnythingOfType("string")).Return(nil).Once()

		setupFileNamesFromPath(mockFileSystemOperator, folderPath, hueyZip)

		zipp = &zipper.Zipper{FileSystemOperator: mockFileSystemOperator}
		zipFullPath, zipError = zipp.ZipFilesInFolder(folderPath)
	})

	It("should return an error", func() {
		Expect(zipError).To(Not(BeNil()))
	})

	It("should return an empty path", func() {
		Expect(zipFullPath).To(BeEmpty())
	})
})

func setupMockFile(
	gt FullGinkgoTInterface, folderName string, fileName string, fsOperator *mocks.IFileSystemOperator) {
	setupMockFileHelper(gt, folderName, fileName, fsOperator, false)
}

func setupMockFileLarge(
	gt FullGinkgoTInterface, folderName string, fileName string, fsOperator *mocks.IFileSystemOperator) {
	setupMockFileHelper(gt, folderName, fileName, fsOperator, true)
}

func setupMockFileHelper(gt FullGinkgoTInterface, folderName string, fileName string, fsOperator *mocks.IFileSystemOperator, largeFile bool) {
	archiveFileInfo := mocks.NewIArchiveFileInfo(gt)
	archiveFileInfo.EXPECT().IsDir().Return(false).Once()

	if largeFile {
		const fileSizeLimit = 524288010 // 500mb plus 10 bytes

		archiveFileInfo.EXPECT().Size().Return(fileSizeLimit).Once()

		filePath := path.Join(folderName, fileName)
		fsOperator.EXPECT().GetInfo(filePath).Return(archiveFileInfo, nil).Once()
		// here we try to delete the zip file if it exists - but it won't
		fsOperator.EXPECT().DeleteFile(mock.AnythingOfType("string")).Return(fmt.Errorf("file does not exist")).Once()
	} else {
		setupFileInfoExpectations(archiveFileInfo, fileName)

		filePath := path.Join(folderName, fileName)
		fsOperator.EXPECT().GetInfo(filePath).Return(archiveFileInfo, nil).Once()
		fsOperator.EXPECT().OpenFile(filePath).Return(createMockFile(gt), nil).Once()
	}
}

func setupFileInfoExpectations(archiveFileInfo *mocks.IArchiveFileInfo, fileName string) {
	archiveFileInfo.EXPECT().ModTime().Return(time.Now().Add(time.Duration(-60) * time.Minute)).Once()
	archiveFileInfo.EXPECT().Mode().Return(os.FileMode(int(0444))).Once()
	archiveFileInfo.EXPECT().Name().Return(fileName).Once()
	archiveFileInfo.EXPECT().Size().Return(666).Times(2)
}

func setupMockDirFile(gt FullGinkgoTInterface, folderName string, fileName string, fsOperator *mocks.IFileSystemOperator) {
	archiveFileInfo := mocks.NewIArchiveFileInfo(gt)
	archiveFileInfo.EXPECT().IsDir().Return(true).Once()

	filePath := path.Join(folderName, fileName)
	fsOperator.EXPECT().GetInfo(filePath).Return(archiveFileInfo, nil).Once()
}

func createMockFile(gt FullGinkgoTInterface) *mocks.IArchiveFile {
	archiveContentFile := mocks.NewIArchiveFile(gt)
	archiveContentFile.EXPECT().Close().Return(nil).Once()
	archiveContentFile.EXPECT().Read(mock.AnythingOfType("[]uint8")).Return(1, nil).Once()
	archiveContentFile.EXPECT().Read(mock.AnythingOfType("[]uint8")).Return(1, io.EOF).Once()
	return archiveContentFile
}

func setupFileNamesFromPath(fsOperator *mocks.IFileSystemOperator, pathArgument string, fileNamesToReturn ...string) {
	fsOperator.EXPECT().FileNamesFromPath(pathArgument).Return(fileNamesToReturn, nil)
}

func setupZipFile(gt FullGinkgoTInterface, fsOperator *mocks.IFileSystemOperator) {
	zipFile := mocks.NewIArchiveFile(gt)

	fsOperator.EXPECT().CreateFile(mock.AnythingOfType("string")).Return(zipFile, nil).Once()

	zipFile.EXPECT().Write(mock.AnythingOfType("[]uint8")).Return(1000, nil).Once()
	zipFile.EXPECT().Close().Return(nil).Once()
}

func openZip(filePath string) (*zip.ReadCloser, error) {
	return zip.OpenReader(filePath)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b) + ".mp3"
}

func generateRandomStringSlice(count int) []string {
	var result []string
	for i := 0; i < count; i++ {
		result = append(result, generateRandomString(10))
	}
	return result
}
