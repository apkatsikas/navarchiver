package filter

import (
	"fmt"
	"maps"
	"path/filepath"
	"strings"

	"github.com/apkatsikas/archiver/db"
)

type FilterService struct {
}

type UploadType int

const (
	NewMedia UploadType = iota
	UpdatedMedia
)

type PathIdentifier struct {
	UploadType
	Id       string
	BasePath string
}

type IdentifiedPaths map[string]PathIdentifier

func (fs *FilterService) IdentifiedPaths(mediaFiles []db.MediaFile, uploadType UploadType) IdentifiedPaths {
	identifiedPaths := make(IdentifiedPaths)
	for _, mediaFile := range mediaFiles {

		pathDirectory := filepath.Dir(mediaFile.Path)
		basePath := filepath.Base(pathDirectory)

		if _, pathExists := identifiedPaths[pathDirectory]; pathExists {
			if fs.isIdLowerThanExistingId(mediaFile.Id, identifiedPaths[pathDirectory].Id) {
				identifiedPaths[pathDirectory] = PathIdentifier{
					Id:         mediaFile.Id,
					UploadType: uploadType,
					BasePath:   basePath,
				}
			}
		} else {
			identifiedPaths[pathDirectory] = PathIdentifier{Id: mediaFile.Id, UploadType: uploadType, BasePath: basePath}
		}
	}

	return identifiedPaths
}

func (fs *FilterService) UpdatedAndNewIdentifiedPaths(newMediaFiles []db.MediaFile, updatedMediaFiles []db.MediaFile) IdentifiedPaths {
	newIdentifiedPaths := fs.IdentifiedPaths(newMediaFiles, NewMedia)
	updatedIdentifiedPaths := fs.IdentifiedPaths(updatedMediaFiles, UpdatedMedia)

	fs.removeDupesFromUpdatedFiles(updatedIdentifiedPaths, newIdentifiedPaths)

	maps.Copy(newIdentifiedPaths, updatedIdentifiedPaths)

	return newIdentifiedPaths
}

func (fs *FilterService) UploadDestination(pathIdentifider PathIdentifier) string {
	return fmt.Sprintf("%v%v.zip", pathIdentifider.BasePath, pathIdentifider.Id)
}

// If a path is considered both new AND updated, we can just consider it new
func (fs *FilterService) removeDupesFromUpdatedFiles(updatedIdentifiedPaths IdentifiedPaths, newIdentifiedPaths IdentifiedPaths) {
	for updatedPath := range updatedIdentifiedPaths {
		if _, pathExists := newIdentifiedPaths[updatedPath]; pathExists {
			delete(updatedIdentifiedPaths, updatedPath)
		}
	}

	maps.Copy(newIdentifiedPaths, updatedIdentifiedPaths)
}

func (fs *FilterService) isIdLowerThanExistingId(id string, existingId string) bool {
	return strings.Compare(id, existingId) < 0
}
