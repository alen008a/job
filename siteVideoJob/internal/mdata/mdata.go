package mdata

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"strings"
	"sync"
)

var (
	Cjson = jsoniter.ConfigCompatibleWithStandardLibrary
)

// Logger 接口实现
type Logger struct{}

func (l *Logger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf("系统监控日志 - "+format, a...))
}

func (l *Logger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf("系统监控日志 - "+format, a...))
}

// Console 调度中心控制台日志打印结构
type Console struct {
	mu     sync.Mutex
	logTxt strings.Builder
}

func (c *Console) Log(format string, args ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logTxt.WriteString(fmt.Sprintf(format, args...))
}

func (c *Console) GetLog() string {
	return c.logTxt.String()
}
