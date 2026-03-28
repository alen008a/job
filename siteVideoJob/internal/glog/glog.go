package glog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"siteVideoJob/internal/glog/bot"
	"strconv"

	. "siteVideoJob/internal/glog/log"
)

func Info(args ...interface{}) {
	ZapLog.Named(funcName()).Info(args...)
}

func Infof(template string, args ...interface{}) {
	ZapLog.Named(funcName()).Infof(template, args...)
}

func Warn(args ...interface{}) {
	ZapLog.Named(funcName()).Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	ZapLog.Named(funcName()).Warnf(template, args...)
}

func Error(args ...interface{}) {
	ZapLog.Named(funcName()).Error(args...)
}

func Debug(args ...interface{}) {
	ZapLog.Named(funcName()).Debug(args...)
}

func Errorf(template string, args ...interface{}) {
	ZapLog.Named(funcName()).Errorf(template, args...)
}

func funcName() string {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	return path.Base(funcName)
}

// 非常重要日志警告 用的预设推送间隔
func Emergency(template string, args ...interface{}) {
	EmergencyWithTimeout(template, 0, args...)
}

// 自订推送间格
func EmergencyWithTimeout(template string, timeout int64, args ...interface{}) {
	s1, s2 := funcName4Emergency()

	if timeout == 0 {
		bot.SendDefault(bot.SlackTemplate, s2, fmt.Sprintf(template, args...))
	} else {
		bot.SendDefaultWithTimeout(bot.SlackTemplate, s2, fmt.Sprintf(template, args...), timeout)
	}
	ZapLog.Named(s1).Warnf(template, args...)
}

func funcName4Emergency() (string, string) {
	pc, f, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()

	//获取上一层的stack
	index := lastIndexByte(f, os.PathSeparator)
	if index != -1 {
		f = f[index+1:]
	}
	return path.Base(funcName),
		path.Base(funcName) + " " + f + ":" + strconv.Itoa(line) + " "
}

func lastIndexByte(s string, c byte) int {
	var count int
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			count++
		}

		if count == 2 {
			return i
		}
	}
	return -1
}
