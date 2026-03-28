package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

var parser = cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
var C = cron.New(cron.WithParser(GetParser()),
	cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
)

var ErrNoFind = fmt.Errorf("no find!")

// Start 开始定时任务
func Start() {
	C.Start()
}

// Stop 停止定时任务
func Stop() {
	C.Stop()
}

// AddFunc 1.删除将来运行的任务, 不影响当前任务. 2.重新添加一个新任务
func AddFunc(oldId cron.EntryID, spec string, job func()) (cron.EntryID, error) {
	if oldId != 0 {
		C.Remove(oldId)
	}

	return C.AddFunc(spec, job)
}

func NextTime(id cron.EntryID) (time.Time, error) {
	entry := C.Entry(id)
	if entry.ID == 0 {
		return time.Time{}, ErrNoFind
	}
	return entry.Next, nil
}

func PreTime(id cron.EntryID) (time.Time, error) {
	entry := C.Entry(id)
	if entry.ID == 0 {
		return time.Time{}, ErrNoFind
	}
	return entry.Prev, nil
}

func GetParser() cron.Parser {
	return parser
}
