package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"siteVideoJob/config"
	"siteVideoJob/db"
	"siteVideoJob/internal/glog"
	"siteVideoJob/internal/glog/log"
	"siteVideoJob/internal/mdata"
	"siteVideoJob/internal/middleware"
	"siteVideoJob/jobs"
	"siteVideoJob/lib/rp"
	"siteVideoJob/utils"
	"siteVideoJob/xxl"
	"syscall"
	"time"
)

var (
	confPath = flag.String("config", "./config/app.local.ini", "profilePath")
	CommitID string
)

func Init() error {

	flag.Parse()

	err := config.InitConfig(*confPath)
	if err != nil {
		return fmt.Errorf("init config is err: %v", err)
	}

	return nil
}

func recoverHandle() {
	if err := recover(); err != nil {
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		glog.Emergency("err=%v panic ==> %s\n", err, string(buf[:n]))
	}
}

func main() {
	defer recoverHandle()
	time.Local = time.FixedZone("CST", 3600*8)
	err := Init()
	if err != nil {
		panic(err)
	}
	logWriter := log.InitLog()
	if logWriter == nil {
		panic("init log failed")
	}

	glog.Info("开始连接数据库")
	st := time.Now().UnixNano()
	err = db.InitDB()
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		glog.Error("init db err: ", err)
		return
	}
	et := time.Now().UnixNano()
	glog.Infof("数据库连接成功，耗时%d毫秒", (et-st)/1e6)

	// 定时任务
	defer func() {
		rp.ReleaseGlobal()
	}()
	//启动协程池
	rp.InitGlobal(300)

	executeIP := config.GetConfig().Application.ExecutorIp
	if executeIP == "" {
		executeIP, err = utils.GetInternalIPv4()
		if err != nil {
			glog.Errorf("get ip v4 err %v", err.Error())
		}
	}
	execute := xxl.CreateExecutor(
		xxl.ServerAddr(config.GetConfig().Application.AdminServer),    // xxl-job-admin 服务器地址
		xxl.AccessToken(config.GetConfig().Application.AccessToken),   //请求令牌(默认为空)
		xxl.RegistryKey(config.GetConfig().Application.ExecutorName),  //执行器名称
		xxl.ExecutorIp(executeIP),                                     //可自动获取
		xxl.ExecutorPort(config.GetConfig().Application.ExecutorPort), //该脚本执行器端口号默认9999（非必填）
		xxl.SetLogger(mdata.Logger{}),                                 //自定义日志
	)

	execute.Init()
	//设置使用自定义中间件
	execute.Use(middleware.CustomMiddleware)
	//设置日志查看handler
	execute.LogHandler(xxl.GetDBLogHandle)
	//注册任务列表
	jobs.RegisterExecutors(execute)

	glog.Info(execute.Run())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	glog.Infof("START AT: %s, GIT: %s", time.Now(), CommitID)

	glog.Info("siteVideoJob task start")
	<-done
	glog.Info("siteVideoJob task end")
}
