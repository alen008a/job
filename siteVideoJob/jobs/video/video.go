package video

import (
	"siteVideoJob/internal/context"
	"siteVideoJob/mdata"
	"siteVideoJob/service"
	"siteVideoJob/xxl"
	"strconv"
)

func ClearFinishedMatches(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("ClearFinishedMatches 开始执行")
	_, err := service.DeleteLastDayMatch(c)
	if err != nil {
		c.Infof("ClearFinishedMatches err = %v", err)
		return c.GetConsoleLog()
	}
	c.Infof("ClearFinishedMatches 执行结束")
	return c.GetConsoleLog()
}

// SiteAdminUpdateLiveVideo 热门直播站点后台视频定时更新数据
func SiteAdminUpdateLiveVideo(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("SiteAdminUpdateLiveVideo 开始执行")
	service.SiteAdminUpdateLiveVideo(c)
	c.Infof("SiteAdminUpdateLiveVideo 执行结束")
	return c.GetConsoleLog()
}

func PullLiveVideoEvents(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullLiveVideoEvents 开始执行")
	service.PullLiveEvents(c)
	c.Infof("PullLiveVideoEvents 执行结束")
	return c.GetConsoleLog()
}

func PullLiveVideoOddsEvents(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullLiveVideoOddsEvents 开始执行")
	service.PullLiveOdds(c)
	c.Infof("PullLiveVideoOddsEvents 执行结束")
	return c.GetConsoleLog()
}

func SetLiveVideoEvents(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullLiveVideoEvents 开始执行")
	service.SetLiveVideoData(c)
	c.Infof("PullLiveVideoEvents 执行结束")
	return c.GetConsoleLog()
}

func PullAnchorEvents(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullAnchorEvents 开始执行")
	service.PullAnchorEvents(c)
	c.Infof("PullAnchorEvents 执行结束")
	return c.GetConsoleLog()
}

func UpdateAnchorEventList(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("UpdateAnchorEventList 开始执行")
	service.UpdateAnchorEventList(c)
	c.Infof("UpdateAnchorEventList 执行结束")
	return c.GetConsoleLog()
}

func GetActivityContestThemeVideo(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("GetActivityContestThemeVideo 开始执行")
	service.GetActivityContestThemeVideo(c)
	c.Infof("GetActivityContestThemeVideo 执行结束")
	return c.GetConsoleLog()
}
