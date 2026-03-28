package sitemsg

import (
	"fmt"
	"runtime"
	confi "siteLetterJob/config"
	"siteLetterJob/internal/context"
	"siteLetterJob/lib/httpclient"
	"siteLetterJob/mdata"
	"siteLetterJob/mdata/sitemsg"
	"siteLetterJob/utils"
	"strings"
)

var (
	Dev_ALL           = sitemsg.DeviceTypeEnum{Code: "all", CodeName: "所有"}
	DEFAULT           = sitemsg.PushTypeEnum{Code: 0, CodeName: "默认"}
	PUSH_ALL          = sitemsg.PushTypeEnum{Code: 1, CodeName: "广播所有人"}
	PUSH_LABEL        = sitemsg.PushTypeEnum{Code: 2, CodeName: "设备标签"}
	PUSH_ALIAS        = sitemsg.PushTypeEnum{Code: 3, CodeName: "设备别名"}
	PUSH_REGISTRATION = sitemsg.PushTypeEnum{Code: 4, CodeName: "Registration ID"}
)

func recoverHandle(c *context.Context) {
	if err := recover(); err != nil {
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		tmpStr := fmt.Sprintf("err=%v panic ==> %s\n", err, string(buf[:n]))
		c.Error(tmpStr)
		fmt.Println(tmpStr)
	}
}

// 调取adam的接口完成极光推送
func SendPush(c *context.Context, pushReqx *sitemsg.SendPushReq) {
	defer recoverHandle(c)
	// todo 极光推送处理
	pushReq := sitemsg.PushReq{}
	pushReq.SiteId = pushReqx.SiteId
	pushReq.MemberIds = pushReqx.MemberIds
	pushReq.PlatformList = pushReqx.PlatformList
	pushReq.PushTitle = pushReqx.PushTitle
	pushReq.PushContent = pushReqx.PushContent
	pushReq.Status = 1
	pushReq.TimeLive = 89600
	pushReq.Type = pushReqx.Type
	pushReq.Alias = []string{}
	c.Info("JgSend start : ", pushReq)
	jsonStr, _ := mdata.Cjson.MarshalToString(pushReq)
	c.Info("JgSend POSTJson url : ", confi.GetConfig().VerifyCodeDomain+confi.GetConfig().PushUrl, " str : ", jsonStr)

	retry := confi.GetConfig().PushRetry
	if retry <= 0 {
		retry = 3
	}
	for j := 1; j <= retry; j++ {
		resB, err := httpclient.ProxyPostJson(confi.GetConfig().VerifyCodeDomain+confi.GetConfig().PushUrl, jsonStr, map[string]string{mdata.HeaderSite: pushReqx.SiteId})
		if err != nil {
			c.Errorf("当前第%d次发送,发送参数=%s,header=%v,ulr=%s,网络请求错误=%v", j, jsonStr, nil,
				confi.GetConfig().VerifyCodeDomain+confi.GetConfig().PushUrl, err)
			//客户端超时时间为30秒，如果还是返回超时，不进行重试，作为成功来处理。如果没收到推送，手动重推一下。
			//将来可能需要再次优化
			if strings.Contains(err.Error(), "context deadline exceeded") {
				break
			}
			continue
		}
		req := map[string]interface{}{}
		err = mdata.Cjson.Unmarshal(resB, &req)
		if err != nil {
			c.Errorf("当前第%d次发送,发送参数=%s,header=%v,ulr=%s,发送结果反序列化错误=%v", j, jsonStr, nil,
				confi.GetConfig().VerifyCodeDomain+confi.GetConfig().PushUrl, err)
			continue
		}
		c.Infof("当前第%d次发送,发送参数=%s,header=%v,ulr=%s,返回数据=%s", j, jsonStr, nil,
			confi.GetConfig().VerifyCodeDomain+confi.GetConfig().PushUrl, string(resB))
		if utils.InterfaceToInt(utils.InterfaceToString(req["status_code"])) == 6000 {
			c.Info("JgSend POSTJson success data %s res : +%s", jsonStr, string(resB))
			break
		}
	}
	return
}
