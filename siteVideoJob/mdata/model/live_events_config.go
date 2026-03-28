package model

type LiveEventsConfig struct {
	Id             int    `gorm:"column:id" json:"id"` // 自增主键id
	SiteId         int    `gorm:"column:site_id" json:"siteId"`
	MatchId        int64  `gorm:"column:match_id" json:"matchId"`                // 赛事id
	MatchName      string `gorm:"column:match_name" json:"matchName"`            // 赛事名称
	MatchClass     string `gorm:"column:match_class" json:"matchClass"`          // 赛事分类
	StartAt        string `gorm:"column:start_at" json:"startAt"`                // 开赛日期
	HomeName       string `gorm:"column:home_name" json:"homeName"`              // 主队名称
	HomeLogo       string `gorm:"column:home_logo" json:"homeLogo"`              // 主队队标
	VisitName      string `gorm:"column:visit_name" json:"visitName"`            // 客队名称
	VisitLogo      string `gorm:"column:visit_logo" json:"visitLogo"`            // 客队队标
	Sort           int    `gorm:"column:sort" json:"sort"`                       // 排序
	Status         int    `gorm:"column:status" json:"status"`                   // 状态 0启用 1停用
	DisplayStartAt string `gorm:"column:display_start_at" json:"displayStartAt"` // 展示开始时间
	DisplayEndAt   string `gorm:"column:display_end_at" json:"displayEndAt"`     // 展示结束时间
	DelFlag        int    `gorm:"column:del_flag" json:"delFlag"`                // 是否删除 0未删除 1删除
	CreatedAt      string `gorm:"column:created_at;<-:false" json:"createdAt"`   // 入库时间
	CreatedBy      string `gorm:"column:created_by" json:"createdBy"`            // 创建人
	UpdatedAt      string `gorm:"column:updated_at;<-:false" json:"updatedAt"`   // 修改人
	UpdatedBy      string `gorm:"column:updated_by" json:"updatedBy"`            // 修改时间
	VenueName      string `gorm:"column:venue_name" json:"venueName"`            // 场馆名称
	DataSource     int    `gorm:"column:data_source" json:"dataSource"`          // 数据来源 0-拉取三方数据 1-后台手动添加
	Cate           string `gorm:"column:cate" json:"cate"`                       // 赛事类型
	LeagueId       int    `gorm:"column:league_id" json:"leagueId"`              // 联赛id
	LeagueLogo     string `gorm:"column:league_logo" json:"leagueLogo"`          // 联赛图标
	PlateId        int    `gorm:"column:plate_id" json:"plateId"`                // 盘类id
	BallClassId    int    `gorm:"column:ball_class_id" json:"ballClassId"`       // 球类id
	MatchNameBy    string `gorm:"column:match_name_by" json:"matchNameBy"`       // 赛事名称操作人
	MatchClassBy   string `gorm:"column:match_class_by" json:"matchClassBy"`     // 赛事分类操作人
	StartAtBy      string `gorm:"column:start_at_by" json:"startAtBy"`           // 开赛日期操作人
	HomeNameBy     string `gorm:"column:home_name_by" json:"homeNameBy"`         // 主队名称操作人
	HomeLogoBy     string `gorm:"column:home_logo_by" json:"homeLogoBy"`         // 主队队标操作人
	VisitNameBy    string `gorm:"column:visit_name_by" json:"visitNameBy"`       // 客队名称操作人
	VisitLogoBy    string `gorm:"column:visit_logo_by" json:"visitLogoBy"`       // 客队队标操作人
	HomeSticky     int    `gorm:"column:home_sticky" json:"homeSticky"`          // 是否首页置顶 1置顶 0不置顶
	MajorEvent     int    `gorm:"column:major_event" json:"majorEvent"`          // 是否重要赛事 1是 0不是
	LiveStatus     int    `gorm:"column:live_status" json:"liveStatus"`          // 直播状态 0 未开始 1进行中 2结束
}

func (*LiveEventsConfig) TableName() string {
	return "live_events_config"
}
