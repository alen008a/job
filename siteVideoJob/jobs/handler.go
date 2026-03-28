package jobs

import (
	"siteVideoJob/jobs/crty"
	"siteVideoJob/jobs/lsty"
	"siteVideoJob/jobs/video"
	"siteVideoJob/mdata/namespace"

	"siteVideoJob/xxl"
)

const (
	jobPrefix     = "videoJob."
	lstyJobPrefix = "lstyJob."
)

// RegisterExecutors  注册执行器列表
func RegisterExecutors(execute *xxl.Executor) {
	execute.RegTask(jobPrefix+namespace.ClearFinishedMatches, video.ClearFinishedMatches)
	execute.RegTask(jobPrefix+namespace.SiteAdminUpdateLiveVideo, video.SiteAdminUpdateLiveVideo)
	execute.RegTask(jobPrefix+namespace.PullLiveVideoEvents, video.PullLiveVideoEvents)
	execute.RegTask(jobPrefix+namespace.SetLiveVideoEvents, video.SetLiveVideoEvents)
	execute.RegTask(jobPrefix+namespace.PullLiveVideoOddsEvents, video.PullLiveVideoOddsEvents)
	execute.RegTask(jobPrefix+namespace.PullAnchorEvents, video.PullAnchorEvents)
	execute.RegTask(jobPrefix+namespace.UpdateAnchorEventList, video.UpdateAnchorEventList)
	execute.RegTask(jobPrefix+namespace.GetActivityContestThemeVideo, video.GetActivityContestThemeVideo)
	execute.RegTask(jobPrefix+namespace.CRTYPullMatch, crty.PullCRTYMatch)
	execute.RegTask(jobPrefix+namespace.CRTYSyncVideoUrl, crty.DoSetCrtyFromLsty)

	execute.RegTask(lstyJobPrefix+namespace.LSTYPullCountry, lsty.PullLSTYCountry)
	execute.RegTask(lstyJobPrefix+namespace.LSTYPullCategory, lsty.PullLSTYCategory)
	execute.RegTask(lstyJobPrefix+namespace.LSTYPullTeam, lsty.PullLSTYTeam)
	execute.RegTask(lstyJobPrefix+namespace.LSTYPullCompetition, lsty.PullLSTYCompetition)
	execute.RegTask(lstyJobPrefix+namespace.LSTYPullMatch, lsty.PullLSTYMatch)
	execute.RegTask(lstyJobPrefix+namespace.LSTYPullLanguage, lsty.PullLSTYLanguage)
	execute.RegTask(lstyJobPrefix+namespace.LSTYPullVideoUrl, lsty.PullLSTYVideoUrl)

}
