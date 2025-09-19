package runner

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/apkatsikas/archiver/db"
	"github.com/apkatsikas/archiver/fileutil"
	"github.com/apkatsikas/archiver/filter"
	storageclient "github.com/apkatsikas/archiver/storage-client"
	"github.com/apkatsikas/archiver/zipper"
	_ "github.com/mattn/go-sqlite3"
)

type Runner struct {
	*filter.FilterService
	StorageClient storageclient.IStorageClient
	*db.MusicFoldersRepository
	*db.ArchiveRunRepository
	*db.AdminRepository
	*db.LibraryRepository
	*zipper.Zipper
	FileSystemOperator fileutil.IFileSystemOperator
}

const navidromeBackupDB = "navidrome-backup.sqlite"

func (r *Runner) RunScheduled() error {
	err := r.ArchiveRunRepository.CreateTable()
	if err != nil {
		return fmt.Errorf("failed to get CreateTable archive run: %v", err)
	}

	lastRun, err := r.ArchiveRunRepository.LastRun()
	if err != nil {
		return fmt.Errorf("failed to get last archive run: %v", err)
	}

	if lastRun == nil {
		err := r.ArchiveRunRepository.UpdateLastRun(time.Now())
		if err != nil {
			return fmt.Errorf("failed to update last archive run: %v", err)
		}
		return nil
	}

	newMediaFiles, err := r.MusicFoldersRepository.NewMediaFilesSinceDate(lastRun.LastRun)
	if err != nil {
		return fmt.Errorf("failed to get media files: %v", err)
	}
	absoluteNewMediaFiles, err := r.absoluteMediaFiles(newMediaFiles)
	if err != nil {
		return fmt.Errorf("failed to get absolute media files: %v", err)
	}

	updatedMediaFiles, err := r.MusicFoldersRepository.UpdatedMediaFilesSinceDate(lastRun.LastRun)
	if err != nil {
		return fmt.Errorf("failed to get updated media files: %v", err)
	}
	absoluteUpdatedMediaFiles, err := r.absoluteMediaFiles(updatedMediaFiles)
	if err != nil {
		return fmt.Errorf("failed to get absolute updated media files: %v", err)
	}

	identifiedPaths := r.FilterService.UpdatedAndNewIdentifiedPaths(
		absoluteNewMediaFiles, absoluteUpdatedMediaFiles)

	zips, err := r.zipFiles(identifiedPaths)
	if err != nil {
		return fmt.Errorf("failed to zip files: %v", err)
	}

	err = r.handleStorage(zips)
	if err != nil {
		return fmt.Errorf("failed to handle storage: %v", err)
	}

	if len(identifiedPaths) > 0 {
		if err := r.FileSystemOperator.DeleteFile(navidromeBackupDB); err != nil {
			log.Printf("failed to delete navidrome backup DB: %v", err)
		}
		if err := r.AdminRepository.CreateBackup(navidromeBackupDB); err != nil {
			return fmt.Errorf("failed to vacuum: %v", err)
		}

		if err := r.StorageClient.ReplaceFile(navidromeBackupDB, navidromeBackupDB); err != nil {
			if err := r.StorageClient.UploadNewFile(navidromeBackupDB, navidromeBackupDB); err != nil {
				return fmt.Errorf("failed to navidrome backup DB to storage: %v", err)
			}
		}
	}

	err = r.ArchiveRunRepository.UpdateLastRun(time.Now())
	if err != nil {
		return fmt.Errorf("failed to update last archive run: %v", err)
	}

	return nil
}

func (r *Runner) RunBatch(jsonPath string) error {
	data, err := r.FileSystemOperator.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read json: %v", err)
	}

	var identifiedPaths filter.IdentifiedPaths
	err = json.Unmarshal(data, &identifiedPaths)
	if err != nil {
		return fmt.Errorf("got an error trying to unmarshal: %v", err)
	}

	zips, err := r.zipFiles(identifiedPaths)
	if err != nil {
		return fmt.Errorf("failed to zip files: %v", err)
	}

	err = r.handleStorage(zips)
	if err != nil {
		return fmt.Errorf("failed to handle storage: %v", err)
	}

	return nil
}

func (r *Runner) BuildLedger(destination string) error {
	mediaFiles, err := r.MusicFoldersRepository.AllMediaFiles()
	if err != nil {
		return fmt.Errorf("failed get media files from repository: %v", err)
	}
	absoluteMediaFiles, err := r.absoluteMediaFiles(mediaFiles)
	if err != nil {
		return fmt.Errorf("failed to get absolute media files: %v", err)
	}

	identifiedPaths := r.FilterService.IdentifiedPaths(absoluteMediaFiles, filter.NewMedia)

	identedIdentifiedPaths, err := json.MarshalIndent(identifiedPaths, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to MarshalIdent: %v", err)
	}
	err = r.FileSystemOperator.WriteNewFile(destination, identedIdentifiedPaths)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}

func (r *Runner) absoluteMediaFiles(mediaFiles []db.MediaFile) ([]db.MediaFile, error) {
	var absoluteMediaFiles []db.MediaFile
	for _, mediaFile := range mediaFiles {
		if filepath.IsAbs(mediaFile.Path) {
			absoluteMediaFiles = append(absoluteMediaFiles, mediaFile)
		} else {
			library, err := r.LibraryRepository.LibraryById(mediaFile.LibraryId)
			if library == nil {
				return nil, fmt.Errorf(
					"could not find library with ID of %v", mediaFile.LibraryId)
			}
			if err != nil {
				return nil, fmt.Errorf(
					"could not find library with ID of %v, error was %v", mediaFile.LibraryId, err)
			}
			absoluteMediaFiles = append(absoluteMediaFiles, db.MediaFile{
				Id:        mediaFile.Id,
				CreatedAt: mediaFile.CreatedAt,
				UpdatedAt: mediaFile.UpdatedAt,
				LibraryId: mediaFile.LibraryId,
				Path:      filepath.Join(library.Path, mediaFile.Path),
			})
		}
	}
	return absoluteMediaFiles, nil
}

func (r *Runner) zipFiles(identifiedPaths filter.IdentifiedPaths) (filter.IdentifiedPaths, error) {
	log.Printf("Zipping %v paths", len(identifiedPaths))
	zips := make(filter.IdentifiedPaths)

	for path, pathId := range identifiedPaths {
		log.Printf("Zipping %v", pathId.BasePath)
		zipPath, err := r.Zipper.ZipFilesInFolder(path)

		if err != nil {
			return nil, err
		}
		zips[zipPath] = filter.PathIdentifier{
			UploadType: pathId.UploadType,
			Id:         pathId.Id,
			BasePath:   pathId.BasePath,
		}
	}
	return zips, nil
}

func (r *Runner) handleStorage(zips filter.IdentifiedPaths) error {
	log.Printf("handleStorage for %v zips", len(zips))
	for path, pathIdentifier := range zips {
		log.Printf("%v is upload type %v", pathIdentifier.BasePath, pathIdentifier.UploadType)
		var err error

		destination := r.FilterService.UploadDestination(pathIdentifier)

		switch pathIdentifier.UploadType {
		case filter.NewMedia:
			err = r.StorageClient.UploadNewFile(path, destination)
		case filter.UpdatedMedia:
			err = r.StorageClient.ReplaceFile(path, destination)
		}

		if err != nil {
			return fmt.Errorf("failed to send %v to storage: %v", path, err)
		}

		err = r.FileSystemOperator.DeleteFile(path)
		if err != nil {
			return fmt.Errorf("failed to delete %v: %v", pathIdentifier.BasePath, err)
		}
		log.Printf("Finished handleStorage for %v", pathIdentifier.BasePath)
	}
	return nil
}
