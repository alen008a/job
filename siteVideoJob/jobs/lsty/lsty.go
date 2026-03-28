package lsty

import (
	"siteVideoJob/internal/context"
	"siteVideoJob/mdata"
	"siteVideoJob/service"
	"siteVideoJob/xxl"
	"strconv"
)

// PullLSTYCountry 拉取雷速国家列表
func PullLSTYCountry(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams *xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullLSTYCountry 开始执行")
	msg = service.PullLSTYCountryList(c)
	c.Infof("PullLSTYCountry 执行结束 result:%s", msg)
	c.Console("PullLSTYCountry 执行结束 result:%s", msg)
	return c.GetConsoleLog()
}

// PullLSTYCategory 拉取雷速分类列表
func PullLSTYCategory(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams *xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullLSTYCategory 开始执行")
	service.PullLSTYCategoryList(c)
	c.Infof("PullLSTYCategory 执行结束 result:%s", msg)
	c.Console("PullLSTYCategory 执行结束 result:%s", msg)
	return c.GetConsoleLog()
}

// PullLSTYTeam 拉取雷速队伍列表
func PullLSTYTeam(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams *xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullLSTYTeam 开始执行")
	msg = service.PullLSTYTeamList(c, executorParams)
	c.Infof("PullLSTYTeam 执行结束 result:%s", msg)
	c.Console("PullLSTYTeam 执行结束 result:%s", msg)
	return c.GetConsoleLog()
}

// PullLSTYCompetition 拉取雷速联赛列表
func PullLSTYCompetition(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams *xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullLSTYCompetition 开始执行")
	msg = service.PullLSTYCompetitionList(c, executorParams)
	c.Infof("PullLSTYCompetition 执行结束 result:%s", msg)
	c.Console("PullLSTYCompetition 执行结束 result:%s", msg)
	return c.GetConsoleLog()
}

// PullLSTYMatch 拉取雷速比赛列表
func PullLSTYMatch(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams *xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullLSTYMatch 开始执行")
	msg = service.PullLSTYMatchList(c, executorParams)
	c.Infof("PullLSTYMatch 执行结束 result:%s", msg)
	c.Console("PullLSTYMatch 执行结束 result:%s", msg)
	return c.GetConsoleLog()
}

// PullLSTYLanguage 拉取雷速语言列表
func PullLSTYLanguage(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams *xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullLSTYLanguage 开始执行")
	msg = service.PullLSTYLanguageList(c, executorParams)
	c.Infof("PullLSTYLanguage 执行结束 result:%s", msg)
	c.Console("PullLSTYLanguage 执行结束 result:%s", msg)
	return c.GetConsoleLog()
}

// PullLSTYVideoUrl 拉取雷速视频源列表
func PullLSTYVideoUrl(c *context.Context, param *xxl.RunReq) (msg string) {
	var executorParams *xxl.ExecutorParams
	mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
	c.SiteId = strconv.Itoa(executorParams.SiteId)
	c.Infof("PullLSTYVideoUrl 开始执行")
	msg = service.PullLSTYVideoUrlList(c)
	c.Infof("PullLSTYVideoUrl 执行结束 result:%s", msg)
	c.Console("PullLSTYVideoUrl 执行结束 result:%s", msg)
	return c.GetConsoleLog()
}
