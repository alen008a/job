package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"siteLetterJob/config"
	"siteLetterJob/db"
	"siteLetterJob/internal/glog"
	"siteLetterJob/internal/glog/log"
	"siteLetterJob/lib"
	"siteLetterJob/lib/kfk"
	"siteLetterJob/service/sitemsg"
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
		panic("main Init error:" + err.Error())
	}
	logWriter := log.InitLog()
	if logWriter == nil {
		panic("init log failed")
	}

	err = db.InitDB()
	defer db.Close()
	if err != nil {
		glog.Error("init db err: ", err)
		return
	}

	err = lib.InitLib()
	if err != nil {
		panic("lib init error:" + err.Error())
	}

	// kfk消息队列监听
	conf := config.GetConfig()
	go kfk.RegisterKfkConsumerTopicsListener(conf.KafkaAddrV2, conf.SiteMsgKafkaTopic, conf.SiteMsgKafkaGroup, "2", sitemsg.SiteLetterConsumeHandle)
	// 站内模板消息消费监听
	go kfk.RegisterKfkConsumerTopicsListener(conf.KafkaAddrV2, conf.SiteMsgTemplateTopic, conf.SiteMsgTemplateGroup, "2", sitemsg.TemplateConsumer)
	// 站内广播消息消费监听
	go kfk.RegisterKfkConsumerTopicsListener(conf.KafkaAddrV2, conf.SiteInnerMsgBrTopic, conf.SiteInnerMsgBrGroup, "2", sitemsg.MsgBrConsumer)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	glog.Infof("START AT: %s, GIT: %s", time.Now(), CommitID)

	glog.Info("siteLetterJob task start")
	<-done
	glog.Info("siteLetterJob task end")
}
