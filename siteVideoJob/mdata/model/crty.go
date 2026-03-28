package model

type PullAllGameReq struct {
	GameDate string `json:"gameDate"` //日期格式
}

type MatchCrtyVideoData struct {
	MatchId    string `gorm:"column:match_id" json:"matchId"`       // 赛事id
	LeagueId   string `gorm:"column:league_id" json:"leagueId"`     // 联赛id
	MatchClass string `gorm:"column:match_class" json:"matchClass"` // 赛事分类
	StartAt    string `gorm:"column:start_at" json:"startAt"`       // 开赛日期
	HomeLogo   string `gorm:"column:home_logo" json:"homeLogo"`     // 主队队标
	VisitLogo  string `gorm:"column:visit_logo" json:"visitLogo"`   // 客队队标
	LeagueLogo string `gorm:"column:league_logo" json:"leagueLogo"` // 联赛图标
	PushUrl1   string `gorm:"column:push_url1" json:"pushurl1"`     // SD stream address
	PushUrl2   string `gorm:"column:push_url2" json:"pushurl2"`     // English HD stream address, not empty
}

// MatchCrtyOriginData CRTY体育赛事源表
type MatchCrtyOriginData struct {
	Id         int    `gorm:"column:id" json:"id"`                         // 自增主键id
	MatchId    string `gorm:"column:match_id" json:"matchId"`              // 赛事id
	LeagueId   string `gorm:"column:league_id" json:"leagueId"`            // 联赛id
	MatchName  string `gorm:"column:match_name" json:"matchName"`          // 赛事名称
	VenueName  string `gorm:"column:venue_name" json:"venueName"`          // 场馆名称
	MatchClass string `gorm:"column:match_class" json:"matchClass"`        // 赛事分类
	StartAt    string `gorm:"column:start_at" json:"startAt"`              // 开赛日期
	HomeName   string `gorm:"column:home_name" json:"homeName"`            // 主队名称
	HomeLogo   string `gorm:"column:home_logo" json:"homeLogo"`            // 主队队标
	VisitName  string `gorm:"column:visit_name" json:"visitName"`          // 客队名称
	VisitLogo  string `gorm:"column:visit_logo" json:"visitLogo"`          // 客队队标
	LeagueLogo string `gorm:"column:league_logo" json:"leagueLogo"`        // 联赛图标
	LiveStatus int    `gorm:"column:live_status" json:"liveStatus"`        // 直播状态 0 未开始 1进行中 2结束
	CreatedAt  string `gorm:"column:created_at;<-:false" json:"createdAt"` // 入库时间
	CreatedBy  string `gorm:"column:created_by" json:"createdBy"`          // 创建人
	UpdatedAt  string `gorm:"column:updated_at;<-:false" json:"updatedAt"` // 修改人
	UpdatedBy  string `gorm:"column:updated_by" json:"updatedBy"`          // 修改时间
}

// TableName 设置数据库表名称
func (MatchCrtyOriginData) TableName() string {
	return "match_crty_origin_data"
}
