package flagutil

import (
	"flag"
	"fmt"
	"sync"
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
	RunMode RunMode
}

func (fu *FlagUtil) Setup() {
	flag.Var(&fu.RunMode, "runMode",
		"Which mode to run archiver in - valid values are "+
			"'scheduled', 'batch' or 'ledger' - default is scheduled")
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
