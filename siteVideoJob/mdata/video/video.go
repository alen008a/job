package video

import (
	"encoding/json"
	"time"
)



type BkResp struct {
	StatusCode int        `json:"status_code"`
	Status     string     `json:"status"`
	Data       []ApiMatch `json:"data"`
}

type BkOddsResp struct {
	StatusCode int       `json:"status_code"`
	Status     string    `json:"status"`
	Data       []PlateVo `json:"data"`
}

type VideoUrlSourceResp struct {
	StatusCode int            `json:"status_code"`
	Status     string         `json:"status"`
	Data       VideoSourceMap `json:"data"`
}

type ApiMatch struct {
	AniID         int         `json:"ani_id"`
	VideoUrls     []*VideoUrl `json:"videoUrls"` //视频链接地址
	ChannelID     int         `json:"channel_id"`
	Channel       string      `json:"channel"`
	DisType       int         `json:"dis_type"`
	Eid           int64       `json:"eid"`
	League        string      `json:"league"`
	LeagueId      string      `json:"league_id"`
	LeagueLogo    string      `json:"league_logo"`
	Cate          int         `json:"cate"`    //默认0 表示体育场馆，1为电竞场馆
	PlateId       int         `json:"plateId"` //盘类ID 1今日 2早盘 3串关 4滚球
	BallClassId   int         `json:"ballClassId"`
	SportType     string      `json:"sport_type"`
	StartDatetime string      `json:"start_datetime"`
	StartTime     int         `json:"start_time"`
	Status        int         `json:"status"`
	StreamID      int64       `json:"stream_id"`
	Team1         string      `json:"team1"`
	Team1Logo     string      `json:"team1_logo"`
	Team2         string      `json:"team2"`
	Team2Logo     string      `json:"team2_logo"`
	Type          int         `json:"type"`
	AnchorStatus  string      `json:"anchor_status"`
	AnchorName    string      `json:"anchor_name"`
	Screenshot    string      `json:"screenshot"`
	Team1Score    string      `json:"team1_score"`
	Team2Score    string      `json:"team2_score"`
	Team1en       string      `json:"team1en"`
	Team2en       string      `json:"team2en"`
}

type VideoSourceMap map[string][]*VideoUrl

type VideoUrl struct {
	VideoType string `json:"videoType"` //视频或动画类型：1:Video 2:Animation
	Path      string `json:"path"`
	PlayType  string `json:"playType"` //地址类型 f-流地址 p-播放页面
}

type BkAnchorResp struct {
	Message    string                    `json:"message"`
	StatusCode int                       `json:"status_code"`
	Status     string                    `json:"status"`
	Data       map[string][]BkAnchorData //两种键名 not_started started 数据结构一致
}

type BkAnchorData struct {
	SportType string `json:"sport_type"`
	Events    []AnchorEvents
}

func (m AnchorEvents) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

type AnchorList []AnchorEvents

func (a AnchorList) Len() int {
	return len(a)
}

func (a AnchorList) Swap(i, j int) {
	{
		a[i], a[j] = a[j], a[i]
	}
}

func (a AnchorList) Less(i, j int) bool {
	return a[i].TimeStamp < a[j].TimeStamp
}

type AnchorEvents struct {
	Cate             string         `json:"cate"`
	Eid              string         `json:"eid"`
	IsStarted        string         `json:"is_started"`
	League           string         `json:"league"`
	LeagueID         string         `json:"league_id"`
	PlateId          int            `json:"plateId"` //盘类ID 1今日 2早盘 3串关 4滚球
	BallClassId      int            `json:"ballClassId"`
	Leagueen         string         `json:"leagueen"`
	StartDate        string         `json:"start_date"`
	Team1            string         `json:"team1"`
	Team1En          string         `json:"team1en"`
	Team2            string         `json:"team2"`
	Team2En          string         `json:"team2en"`
	GuestTeamLogoURL string         `json:"guestTeamLogoUrl"`
	HomeTeamLogoURL  string         `json:"homeTeamLogoUrl"`
	Team1Score       string         `json:"team1_score"`
	Team2Score       string         `json:"team2_score"`
	AnchorVideo      []AnchorVideos `json:"anchorVideo"`
	VenueName        string         `json:"venue_name"` //接口新增字段 用来区分场馆
	TimeStamp        int64          `json:"time_stamp"` //新增字段  便于排序
}

func (m AnchorEvents) Expired() bool {
	return (time.Now().Unix() - m.TimeStamp) > 24*3600
}

type AnchorVideos struct {
	Hd            string     `json:"HD"`
	Anchors       AnchorInfo `json:"anchor"`
	Online        string     `json:"online"`
	PathFlv       string     `json:"path_flv"`
	PathM3U8      string     `json:"path_m3u8"`
	Screenshot    string     `json:"screenshot"`
	MainVideoMode string     `json:"main_video_mode"`
	Status        string     `json:"status"`
	AnchorStatus  string     `json:"anchor_status"`
}

type AnchorInfo struct {
	Grade           string `json:"grade"`
	ID              string `json:"id"`
	LogoRectangle   string `json:"logo_rectangle"`
	LogoSquare      string `json:"logo_square"`
	Nickname        string `json:"nickname"`
	PersonalProfile string `json:"personal_profile"`
}

type AnchorListVo struct {
	VID          string //场馆名+EventID
	Nickname     string //主播昵称
	AnchorsID    string //主播ID
	AnchorStatus string //直播状态
	LogoSquare   string //主播logo
	EventID      string //赛事id
	Cate         string
	VenueName    string //场馆名
	League       string //联赛名
	StartDate    string //开始时间
	Team1        string //主队
	Team2        string //客队
	//UpdateAt     string //操作时间  动态字段放到单独key管理
	//UpdateBy     string //操作人
	//ShowStatus   int    //展示状态 默认不展示
}

type VideoSourceVo struct {
	VideoUrls        []*VideoUrl `json:"videoUrls"` //视频链接地址
	Cate             int         `json:"cate"`
	MatchClass       string      `json:"matchClass"`
	EventId          int64       `json:"eventId"`
	GuestTeamLogoUrl string      `json:"guestTeamLogoUrl"`
	HomeTeamLogoUrl  string      `json:"homeTeamLogoUrl"`
	League           string      `json:"league"`
	LeagueId         string      `json:"leagueId"`
	PlateId          int         `json:"plateId"` //盘类ID 1今日 2早盘 3串关 4滚球
	BallClassId      int         `json:"ballClassId"`
	LeagueLogoUrl    string      `json:"leagueLogoUrl"`
	StartDate        string      `json:"startDate"`
	StatusDes        string      `json:"statusDes"`
	StreamId         int64       `json:"streamId"`
	Team1            string      `json:"team1"`
	Team1en          string      `json:"team1en"`
	Team2            string      `json:"team2"`
	Team2en          string      `json:"team2en"`
	VenueName        string      `json:"venueName"`
	LiveStatus       int         `json:"liveStatus"` //直播状态 0 未开始 1进行中 2结束
	AnchorStatus     string      `json:"anchorStatus"`
	AnchorName       string      `json:"anchorName"`
}

type MatchDataVo struct {
	VideoUrls    []*VideoUrl `json:"videoUrls"` //视频链接地址
	LeagueId     string      `json:"leagueId"`
	LeagueLogo   string      `json:"league_logo"`
	MatchStatus  int         `json:"matchStatus"`
	HomeLogo     string      `json:"homeLogo"`
	VisitLogo    string      `json:"visitLogo"`
	StartTime    string      `json:"startTime"`
	AnchorStatus string      `json:"anchorStatus"`
	AnchorName   string      `json:"anchorName"`
	PlateId      int         `json:"plateId"` //盘类ID 1今日 2早盘 3串关 4滚球
	BallClassId  int         `json:"ballClassId"`
}

type VenueStatusInfo struct {
	Status    int    `json:"status"`
	IsDisplay int    `json:"isDisplay"`
	EnName    string `json:"enName"`
}

type LiveEventsDto struct {
	SiteId    int    `json:"siteId"`
	Status    int    `json:"status"`
	VenueName string `json:"venueName"`
	StartAt   string `json:"startAt"`
	DelStatus int    `json:"delStatus"`
}

type LiveEvents struct {
	Id                int64       `json:"id" gorm:"column:id"`
	SiteId            int         `json:"siteId" gorm:"column:site_id"`
	MatchId           int64       `json:"matchId" gorm:"column:match_id"`
	LeagueId          int64       `json:"leagueId" gorm:"column:league_id"`
	MatchName         string      `json:"matchName" gorm:"column:match_name"`
	MatchClass        string      `json:"matchClass" gorm:"column:match_class"`
	VideoUrls         []*VideoUrl `json:"videoUrls" gorm:"-"` //视频链接地址
	StartAt           string      `json:"startAt" gorm:"column:start_at"`
	HomeName          string      `json:"homeName" gorm:"column:home_name"`
	HomeLogo          string      `json:"homeLogo" gorm:"column:home_logo"`
	HomeScore         string      `json:"homeScore" gorm:"-"`
	VisitName         string      `json:"visitName" gorm:"column:visit_name"`
	VisitLogo         string      `json:"visitLogo" gorm:"column:visit_logo"`
	AwayScore         string      `json:"awayScore" gorm:"-"`
	Sort              int         `json:"sort" gorm:"column:sort"`
	SortByBallClassId int         `json:"sortByBallClassId" gorm:"-"`
	Status            int         `json:"status" gorm:"column:status"`
	DelFlag           int         `json:"delFlag" gorm:"column:del_flag"` // 是否删除 0未删除 1删除
	DisplayStartAt    string      `json:"displayStartAt" gorm:"column:display_start_at"`
	DisplayEndAt      string      `json:"displayEndAt" gorm:"column:display_end_at"`
	CreatedAt         string      `json:"createdAt" gorm:"column:created_at"`
	CreatedBy         string      `json:"createdBy" gorm:"column:created_by"`
	UpdatedAt         string      `json:"updatedAt" gorm:"column:updated_at"`
	UpdatedBy         string      `json:"updatedBy" gorm:"column:updated_by"`
	OpenTime          string      `json:"openTime" gorm:"-"`
	VenueName         string      `json:"venueName" gorm:"column:venue_name"`
	Cate              string      `json:"cate" gorm:"column:cate"` //默认0 表示体育场馆，1为电竞场馆
	GroupBy           string      `json:"groupBy" gorm:"-"`
	PlateId           int         `json:"plateId" gorm:"column:plate_id"`
	BallClassId       int         `json:"ballClassId" gorm:"column:ball_class_id"`
	MajorEvent        int         `json:"majorEvent" gorm:"column:major_event"`
	HomeSticky        int         `json:"homeSticky" gorm:"column:home_sticky"`
	OddsInfo          interface{} `json:"oddsInfo" gorm:"-"`
	PeriodStatus      string      `json:"periodStatus" gorm:"-"`
	IsStart           bool        `json:"isStart" gorm:"-"`
	Week              string      `json:"week" gorm:"-"`
	GameTime          string      `json:"gameTime" gorm:"-"`
	LiveStatus        int         `json:"liveStatus" gorm:"column:live_status"`
	MatchLogo         string      `json:"leagueLogo" gorm:"column:league_logo"`
	AnchorStatus      string      `json:"anchorStatus" gorm:"-"`
	AnchorName        string      `json:"anchorName" gorm:"-"`
}

type LiveEventList []LiveEvents

func (s LiveEventList) Len() int      { return len(s) }
func (s LiveEventList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s LiveEventList) Less(i, j int) bool {
	ei := s[i]
	ej := s[j]
	if ei.MajorEvent > 0 || ej.MajorEvent > 0 {
		return ei.MajorEvent > ej.MajorEvent
	}
	if ei.Sort > 0 || ej.Sort > 0 {
		return ei.Sort > ej.Sort
	}
	return ei.StartAt < ej.StartAt
}

type PlateVo struct {
	Gid          string      `json:"gid"`       //赛事id
	League       string      `json:"league"`    //联赛名称
	Kid          string      `json:"kid"`       //联赛ID
	P            string      `json:"p"`         //赛事时间(数字) 0全部 1今日 2早盘 3串关 4滚球
	CateId       string      `json:"t"`         //赛种ID
	HomeName     string      `json:"h"`         //主队名称
	AwayName     string      `json:"a"`         //客队名称
	Platform     string      `json:"platform"`  //场馆名称
	Category     string      `json:"sportName"` //球类类别
	HomeScore    string      `json:"homeScore"`
	AwayScore    string      `json:"awayScore"`
	LiveStatus   string      `json:"liveStatus"` // 赛事状态 0 未开始 1进行中 2结束
	GameTime     string      `json:"gameTime"`
	OpenTime     string      `json:"openTime"` // 开赛的开始时间
	PeriodStatus string      `json:"periodStatus"`
	Odds         interface{} `json:"odds"`
}

type BkResponseForSourceUrl struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Data       ResponseSourceUrl
}

type ResponseSourceUrl struct {
	Status       int                 `json:"status"`
	URL          SourceUrlDetail     `json:"url"`
	List         []SourceUrlDetail   `json:"list"`
	AnimationURL string              `json:"animation_url"`
	AnchorM3U8   string              `json:"anchor_m3u8_url"`
	AnchorFlv    string              `json:"anchor_flv_url"`
	AnchorStatus string              `json:"anchor_status"`
	AnchorName   string              `json:"anchor_name"`
	SwitchBTN    int                 `json:"switchbtn"`
	AnchorList   []map[string]string `json:"anchor_list"`
	PlateVo
}
type SourceUrlDetail struct {
	FlvURL  string      `json:"flv_url"`
	PlayURL string      `json:"play_url"`
	Hd      interface{} `json:"hd"`
	Name    string      `json:"name"`
}

type LoginData struct {
	UserId      string `json:"userId"`
	Token       string `json:"token"`
	ApiDomain   string `json:"apiDomain"`
	Domain      string `json:"domain"`
	ImgDomain   string `json:"imgDomain"`
	EsImgDomain string `json:"esImgDomain"`
}
