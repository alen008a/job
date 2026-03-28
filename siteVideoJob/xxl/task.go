package xxl

import (
	"context"
	"fmt"
	"runtime/debug"
	ctx "siteVideoJob/internal/context"
	"siteVideoJob/internal/mdata"
)

// TaskFunc 任务执行函数
type TaskFunc func(cxt *ctx.Context, param *RunReq) string

// Task 任务
type Task struct {
	Id        int64
	Name      string
	Ext       *ctx.Context
	Param     *RunReq
	fn        TaskFunc
	Cancel    context.CancelFunc
	StartTime int64
	EndTime   int64
	//日志
	Log mdata.Logger
}

// Run 运行任务
func (t *Task) Run(callback func(code int64, msg string)) {
	defer func(cancel func()) {
		if err := recover(); err != nil {
			t.Ext.ConsoleErr(t.getInfo()+" panic: %v", err)
			debug.PrintStack() //堆栈跟踪
			callback(FailureCode, t.Ext.GetConsoleLog())
			cancel()
		}
	}(t.Ext.CancelFunc)
	msg := t.fn(t.Ext, t.Param)
	callback(SuccessCode, msg)
	return
}

// Info 任务信息
func (t *Task) getInfo() string {
	return fmt.Sprintf("任务ID: [%d]<br>任务名称: [%s]<br>参数: %s<br>", t.Id, t.Name, t.Param.ExecutorParams)
}
