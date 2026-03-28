package namespace

type NacosNamespace = string

const (
	ClearFinishedMatches         = "ClearFinishedMatches"
	SiteAdminUpdateLiveVideo     = "SiteAdminUpdateLiveVideo"
	PullLiveVideoEvents          = "PullLiveVideoEvents"
	PullLiveVideoOddsEvents      = "PullLiveVideoOddsEvents" // 拉取播控视频赔率数据
	SetLiveVideoEvents           = "SetLiveVideoEvents"
	PullAnchorEvents             = "PullAnchorEvents"
	UpdateAnchorEventList        = "UpdateAnchorEventList"
	GetActivityContestThemeVideo = "GetActivityContestThemeVideo"

	LSTYPullCountry     = "LSTYPullCountry"
	LSTYPullCategory    = "LSTYPullCategory"
	LSTYPullTeam        = "LSTYPullTeam"
	LSTYPullCompetition = "LSTYPullCompetition"
	LSTYPullMatch       = "LSTYPullMatch"
	LSTYPullLanguage    = "LSTYPullLanguage"
	LSTYPullVideoUrl    = "LSTYPullVideoUrl"
	CRTYPullMatch       = "CRTYPullMatch"
	CRTYSyncVideoUrl    = "CRTYSyncVideoUrl"
)

// Global为主站通用配置，其他项目可直接使用该配置
// Global要自动替换为app.dev.ini里面的Prefix
const (
	Application       NacosNamespace = "Application"
	Logger                           = "Logger"
	ControlSlave                     = "Global.Database.ControlSlave"
	Site                             = "Global.Database.Site"
	SiteSlave                        = "Global.Database.SiteSlave"
	Video                            = "Global.Database.Video"
	VideoSlave                       = "Global.Database.VideoSlave"
	XxlJobDbNamespace                = "Global.Database.XxlJob"
	Kafka                            = "Global.MQ.Kafka"
	Common                           = "Global.Config.Common"
	NameMapping                      = "Global.Mapping.NameMapping"
	RedisCore                        = "Global.Redis.RedisCore"
	RedisGame                        = "Global.Redis.RedisGame"
)
