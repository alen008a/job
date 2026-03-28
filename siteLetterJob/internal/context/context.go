package context

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"siteLetterJob/config"
	"siteLetterJob/internal/glog/bot"
	. "siteLetterJob/internal/glog/log"
	"strconv"
)

type Handler func(ctx *Context)

type Context struct {
	SiteId       string
	Trace        string
	KafkaVersion string
	ServerInfo   string
}

func (c *Context) Info(args ...interface{}) {
	ZapLog.Named(c.funcName()).Info(args...)
}

func (c *Context) Infof(template string, args ...interface{}) {
	ZapLog.Named(c.funcName()).Infof(template, args...)
}

func (c *Context) Warn(args ...interface{}) {
	ZapLog.Named(c.funcName()).Warn(args...)
}

func (c *Context) Warnf(template string, args ...interface{}) {
	ZapLog.Named(c.funcName()).Warnf(template, args...)
}

func (c *Context) Error(args ...interface{}) {
	ZapLog.Named(c.funcName()).Error(args...)
}

func (c *Context) Debug(args ...interface{}) {
	ZapLog.Named(c.funcName()).Debug(args...)
}

func (c *Context) Debugf(template string, args ...interface{}) {
	ZapLog.Named(c.funcName()).Debugf(template, args...)
}

func (c *Context) Errorf(template string, args ...interface{}) {
	ZapLog.Named(c.funcName()).Errorf(template, args...)
}

func (c *Context) funcName() string {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s, Trace: %s, siteId: %s kafkaVersion: %s serverInfo: %s", path.Base(funcName), c.Trace, c.SiteId, c.KafkaVersion, c.ServerInfo)
}

// 非常重要日志警告
func (c *Context) Emergency(template string, args ...interface{}) {
	s1, s2 := funcName4Emergency(c.Trace)
	bot.Send(bot.SlackTemplate, s2, fmt.Sprintf(template, args...), c.SiteId)
	ZapLog.Named(s1).Warnf(template, args...)
}

// 仅提供给中间件使用，告警慢查询接口
func (c *Context) Tracef(latency float64, template, path string, args ...interface{}) {
	if latency > 20 {
		bot.SendMid(fmt.Sprintf(bot.SlowTemplate, config.GetApp().AppID, c.SiteId+"-"+config.GetApp().Env, path, latency, c.Trace), path)
	}
	ZapLog.Named(c.funcName()).Infof(template, args...)
}

func funcName4Emergency(trace string) (string, string) {
	pc, f, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()

	//获取上一层的stack
	index := lastIndexByte(f, os.PathSeparator)
	if index != -1 {
		f = f[index+1:]
	}
	return path.Base(funcName) + " " + trace,
		path.Base(funcName) + " " + trace + " " + f + ":" + strconv.Itoa(line) + " "
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
