package glog

import (
	ctx "context"
	"fmt"
	"os"
	"path"
	"runtime"
	"siteVideoJob/config"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//在线上环境中，后期sql日志过多需要去掉的时候，可以选择是否打开某些关键的sql日志
type dbLog struct {
	logger.Config
	traceSwitch bool
}

func NewDBLog(config logger.Config) *dbLog {
	return &dbLog{Config: config, traceSwitch: true}
}

func (l *dbLog) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	newLogger.traceSwitch = true
	//超过3秒都算慢查询
	newLogger.SlowThreshold = time.Second * 1
	return &newLogger
}

func (l *dbLog) Info(ctx ctx.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Info {
		Info(fmt.Sprintf("DBLOG |"+s+"\n", i...))
	}
}

func (l *dbLog) Warn(ctx ctx.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Warn {
		Warn(fmt.Sprintf("DBLOG |"+s+"\n", i...))
	}
}

func (l *dbLog) Error(ctx ctx.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Error {
		Error(fmt.Sprintf("DBLOG |"+s+"\n", i...))
	}
}

func (l *dbLog) Trace(ctx ctx.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent || !l.traceSwitch {
		return
	}

	//耗时
	timestamp := time.Since(begin)

	sql, rows := fc()

	if strings.ToLower(config.GetEnv()) == "dev" {
		Infof("执行：%s，影响行：%d", sql, rows)
	}

	//如果报错直接打印日志
	if err != nil && err != gorm.ErrRecordNotFound {
		Emergency(
			"DBLOG |pos=%s |err=%v |sql=%s |rows=%d |timestamp=%.2fs",
			funcName4Gorm(),
			err,
			sql,
			rows,
			timestamp.Seconds(),
		)
		return
	}

	//产生慢查询
	if timestamp > l.SlowThreshold {
		//记录不同标准的慢sql
		var standard = "PURPLE"
		switch {
		case timestamp > l.SlowThreshold*5:
			standard = "RED"
			Emergency(
				"DBLOG |SLOW |pos=%s |standard=%s |sql=%s |rows=%d |timestamp=%.2fs",
				funcName4Gorm(),
				standard,
				sql,
				rows,
				timestamp.Seconds(),
			)
			return
		case timestamp > l.SlowThreshold*4:
			standard = "ORANGE"
		case timestamp > l.SlowThreshold*2:
			standard = "YELLOW"
		}

		Warnf(
			"DBLOG |SLOW |pos=%s |standard=%s |sql=%s |rows=%d |timestamp=%.2fs",
			funcName4Gorm(),
			standard,
			sql,
			rows,
			timestamp.Seconds(),
		)
		return
	}

	//如果不是慢日志，只有数据变更操作的日志才记录
	if rows > 0 && (strings.Contains(sql, "INSERT") || strings.Contains(sql, "UPDATE") || strings.Contains(sql, "DELETE")) {
		Infof(
			"DBLOG |pos=%s|sql=%s |rows=%d |timestamp=%.2fs",
			funcName4Gorm(),
			sql,
			rows,
			timestamp.Seconds(),
		)
		return
	}

	//更新0的日志额外重点记录
	if strings.Contains(sql, "UPDATE") {
		Warnf(
			"DBLOG |pos=%s |sql=%s |rows=%d |timestamp=%.2fs",
			funcName4Gorm(),
			sql,
			rows,
			timestamp.Seconds(),
		)
		return
	}
}

func funcName4Gorm() string {
	pc, f, line, _ := runtime.Caller(4)
	funcName := runtime.FuncForPC(pc).Name()

	//获取上一层的stack
	index := lastIndexByte(f, os.PathSeparator)
	if index != -1 {
		f = f[index+1:]
	}
	return path.Base(funcName) + " " + f + ":" + strconv.Itoa(line) + " "
}
