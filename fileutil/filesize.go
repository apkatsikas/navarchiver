package fileutil

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

type FileSize uint

func (fs FileSize) String() string {
	return humanize.Bytes(uint64(fs))
}

func (fs *FileSize) Set(value string) error {
	convertedUint, err := humanize.ParseBytes(value)
	if err != nil {
		return fmt.Errorf("failed to convert %v to uint: %v", convertedUint, err)
	}
	if convertedUint <= 0 {
		return fmt.Errorf("file size limit must be at least 1 byte")
	}
	*fs = FileSize(convertedUint)
	return nil
}
