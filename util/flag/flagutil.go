package flagutil

import (
	"flag"
	"fmt"
	"sync"

	"github.com/apkatsikas/archiver/fileutil"
)

type RunMode string

func (rm *RunMode) String() string {
	return string(*rm)
}

func (rm *RunMode) Set(value string) error {
	switch value {
	case string(RunModeScheduled), string(RunModeBatch), string(RunModeLedger):
		*rm = RunMode(value)
		return nil
	default:
		return fmt.Errorf("invalid value for runMode: %s", value)
	}
}

const (
	RunModeScheduled RunMode = "scheduled"
	RunModeBatch     RunMode = "batch"
	RunModeLedger    RunMode = "ledger"
)

type FlagUtil struct {
	RunMode        RunMode
	FileSizeLimit  fileutil.FileSize
	FileCountLimit fileutil.FileCount
}

func (fu *FlagUtil) Setup() {
	flag.Var(&fu.RunMode, "runMode",
		"Which mode to run archiver in - valid values are "+
			"'scheduled', 'batch' or 'ledger' - default is scheduled")
	flag.Var(&fu.FileSizeLimit, "fileSizeLimit", "Maximum size for a file, if exceeded the archiver will throw an error")
	flag.Var(&fu.FileCountLimit, "fileCountLimit", "Maximum number of files allowed in a folder, if exeeded the archiver will throw an error")
	flag.Parse()
}

var (
	fu     *FlagUtil
	fuOnce sync.Once
)

func Get() *FlagUtil {
	if fu == nil {
		fuOnce.Do(func() {
			fu = &FlagUtil{}
		})
	}
	return fu
}
