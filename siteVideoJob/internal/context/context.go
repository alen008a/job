package context

import (
	"context"
	"fmt"
	"github.com/rs/xid"
	"os"
	"path"
	"runtime"
	"siteVideoJob/config"
	"siteVideoJob/internal/glog/bot"
	"siteVideoJob/internal/glog/log"
	"strconv"
	"strings"
	"time"
)

type Handler func(ctx *Context)

type Context struct {
	Ctx        *context.Context
	CancelFunc context.CancelFunc
	logs       strings.Builder
	SiteId     string
	Trace      string
}

func Background(timeout time.Duration) *Context {
	background := context.Background()

	if timeout > 0 {
		ctx, cancelFunc := context.WithTimeout(background, timeout)
		return &Context{Ctx: &ctx, CancelFunc: cancelFunc, logs: strings.Builder{}, Trace: xid.New().String()}
	}
	ctx, cancelFunc := context.WithCancel(background)
	return &Context{Ctx: &ctx, CancelFunc: cancelFunc, logs: strings.Builder{}, Trace: xid.New().String()}
}

func (c *Context) Info(args ...interface{}) {
	log.ZapLog.Named(c.funcName()).Info(args...)
}

func (c *Context) Infof(template string, args ...interface{}) {
	log.ZapLog.Named(c.funcName()).Infof(template, args...)
}

func (c *Context) Warn(args ...interface{}) {
	log.ZapLog.Named(c.funcName()).Warn(args...)
}

func (c *Context) Warnf(template string, args ...interface{}) {
	log.ZapLog.Named(c.funcName()).Warnf(template, args...)
}

func (c *Context) Error(args ...interface{}) {
	log.ZapLog.Named(c.funcName()).Error(args...)
}

func (c *Context) Debug(args ...interface{}) {
	log.ZapLog.Named(c.funcName()).Debug(args...)
}

func (c *Context) Debugf(template string, args ...interface{}) {
	log.ZapLog.Named(c.funcName()).Debugf(template, args...)
}

func (c *Context) Errorf(template string, args ...interface{}) {
	log.ZapLog.Named(c.funcName()).Errorf(template, args...)
}

func (c *Context) funcName() string {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s, TraceId: %s, siteId: %s", path.Base(funcName), c.Trace, c.SiteId)
}

// 非常重要日志警告
func (c *Context) Emergency(template string, args ...interface{}) {
	s1, s2 := funcName4Emergency(c.Trace)
	bot.Send(c.SiteId, bot.SlackTemplate, s2, fmt.Sprintf(template, args...))
	log.ZapLog.Named(s1).Warnf(template, args...)
}

// 仅提供给中间件使用，告警慢查询接口
func (c *Context) Tracef(latency float64, template, path string, args ...interface{}) {
	if latency > 20 {
		bot.SendMid(c.SiteId, fmt.Sprintf(bot.SlowTemplate, config.GetApp().AppID, c.SiteId+config.GetApp().Env, path, latency, c.Trace), path)
	}
	log.ZapLog.Named(c.funcName()).Infof(template, args...)
}

func funcName4Emergency(args ...string) (string, string) {
	pc, f, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()

	//获取上一层的stack
	index := lastIndexByte(f, os.PathSeparator)
	if index != -1 {
		f = f[index+1:]
	}
	return path.Base(funcName) + " " + strings.Join(args, " "),
		path.Base(funcName) + " " + strings.Join(args, " ") + " " + f + ":" + strconv.Itoa(line) + " "
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

// Console Log 调度中心控制台输出，请勿打印敏感信息
func (c *Context) Console(msg string, args ...interface{}) {
	if len(args) == 0 {
		c.logs.WriteString(msg)
	} else {
		c.logs.WriteString(fmt.Sprintf(msg, args...))
	}
	// 换行
	c.logs.WriteString("<br>")
	// 文件同步到日志文件
	c.Infof(msg, args...)
}

// ConsoleErr 记录错误日志信息
func (c *Context) ConsoleErr(msg string, args ...interface{}) {
	c.logs.WriteString("<text style='color:red'>")
	if len(args) == 0 {
		c.logs.WriteString(msg)
	} else {
		c.logs.WriteString(fmt.Sprintf(msg, args...))
	}
	c.logs.WriteString("</text>")
	// 换行
	c.logs.WriteString("<br>")
	// 文件同步到日志文件
	c.Errorf(msg, args...)
}

// GetConsoleLog 获取调度中心控制台输出
func (c *Context) GetConsoleLog() string {
	return c.logs.String()
}
