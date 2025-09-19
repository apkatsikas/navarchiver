package filter_test

import (
	"github.com/apkatsikas/archiver/db"
	"github.com/apkatsikas/archiver/filter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IdentifiedPaths with duplicate paths", func() {
	var expectedIdentifiedPaths = filter.IdentifiedPaths{
		"/path/to/stuff":      filter.PathIdentifier{Id: "1abde123", UploadType: filter.NewMedia, BasePath: "stuff"},
		"/path/to/otherstuff": filter.PathIdentifier{Id: "aab123", UploadType: filter.NewMedia, BasePath: "otherstuff"},
	}

	var fs *filter.FilterService
	var identifiedPaths filter.IdentifiedPaths

	BeforeEach(func() {
		fs = &filter.FilterService{}
		identifiedPaths = fs.IdentifiedPaths([]db.MediaFile{
			{
				Id:   "abde123",
				Path: "/path/to/stuff/01 track.mp3",
			},
			{
				Id:   "1abde123",
				Path: "/path/to/stuff/02 track.mp3",
			},
			{
				Id:   "ba123",
				Path: "/path/to/otherstuff/01 track.mp3",
			},
			{
				Id:   "aab123",
				Path: "/path/to/otherstuff/02 track.mp3",
			},
		}, filter.NewMedia)
	})

	It("will consolidate the records into a map with path as the key using the lowest id", func() {
		Expect(identifiedPaths).To(Equal(expectedIdentifiedPaths))
	})
})

var _ = Describe("IdentifiedPaths with nil input", func() {
	var fs *filter.FilterService

	BeforeEach(func() {
		fs = &filter.FilterService{}
	})

	It("will return empty", func() {
		Expect(fs.IdentifiedPaths(nil, filter.NewMedia)).To(BeEmpty())
	})
})

var _ = Describe("UpdatedAndNewIdentifiedPaths", func() {
	var fs *filter.FilterService
	var identifiedPaths filter.IdentifiedPaths

	BeforeEach(func() {
		fs = &filter.FilterService{}

		newMediaFiles := []db.MediaFile{
			{
				Id:   "abde123",
				Path: "/path/to/stuff/01 track.mp3",
			},
			{
				Id:   "1abde123",
				Path: "/path/to/stuff/02 track.mp3",
			},
			{
				Id:   "ba123",
				Path: "/path/to/otherstuff/01 track.mp3",
			},
			{
				Id:   "aab123",
				Path: "/path/to/otherstuff/02 track.mp3",
			},
		}

		updatedMediaFiles := []db.MediaFile{
			{
				Id:   "abde123",
				Path: "/path/to/stuff/01 track.mp3",
			},
			{
				Id:   "1abde123",
				Path: "/path/to/stuff/02 track.mp3",
			},
			{
				Id:   "ba123",
				Path: "/path/to/otherstuff/01 track.mp3",
			},
			{
				Id:   "aab123",
				Path: "/path/to/otherstuff/02 track.mp3",
			},
			{
				Id:   "666abc420",
				Path: "/path/to/evenmorestuff/01 track.mp3",
			},
			{
				Id:   "0123",
				Path: "/path/to/evenmorestuff/02 track.mp3",
			},
		}

		identifiedPaths = fs.UpdatedAndNewIdentifiedPaths(newMediaFiles, updatedMediaFiles)
	})

	It("will consolidate records into a map with the path as the key using the lowest id"+
		", removing duplicate entries from updated in favor of new media", func() {
		Expect(identifiedPaths).To(Equal(filter.IdentifiedPaths{
			"/path/to/stuff":         {UploadType: filter.NewMedia, Id: "1abde123", BasePath: "stuff"},
			"/path/to/otherstuff":    {UploadType: filter.NewMedia, Id: "aab123", BasePath: "otherstuff"},
			"/path/to/evenmorestuff": {UploadType: filter.UpdatedMedia, Id: "0123", BasePath: "evenmorestuff"},
		}))
	})
})

var _ = Describe("UploadDestination", func() {
	var fs *filter.FilterService

	var pathIdentifier = filter.PathIdentifier{
		UploadType: filter.NewMedia,
		Id:         "abcd1234",
		BasePath:   "testpath",
	}
	var expectedDestination = "testpathabcd1234.zip"

	BeforeEach(func() {
		fs = &filter.FilterService{}
	})

	It("Will return the expected destination", func() {
		Expect(fs.UploadDestination(pathIdentifier)).To(Equal(expectedDestination))
	})
})
