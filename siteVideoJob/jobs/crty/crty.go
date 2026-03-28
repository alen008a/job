package crty

import (
	"siteVideoJob/internal/context"
	"siteVideoJob/mdata"
	"siteVideoJob/service"
	"siteVideoJob/xxl"
	"strconv"
)

// PullCRTYMatch 拉取CR体育比赛列表
func PullCRTYMatch(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams *xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullCRTYMatch 开始执行")
	msg = service.PullCRTYMatchList(c)
	c.Infof("PullCRTYMatch 执行结束 result:%s", msg)
	c.Console("PullCRTYMatch 执行结束 result:%s", msg)
	return c.GetConsoleLog()
}

// DoSetCrtyFromLsty 从LSTY匹配合适的视频源
func DoSetCrtyFromLsty(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams *xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("DoSetCrtyFromLsty 开始执行")
	msg = service.DoSetCrtyFromLsty(c)
	c.Infof("DoSetCrtyFromLsty 执行结束 result:%s", msg)
	c.Console("DoSetCrtyFromLsty 执行结束 result:%s", msg)
	return c.GetConsoleLog()
}
