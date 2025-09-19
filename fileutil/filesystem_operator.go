package fileutil

import (
	"io/fs"
	"os"
	"time"
)

type FileSystemOperator struct {
}

func (fso *FileSystemOperator) FileNamesFromPath(folderPath string) ([]string, error) {
	filepathFolder, err := fso.OpenFile(folderPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := filepathFolder.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	return fileNames, nil
}

func (fso *FileSystemOperator) GetInfo(path string) (fs.FileInfo, error) {
	stattedFile, err := os.Stat(path)

	if err != nil {
		return nil, err
	}
	return stattedFile, nil
}

func (fso *FileSystemOperator) CreateFile(filePath string) (IArchiveFile, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (fso *FileSystemOperator) OpenFile(filePath string) (IArchiveFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (fso *FileSystemOperator) DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func (fso *FileSystemOperator) WriteNewFile(filePath string, data []byte) error {
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (fso *FileSystemOperator) ReadFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//go:generate mockery --name IFileSystemOperator
type IFileSystemOperator interface {
	FileNamesFromPath(folderPath string) ([]string, error)
	GetInfo(path string) (fs.FileInfo, error)
	CreateFile(zipFileFullPath string) (IArchiveFile, error)
	OpenFile(filePath string) (IArchiveFile, error)
	DeleteFile(filePath string) error
	WriteNewFile(filePath string, data []byte) error
	ReadFile(filePath string) ([]byte, error)
}

//go:generate mockery --name IArchiveFile
type IArchiveFile interface {
	Write(p []byte) (n int, err error)
	Close() error
	Read(p []byte) (n int, err error)
	Readdirnames(n int) (names []string, err error)
}

//go:generate mockery --name IArchiveFileInfo
type IArchiveFileInfo interface {
	Name() string
	Size() int64
	Mode() fs.FileMode
	ModTime() time.Time
	IsDir() bool
	Sys() any
}
