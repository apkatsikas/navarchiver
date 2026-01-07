package fileutil

import (
	"fmt"
	"strconv"
)

type FileCount uint

func (fc FileCount) String() string {
	return strconv.FormatUint(uint64(fc), 10)
}

func (fc *FileCount) Set(value string) error {
	convertedUint, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert %v to uint: %v", convertedUint, err)
	}
	if convertedUint <= 0 {
		return fmt.Errorf("file count limit must be greater than 0")
	}
	*fc = FileCount(convertedUint)
	return nil
}
