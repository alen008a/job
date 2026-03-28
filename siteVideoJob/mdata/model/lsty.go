package model

type DataType int

const (
	CategoryDataType    = 1 //类型
	CountryDataType     = 2 //国家
	CompetitionDataType = 3 //联赛
	TeamDataType        = 4 //队伍
)

// CommonReq 公共请求的查询参数
type CommonReq struct {
	Page int    `json:"page"  default:"1000"` // 页码查询，返回查询页码的数据（从1开始），默认是1000
	Time int    `json:"time"`                 // 时间查询，返回更新时间大于等于该时间戳的记录，按更新时间排序
	UUID string `json:"uuid"`                 // UUID 查询，返回查询的 UUID 数据
	Type int    `json:"type"`                 // 1-category, 2-country, 3-competition, 4-team, 5-player, 6-injury(Team injury)
}

// Inquiry 用于表示查询结果
type Inquiry struct {
	Total   int    `json:"total"`              // 返回的数据总量
	Type    string `json:"type"`               // 查询类型，uuid 查询：uuid，page 查询：page，time 查询：time，默认为 page
	UUID    string `json:"uuid,omitempty"`     // UUID 查询值（UUID 查询，字段存在时）
	Page    int    `json:"page,omitempty"`     // 页码查询值（页码查询，字段存在时）
	Time    int    `json:"time,omitempty"`     // 时间查询值，时间戳格式（时间查询，字段存在时）
	MinTime int    `json:"min_time,omitempty"` // 返回数据中的最小时间（updated_at 值）（时间查询，字段存在时）
	MaxTime int    `json:"max_time,omitempty"` // 返回数据中的最大时间（updated_at 值）（时间查询，字段存在时）
}

// Team 球队
type Team struct {
	ID                  string `json:"id"`                    // 球队ID
	CompetitionID       string `json:"competition_id"`        // 赛事ID（所属联赛，杯赛无关）
	CountryID           string `json:"country_id"`            // 国家ID
	Name                string `json:"name"`                  // 球队名称
	ShortName           string `json:"short_name"`            // 球队简称
	Logo                string `json:"logo"`                  // 球队Logo
	National            int    `json:"national"`              // 是否是国家队，1-是，0-否
	CountryLogo         string `json:"country_logo"`          // 国家队Logo（仅国家队存在）
	FoundationTime      int    `json:"foundation_time"`       // 成立时间
	Website             string `json:"website"`               // 官方网站
	CoachID             string `json:"coach_id"`              // 教练ID
	VenueID             string `json:"venue_id"`              // 场地ID
	MarketValue         int    `json:"market_value"`          // 市场价值
	MarketValueCurrency string `json:"market_value_currency"` // 市场价值单位
	TotalPlayers        int    `json:"total_players"`         // 球员总数，-1表示没有数据
	ForeignPlayers      int    `json:"foreign_players"`       // 外籍球员数，-1表示没有数据
	NationalPlayers     int    `json:"national_players"`      // 国家队球员数，-1表示没有数据
	UID                 string `json:"uid"`                   // 统一的球队ID（合并重复球队后对应的ID）
	UpdatedAt           int    `json:"updated_at"`            // 更新时间（时间戳）
}

// TeamResp 队伍响应
type TeamResp struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Data       struct {
		Query   Inquiry `json:"query"`
		Results []*Team `json:"results"`
	} `json:"data"`
}

// MatchTeamData 队伍表
type MatchTeamData struct {
	ID               int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`     // 自增主键id
	TeamID           string `gorm:"column:team_id" json:"teamId"`                     // 队伍id
	UID              string `gorm:"column:uid" json:"uid"`                            // 队伍id（重复队伍合并后对应的id），如果存在则返回
	CountryID        string `gorm:"column:country_id" json:"countryId"`               // 国家ID
	Name             string `gorm:"column:name" json:"name"`                          // 队伍名称
	MatchClass       string `gorm:"column:match_class" json:"matchClass"`             // 赛事分类
	ShortName        string `gorm:"column:short_name" json:"shortName"`               // 队伍简称
	NameZh           string `gorm:"column:name_zh" json:"nameZh"`                     // 队伍简体中文
	NameZht          string `gorm:"column:name_zht" json:"nameZht"`                   // 队伍繁体中文
	UpdatedTimestamp int    `gorm:"column:updated_timestamp" json:"updatedTimestamp"` // 三方更新时间
	Logo             string `gorm:"column:logo" json:"logo"`                          // 队伍logo
	National         int    `gorm:"column:national" json:"national"`                  // 是否是国家队 1-是 0-否
	CountryLogo      string `gorm:"column:country_logo" json:"countryLogo"`           // 国家队LOGO, national=1时存在
	CreatedAt        string `gorm:"column:created_at;<-:false" json:"createdAt"`      // 入库时间
	CreatedBy        string `gorm:"column:created_by" json:"createdBy"`               // 创建人
	UpdatedAt        string `gorm:"column:updated_at;<-:false" json:"updatedAt"`      // 修改时间
	UpdatedBy        string `gorm:"column:updated_by" json:"updatedBy"`               // 修改人
}

// TableName 设置数据库表名称
func (MatchTeamData) TableName() string {
	return "match_team_data"
}

type VideoUrlResp struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Data       []*VideoUrl `json:"data"`
}
type VideoUrl struct {
	SportId   int    `json:"sport_id"` //Balls, 1-football, 2-basketball
	MatchId   string `json:"match_id"`
	MatchTime int    `json:"match_time"`
	Pushurl1  string `json:"pushurl1"`
	Pushurl2  string `json:"pushurl2"`
}

// MatchVideoData 赛事视频表
type MatchVideoData struct {
	ID               int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`     // 自增主键id
	MatchID          string `gorm:"column:match_id" json:"match_id"`                  // 赛事id
	MatchClass       string `gorm:"column:match_class" json:"match_class"`            // 赛事分类
	UpdatedTimestamp int    `gorm:"column:updated_timestamp" json:"updatedTimestamp"` // 三方更新时间
	StartAt          string `gorm:"column:start_at" json:"start_at"`                  // 开赛日期
	PushUrl1         string `gorm:"column:push_url1" json:"pushurl1"`                 // SD stream address
	PushUrl2         string `gorm:"column:push_url2" json:"pushurl2"`                 // English HD stream address, not empty
	CreatedAt        string `gorm:"column:created_at;<-:false" json:"created_at"`     // 入库时间
	CreatedBy        string `gorm:"column:created_by" json:"created_by"`              // 创建人
	UpdatedAt        string `gorm:"column:updated_at;<-:false" json:"updated_at"`     // 修改时间
	UpdatedBy        string `gorm:"column:updated_by" json:"updated_by"`              // 修改人
}

// TableName 设置数据库表名称
func (MatchVideoData) TableName() string {
	return "match_video_data"
}

type CountryResp struct {
	StatusCode int        `json:"status_code"`
	Status     string     `json:"status"`
	Data       []*Country `json:"data"`
}

// Country 国家
type Country struct {
	ID         string `json:"id"`          // 国家ID
	CategoryID string `json:"category_id"` // 分类ID
	Name       string `json:"name"`        // 国家名称
	Logo       string `json:"logo"`        // 国家Logo
	UpdatedAt  int    `json:"updated_at"`  // 更新时间（时间戳）
}

// MatchCountryData 国家表
type MatchCountryData struct {
	ID               int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`     // 自增主键id
	CountryID        string `gorm:"column:country_id" json:"countryId"`               // 国家ID
	CategoryID       string `gorm:"column:category_id" json:"categoryId"`             // 类别ID
	MatchClass       string `gorm:"column:match_class" json:"matchClass"`             // 赛事分类
	Name             string `gorm:"column:name" json:"name"`                          // 国家名称
	NameZh           string `gorm:"column:name_zh" json:"nameZh"`                     // 队伍简体中文
	NameZht          string `gorm:"column:name_zht" json:"nameZht"`                   // 队伍繁体中文
	UpdatedTimestamp int    `gorm:"column:updated_timestamp" json:"updatedTimestamp"` // 三方更新时间
	Logo             string `gorm:"column:logo" json:"logo"`                          // 国家logo
	CreatedAt        string `gorm:"column:created_at;<-:false" json:"createdAt"`      // 入库时间
	CreatedBy        string `gorm:"column:created_by" json:"createdBy"`               // 创建人
	UpdatedAt        string `gorm:"column:updated_at;<-:false" json:"updatedAt"`      // 修改时间
	UpdatedBy        string `gorm:"column:updated_by" json:"updatedBy"`               // 修改人
}

// TableName 设置数据库表名称
func (MatchCountryData) TableName() string {
	return "match_country_data"
}

type CategoryResp struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Data       []*Category `json:"data"`
}

type Category struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	UpdatedAt int    `json:"updated_at"`
}

// MatchCategoryData 类别表
type MatchCategoryData struct {
	ID               int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`     // 自增主键id
	CategoryID       string `gorm:"column:category_id" json:"categoryId"`             // 类别ID
	MatchClass       string `gorm:"column:match_class" json:"matchClass"`             // 赛事分类
	Name             string `gorm:"column:name" json:"name"`                          // 类别名称
	NameZh           string `gorm:"column:name_zh" json:"nameZh"`                     // 队伍简体中文
	NameZht          string `gorm:"column:name_zht" json:"nameZht"`                   // 队伍繁体中文
	UpdatedTimestamp int    `gorm:"column:updated_timestamp" json:"updatedTimestamp"` // 三方更新时间
	CreatedAt        string `gorm:"column:created_at;<-:false" json:"createdAt"`      // 入库时间
	CreatedBy        string `gorm:"column:created_by" json:"createdBy"`               // 创建人
	UpdatedAt        string `gorm:"column:updated_at;<-:false" json:"updatedAt"`      // 修改时间
	UpdatedBy        string `gorm:"column:updated_by" json:"updatedBy"`               // 修改人
}

// TableName 设置数据库表名称
func (MatchCategoryData) TableName() string {
	return "match_category_data"
}

// CompetitionResp 赛事响应
type CompetitionResp struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Data       struct {
		Query   Inquiry        `json:"query"`
		Results []*Competition `json:"results"`
	} `json:"data"`
}

// Competition 比赛数据
type Competition struct {
	ID             string `json:"id"`              // 赛事ID
	CategoryID     string `json:"category_id"`     // 分类ID
	CountryID      string `json:"country_id"`      // 国家ID
	Name           string `json:"name"`            // 赛事名称
	ShortName      string `json:"short_name"`      // 赛事简称
	Logo           string `json:"logo"`            // 赛事Logo
	Type           int    `json:"type"`            // 赛事类型，0-未知，1-联赛，2-杯赛，3-友谊赛
	CurSeasonID    string `json:"cur_season_id"`   // 当前赛季ID
	CurStageID     string `json:"cur_stage_id"`    // 当前阶段ID
	CurRound       int    `json:"cur_round"`       // 当前轮次
	RoundCount     int    `json:"round_count"`     // 总轮数
	PrimaryColor   string `json:"primary_color"`   // 主色调，可忽略
	SecondaryColor string `json:"secondary_color"` // 次色调，可忽略
	UpdatedAt      int    `json:"updated_at"`      // 更新时间（时间戳）
}

// MatchCompetitionData 联赛数据表
type MatchCompetitionData struct {
	ID               string `gorm:"primaryKey;autoIncrement;column:id" json:"id"`     // 自增主键id
	CompetitionID    string `gorm:"column:competition_id" json:"competitionId"`       // 联赛ID
	CategoryID       string `gorm:"column:category_id" json:"categoryId"`             // 分类ID
	MatchClass       string `gorm:"column:match_class" json:"matchClass"`             // 赛事分类
	CountryID        string `gorm:"column:country_id" json:"countryId"`               // 国家ID
	Name             string `gorm:"column:name" json:"name"`                          // 赛事名称
	ShortName        string `gorm:"column:short_name" json:"shortName"`               // 赛事简称
	NameZh           string `gorm:"column:name_zh" json:"nameZh"`                     // 队伍简体中文
	NameZht          string `gorm:"column:name_zht" json:"nameZht"`                   // 队伍繁体中文
	Logo             string `gorm:"column:logo" json:"logo"`                          // 赛事Logo
	Type             int    `gorm:"column:type"  json:"type"`                         // 赛事类型，0-未知，1-联赛，2-杯赛，3-友谊赛
	CurSeasonID      string `gorm:"column:cur_season_id" json:"curSeasonId"`          // 当前赛季ID
	CurStageID       string `gorm:"column:cur_stage_id" json:"curStageId"`            // 当前阶段ID
	CurRound         int    `gorm:"column:cur_round" json:"curRound"`                 // 当前轮次
	RoundCount       int    `gorm:"column:round_count" json:"roundCount"`             // 总轮数
	UpdatedTimestamp int    `gorm:"column:updated_timestamp" json:"updatedTimestamp"` // 三方更新时间
	CreatedAt        string `gorm:"column:created_at;<-:false" json:"created_at"`     // 入库时间
	CreatedBy        string `gorm:"column:created_by" json:"created_by"`              // 创建人
	UpdatedAt        string `gorm:"column:updated_at;<-:false" json:"updated_at"`     // 修改时间
	UpdatedBy        string `gorm:"column:updated_by" json:"updated_by"`              // 修改人
}

// TableName 设置数据库表名称
func (MatchCompetitionData) TableName() string {
	return "match_competition_data"
}

// MatchResp 赛事响应
type MatchResp struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Data       struct {
		Query   Inquiry  `json:"query"`
		Results []*Match `json:"results"`
	} `json:"data"`
}

// Match 赛事
type Match struct {
	ID            string   `json:"id"`             // 比赛ID
	SeasonID      string   `json:"season_id"`      // 赛季ID
	CompetitionID string   `json:"competition_id"` // 赛事ID
	HomeTeamID    string   `json:"home_team_id"`   // 主队ID
	AwayTeamID    string   `json:"away_team_id"`   // 客队ID
	StatusID      int      `json:"status_id"`      // 比赛状态，请参考状态代码 -> 比赛状态
	MatchTime     int      `json:"match_time"`     // 比赛时间
	VenueID       string   `json:"venue_id"`       // 场地ID
	RefereeID     string   `json:"referee_id"`     // 裁判ID
	Neutral       int      `json:"neutral"`        // 是否为中立场，1-是，0-否
	LiveStatus    int      `json:"liveStatus"`     // 直播状态 0 未开始 1进行中 2结束
	Note          string   `json:"note"`           // 备注
	HomeScores    []int    `json:"home_scores"`    // 主队比分字段描述
	AwayScores    []int    `json:"away_scores"`    // 客队比分字段描述
	HomePosition  string   `json:"home_position"`  // 主队排名
	AwayPosition  string   `json:"away_position"`  // 客队排名
	Coverage      Coverage `json:"coverage"`       // 动画、阵容
	Round         Round    `json:"round"`          // 阶段
	RelatedID     string   `json:"related_id"`     // 双回合中另一回合的比赛ID（无数据字段不存在）
	AggScore      []int    `json:"agg_score"`      // 常规时间两回合总比分（含加时，若无数据字段不存在）
	UpdatedAt     int      `json:"updated_at"`     // 更新时间
}

// Coverage 表示动画和阵容信息
type Coverage struct {
	MLive  int `json:"mlive"`  // 是否有动画，1-有，0-无
	Lineup int `json:"lineup"` // 是否有阵容，1-有，0-无
}

// Round 表示阶段信息
type Round struct {
	StageID  string `json:"stage_id"`  // 阶段ID
	GroupNum int    `json:"group_num"` // 所属小组，1-A，2-B，依次类推
	RoundNum int    `json:"round_num"` // 第几轮
}

// MatchOriginData 赛事视频源直播表
type MatchOriginData struct {
	ID               int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`     // 自增主键id
	MatchID          string `gorm:"column:match_id" json:"match_id"`                  // 赛事id
	MatchName        string `gorm:"column:match_name" json:"match_name"`              // 赛事名称
	MatchClass       string `gorm:"column:match_class" json:"match_class"`            // 赛事分类
	StatusCode       int    `json:"column:status_code" json:"statusCode"`             // 比赛状态，请参考状态代码 -> 比赛状态
	HomeId           string `gorm:"column:home_id" json:"homeId"`                     // 主队Id
	HomeName         string `gorm:"column:home_name" json:"home_name"`                // 主队名称
	HomeLogo         string `gorm:"column:home_logo" json:"home_logo"`                // 主队队标
	VisitId          string `gorm:"column:visit_id" json:"visitId"`                   // 客队Id
	VisitName        string `gorm:"column:visit_name" json:"visit_name"`              // 客队名称
	VisitLogo        string `gorm:"column:visit_logo" json:"visit_logo"`              // 客队队标
	VenueName        string `gorm:"column:venue_name" json:"venue_name"`              // 场馆名称
	DataSource       int    `gorm:"column:data_source" json:"data_source"`            // 数据来源 0-拉取三方数据 1-后台手动添加
	LeagueID         string `gorm:"column:league_id" json:"league_id"`                // 联赛id
	LeagueLogo       string `gorm:"column:league_logo" json:"league_logo"`            // 联赛logo
	PlateID          int    `gorm:"column:plate_id" json:"plate_id"`                  // 盘类id  1今日 2早盘 3串关 4滚球 5闭盘
	BallClassID      int    `gorm:"column:ball_class_id" json:"ball_class_id"`        // 球类id
	LiveStatus       int    `gorm:"column:live_status" json:"live_status"`            // 直播状态 0 未开始 1进行中 2结束
	MLive            int    `gorm:"column:m_live" json:"mlive"`                       // 是否有动画，1-有，0-无
	Lineup           int    `gorm:"column:line_up" json:"lineup"`                     // 是否有阵容，1-有，0-无
	StageCode        string `gorm:"column:stage_code" json:"stageCode"`               // 阶段ID
	GroupNum         int    `gorm:"column:group_num" json:"groupNum"`                 // 所属小组，1-A，2-B，依次类推
	RoundNum         int    `gorm:"column:round_num" json:"roundNum"`                 // 第几轮
	UpdatedTimestamp int    `gorm:"column:updated_timestamp" json:"updatedTimestamp"` // 三方更新时间
	StartAt          string `gorm:"column:start_at" json:"start_at"`                  // 开赛日期
	CreatedAt        string `gorm:"column:created_at;<-:false" json:"created_at"`     // 入库时间
	CreatedBy        string `gorm:"column:created_by" json:"created_by"`              // 创建人
	UpdatedAt        string `gorm:"column:updated_at;<-:false" json:"updated_at"`     // 修改时间
	UpdatedBy        string `gorm:"column:updated_by" json:"updated_by"`              // 修改人
}

// TableName 设置数据库表名称
func (MatchOriginData) TableName() string {
	return "match_origin_data"
}

type LanguageResp struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Data       struct {
		Query   Inquiry     `json:"query"`
		Results []*Language `json:"results"`
	} `json:"data"`
}

// Language 语言
type Language struct {
	Id        string `json:"id"`
	NameEn    string `json:"name_en"`    //英文
	NameZh    string `json:"name_zh"`    //简体中文
	NameZht   string `json:"name_zht"`   //繁体中文
	UpdatedAt int    `json:"updated_at"` //更新时间
}

// MatchLanguageData 赛事多语言表
type MatchLanguageData struct {
	ID               int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`     // 自增主键id
	DataId           string `gorm:"column:data_id" json:"dataId"`                     // 对应类型的数据ID
	MatchClass       string `gorm:"column:match_class" json:"match_class"`            // 赛事分类
	DataType         int    `gorm:"column:data_type" json:"dataType"`                 // 语言数据类型：1-category, 2-country, 3-competition, 4-team
	NameEn           string `gorm:"column:name_en" json:"nameEn"`                     // 英文
	NameZh           string `gorm:"column:name_zh" json:"nameZh"`                     // 简体中文
	NameZht          string `gorm:"column:name_zht" json:"nameZht"`                   // 繁体中文
	UpdatedTimestamp int    `gorm:"column:updated_timestamp" json:"updatedTimestamp"` // 三方更新时间
	CreatedAt        string `gorm:"column:created_at;<-:false" json:"created_at"`     // 入库时间
	CreatedBy        string `gorm:"column:created_by" json:"created_by"`              // 创建人
	UpdatedAt        string `gorm:"column:updated_at;<-:false" json:"updated_at"`     // 修改时间
	UpdatedBy        string `gorm:"column:updated_by" json:"updated_by"`              // 修改人
}

// TableName 设置数据库表名称
func (MatchLanguageData) TableName() string {
	return "match_language_data"
}

type MatchVideoUrlData struct {
	ID          int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"` // 自增主键id
	MatchID     string `gorm:"column:match_id" json:"match_id"`              // 赛事id
	MatchClass  string `gorm:"column:match_class" json:"match_class"`        // 赛事分类
	StartAt     string `gorm:"column:start_at" json:"start_at"`              // 开赛日期
	PushUrl1    string `gorm:"column:push_url1" json:"pushurl1"`             // SD stream address
	PushUrl2    string `gorm:"column:push_url2" json:"pushurl2"`             // English HD stream address, not empty
	MatchName   string `gorm:"column:match_name" json:"match_name"`          // 赛事名称
	HomeName    string `gorm:"column:home_name" json:"home_name"`            // 主队名称
	HomeLogo    string `gorm:"column:home_logo" json:"home_logo"`            // 主队队标
	VisitName   string `gorm:"column:visit_name" json:"visit_name"`          // 客队名称
	VisitLogo   string `gorm:"column:visit_logo" json:"visit_logo"`          // 客队队标
	VenueName   string `gorm:"column:venue_name" json:"venue_name"`          // 场馆名称
	LeagueLogo  string `gorm:"column:league_logo" json:"league_logo"`        // 联赛logo
	PlateID     int    `gorm:"column:plate_id" json:"plate_id"`              // 盘类id  1今日 2早盘 3串关 4滚球 5闭盘
	BallClassID int    `gorm:"column:ball_class_id" json:"ball_class_id"`    // 球类id
}
