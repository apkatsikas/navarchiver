package zipper

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/apkatsikas/archiver/fileutil"
)

const (
	defaultFileCountLimit uint = 150
	defaultFileSizeLimit  uint = 524288000 // 500mb
)

type Zipper struct {
	FileSystemOperator fileutil.IFileSystemOperator
	builder            *zipBuilder
	fileCountLimit     uint
	fileSizeLimit      uint
}

type zipBuilder struct {
	createdZip bool
	zip        fileutil.IArchiveFile
	writer     *zip.Writer
	fullPath   string
}

func (zb *zipBuilder) closeBuilder() {
	if zb.createdZip {
		if zb.writer != nil {
			zb.writer.Close()
		}
		if zb.zip != nil {
			zb.zip.Close()
		}
	}
}

func (z *Zipper) SetFileLimits(fileSizeLimit, fileCountLimit uint) {
	z.fileSizeLimit = fileSizeLimit
	z.fileCountLimit = fileCountLimit
}

func (z *Zipper) ZipFilesInFolder(folderPath string) (string, error) {
	if z.fileCountLimit == 0 {
		z.fileCountLimit = defaultFileCountLimit
	}
	if z.fileSizeLimit == 0 {
		z.fileSizeLimit = defaultFileSizeLimit
	}
	z.builder = &zipBuilder{}
	fileNames, err := z.FileSystemOperator.FileNamesFromPath(folderPath)

	if len(fileNames) == 0 {
		return "", fmt.Errorf("0 files found at path %v", folderPath)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get file names from path: %v", err)
	}

	fileCount := len(fileNames)
	if fileCount > int(z.fileCountLimit) {
		return "", fmt.Errorf(
			"got %v files in folder %v, limit is %v", fileCount, folderPath, z.fileCountLimit)
	}

	for _, fileName := range fileNames {
		if strings.HasSuffix(fileName, ".zip") {
			return "", z.closeAndError(fmt.Errorf("folder to zip contained a zip: %v", fileName))
		}
		err = z.addToZip(folderPath, fileName)

		if err != nil {
			return "", z.closeAndError(err)
		}
	}
	z.builder.closeBuilder()
	return z.builder.fullPath, nil
}

func (z *Zipper) closeAndError(err error) error {
	z.builder.closeBuilder()
	deleteErr := z.FileSystemOperator.DeleteFile(z.builder.fullPath)
	if deleteErr != nil {
		log.Printf(
			"ERROR: Got an error when trying to delete zip during closeAndError: %v", deleteErr)
	}
	return err
}

func (z *Zipper) createZip(folderPath string) (fileutil.IArchiveFile, error) {
	z.builder.fullPath = z.zipFullPathName(folderPath)

	zipFile, err := z.FileSystemOperator.CreateFile(z.builder.fullPath)
	if err != nil {
		return nil, err
	}
	return zipFile, nil
}

func (z *Zipper) zipFullPathName(folderPath string) string {
	zipFileName := fmt.Sprintf("%v.zip", filepath.Base(folderPath))
	return filepath.Join(folderPath, "..", zipFileName)
}

func (z *Zipper) addToZip(folderPath string, fileName string) error {
	joinedPath := filepath.Join(folderPath, fileName)

	pathInfo, err := z.FileSystemOperator.GetInfo(joinedPath)

	if err != nil {
		return err
	}

	if !pathInfo.IsDir() {

		fileSize := pathInfo.Size()
		if fileSize > int64(z.fileSizeLimit) {
			return fmt.Errorf("file %v size is %v. Limit is %v", joinedPath,
				fileutil.FileSize(fileSize).String(), fileutil.FileSize(z.fileSizeLimit).String())
		}

		if !z.builder.createdZip {
			zipFile, err := z.createZip(folderPath)
			if err != nil {
				return err
			}
			z.builder.createdZip = true

			z.builder.zip = zipFile
			z.builder.writer = zip.NewWriter(bufio.NewWriter(zipFile))
		}

		fileInfoHeader, err := zip.FileInfoHeader(pathInfo)
		if err != nil {
			return err
		}

		fileInfoHeader.Method = zip.Deflate

		fileInfoHeader.Name, err = filepath.Rel(filepath.Dir(folderPath), joinedPath)
		if err != nil {
			return err
		}

		headerWriter, err := z.builder.writer.CreateHeader(fileInfoHeader)
		if err != nil {
			return err
		}

		fileToCopy, err := z.FileSystemOperator.OpenFile(joinedPath)
		if err != nil {
			return err
		}
		defer fileToCopy.Close()

		_, err = io.Copy(headerWriter, fileToCopy)

		if err != nil {
			return err
		}
	}
	return nil
}
