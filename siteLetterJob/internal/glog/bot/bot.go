package bot

import (
	"fmt"
	"os"
	"siteLetterJob/config"
	"siteLetterJob/db/redisdb/core"
	"siteLetterJob/lib/httpclient"
	"siteLetterJob/mdata"
	"strconv"
	"sync"
	"time"
)

// 告警的前提是不影响业务
var store sync.Map                      //只存代码位置，减少内存占用
var warningMsg = make(chan string, 100) //增加容量，防止业务慢
var cacheKey = "SITE_BOT_WARNING_CACHE_%s"
var timeout int64 = 600
var hostName string

func init() {
	hostName, _ = os.Hostname()
}

func SendDefault(template, position, msg string) {
	Send(template, position, msg, "初始化站点")
}
func Send(template, position, msg, siteId string) {

	if config.GetConfig().WarningCode == "" {
		return
	}

	var (
		now   = time.Now().Unix()
		t, ok = store.Load(position)
	)

	//相同的代码位置，十分钟只推送一次，防止大批量推送导致系统奔溃和漏掉通知
	//如果不存在就从redis去获取,如果redis存在，然后刷到内存，并直接返回
	if !ok {
		rs, err := core.GetKey(cacheKey + position)
		if rs != "" && err == nil {
			t1, _ := strconv.ParseInt(rs, 10, 64)
			if t1 > 0 && now-t1 < timeout {
				//刷到内存里
				store.Store(position, t1)
				return
			}
		}
	}

	if ok && now-t.(int64) < timeout {
		return
	}

	//刷新时间
	if ok && now-t.(int64) > timeout {
		store.Delete(position)
	}

	select {
	case <-time.After(time.Millisecond * 100): //不能因为量大的时候阻塞主要业务
	case warningMsg <- fmt.Sprintf(template, config.GetApp().AppID, siteId+"-"+config.GetApp().Env, hostName, position, msg):
	}

	//如果是第一次就存入
	if !ok {
		store.Store(position, now)
		core.SetExpireKV(fmt.Sprintf(cacheKey, siteId)+position, strconv.FormatInt(now, 10), time.Second*time.Duration(timeout))
	}
}

func SendMid(msg, path string) {

	if config.GetConfig().WarningCode == "" {
		return
	}

	//相同的代码位置，十分钟只推送一次，防止大批量推送导致系统奔溃和漏掉通知
	var (
		now   = time.Now().Unix()
		t, ok = store.Load(path)
	)

	//如果不存在就从redis去获取,如果redis存在，然后刷到内存，并直接返回
	if !ok {
		rs, err := core.GetKey(cacheKey + path)
		if rs != "" && err == nil {
			t1, _ := strconv.ParseInt(rs, 10, 64)
			if t1 > 0 && now-t1 < timeout {
				//刷到内存里
				store.Store(path, t1)
				return
			}
		}
	}

	if ok && now-t.(int64) < timeout {
		return
	}

	//刷新时间
	if ok && now-t.(int64) > timeout {
		store.Delete(path)
	}

	select {
	case <-time.After(time.Millisecond * 100): //不能因为量大的时候阻塞主要业务
	case warningMsg <- msg:
	}

	//如果是第一次就存入
	if !ok {
		store.Store(path, now)
		core.SetExpireKV(cacheKey+path, strconv.FormatInt(now, 10), time.Second*time.Duration(timeout))
	}
}

func send2bot(serviceCode, msg string, botType int) (string, error) {
	data, err := httpclient.POSTJson(
		config.GetConfig().VerifyCodeDomain+"/verifycode/bot/v1/send",
		mdata.MustMarshal2Byte(map[string]interface{}{
			"serviceCode": serviceCode,
			"botType":     botType,
			"msg":         msg,
		}),
		map[string]string{"Content-Type": "application/json", mdata.HeaderSite: "9999"},
		httpclient.GetShortProxyNotifyClient(time.Second*3),
	)
	return string(data), err
}

func init() {
	go watchDog()
}

func watchDog() {
	go func() {
		for v := range warningMsg {
			data, err := send2bot(config.GetConfig().WarningCode, v, 3)
			if err != nil {
				_ = fmt.Errorf("send2bot occur err=%v |data=%v |msg=%s", err, data, v)
			}
		}
	}()
}

const SlackTemplate = `
=============系统警告，开发人员需要注意=============
项目：【%s】
站点：【%s】
主机：【%s】
代码位置：
	%s
告警内容：
	%s

`

const SlowTemplate = `
-------------慢接口警告，开发人员需要注意-------------
项目：【%s】
站点：【%s】
接口地址：%s
接口耗时：%.2fs
日志追踪：%s


`
