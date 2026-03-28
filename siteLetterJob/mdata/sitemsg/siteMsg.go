package sitemsg

// MsgDto 站内信发送消息实体
type MsgDto struct {
	BatchId       string      `json:"batchId"`       // 批次号
	OperationType string      `json:"operationType"` // 操作类型1-新增，2-修改，3-删除，4-启用,5-停用，6-置顶，7-取消置顶
	Status        int         `json:"status"`        // 是否有效0-有效，1无效
	Id            string      `json:"id"`            // 文档id
	MemberId      int         `json:"memberId"`      // 会员id
	SiteId        int         `json:"siteId"`        // 站点id
	GroupType     int         `json:"groupType"`     // 发送群组类型 1指定用户（groupId可为空），2在线用户，3全部用户，4VIP等级用户，5场馆用户
	GroupId       string      `json:"groupId"`       // 群组id,多个用逗号隔开,例如多个vip等级
	Msgtype       int         `json:"msgtype"`       // 消息类型1.通知 2.活动 3.公告 4.赛事公告
	Type          int         `json:"type"`          // 发送方式 1单条，2批量，3群组（所有的用户标题和内容一样）
	Title         string      `json:"title"`         // 消息标题
	Content       string      `json:"content"`       // 消息类容
	IsRead        int         `json:"isRead"`        // 是否已阅读
	Sticky        int         `json:"sticky"`        // 置顶 0未置顶，1已置顶
	SendTime      int64       `json:"sendTime"`      // 发送时间
	UpdateTime    int64       `json:"updateTime"`    // 修改时间
	ImgUrl        string      `json:"imgUrl"`        // 图片地址
	Icon          interface{} `json:"icon"`          // 消息图标
	IconType      string      `json:"iconType"`      // 图标类型 取值 1,2,3，4,5
	UserSystem    string      `json:"userSystem"`    // 用户体系1-会员，2-代理
	PushPlatform  string      `json:"pushPlatform"`  // 推送平台 0全站，1体育,2Web,3 H5,7棋牌，6彩票，8真人，9全站体育App
	PushDevice    string      `json:"pushDevice"`    // 推送设备1 android  0 ios
	Platform      int         `json:"platform"`      // 存储平台
	PushFlag      int         `json:"pushFlag"`      // 推送标记
	VipGrade      string      `json:"vipGrade"`      // vip等级
	VipGradeNum   int         `json:"vipGradeNum"`   // VIP等级数
	PcPath        string      `json:"pcPath"`        // pc图片
	H5Path        string      `json:"h5Path"`        // h5图片
	PcUrl         string      `json:"pcUrl"`         // pcUrl
	H5Url         string      `json:"h5Url"`         // pcUrl
	JumpUrlType   int         `json:"jumpUrlType"`   // 站内跳转类型
	DelFlag       int         `json:"delFlag"`       // 删除标志1已删除，0未删除
	SysStatus     int         `json:"sysStatus"`     // 消息状态1启用，0禁用
	Top           int         `json:"top"`           // 置顶状态1置顶，0未置顶
	Sort          int         `json:"sort"`          // 排序
	PushWay       int         `json:"pushWay"`       // 推送路径选择 默认为0-全部 1-只用极光推送 2-只用WS推送
	ImgTop        int         `json:"imgTop"`        // 图片是否置顶 0=否，1=是
}

// MsgTemplateDto 站内信发送消息实体
type MsgTemplateDto struct {
	SiteId          string   `json:"siteId"`          // 站点id
	MemberId        string   `json:"memberId"`        // 会员id
	TemplateNo      string   `json:"templateNo"`      // 模板编号
	ParamsValues    []string `json:"paramsValues"`    // 模板参数
	SendSiteMessage int      `json:"sendSiteMessage"` // 是否站内信通知 0是，1否
	TemplateType    int      `json:"templateType"`    // 选择模板 0系统模板，1自定义模板
	Title           string   `json:"title"`           // 消息标题
	Msgtype         int      `json:"msgtype"`         // 消息类型1.通知 2.活动 3.公告 4.赛事公告
	Content         string   `json:"content"`         // 消息类容
	Icon            string   `json:"icon"`            // 消息图标
	IconType        string   `json:"iconType"`        // 图标类型 取值 1,2,3，4,5
	BatchId         string   `json:"batchId"`         // 批次号
	PcPath          string   `json:"pcPath"`          // pc图片
	H5Path          string   `json:"h5Path"`          // h5图片
	PcUrl           string   `json:"pcUrl"`           // pcUrl
	H5Url           string   `json:"h5Url"`           // pcUrl
	JumpUrlType     int      `json:"jumpUrlType"`     // 站内跳转类型
}

// SiteInnerTemplateVo 站内信
type SiteInnerTemplateVo struct {
	Id             int64  `json:"id"`             // 文档id
	No             string `json:"no"`             // 模板id
	Title          string `json:"title"`          // 模板标题
	Content        string `json:"content"`        // 模板展示内容
	SysStatus      int    `json:"sysStatus"`      // 状态
	SendScenesCode string `json:"sendScenesCode"` // 发送场景编码
	SendScenesName string `json:"sendScenesName"` // 发送场景中文描述
	SendObjectCode string `json:"sendObjectCode"` // 发送对象编码
	SendObjectName string `json:"sendObjectName"` // 发送对象名称
	Category       int    `json:"category"`       // 类别 1系统默认,2自定义
	IsUse          int    `json:"isUse"`          // 是否使用中0未使用，1使用中
	UpdateTime     string `json:"updateTime"`     // 操作时间
	FunctionCode   string `json:"functionCode"`   // 发送对象编码
	FunctionName   string `json:"functionName"`   // 发送对象名称
	Params         string `json:"params"`         // 开发参数
	OriginalParams string `json:"originalParams"` // 原始开发参数
	Icon           string `json:"icon"`           // 图标key
	PcPath         string `json:"pcPath"`         // PcPath
	H5Path         string `json:"h5Path"`         // H5Path
	PcUrl          string `json:"pcUrl"`          // PcUrl
	H5Url          string `json:"h5Url"`          // H5Url
	JumpUrlType    int    `json:"jumpUrlType"`    // 跳转类型
}

// SendPushReq 极光推送发送请求对象实体类
type SendPushReq struct {
	NoticeStr      string             `json:"noticeStr"`      //提示字符串 仅为了打印日志 如：公告-UUID，通知-UUID
	BWPushSync     bool               `json:"BWPushSync"`     //是否异步发送
	MsgType        string             `json:"msgType"`        //推送消息类型(PushMsgTypeEnum 枚举定义,默认：通知） notification-通知 message-消息
	SiteId         string             `json:"siteId"`         //站点ID
	OtherSiteId    string             `json:"otherSiteId"`    //跨站点发送消息-用于查询极光配置的站点ID
	MemberIds      []int64            `json:"memberIds"`      //会员ID
	Alias          []string           `json:"alias"`          //设备别名
	Audience       string             `json:"audience"`       //广播所有人
	Extras         string             `json:"extras"`         //扩展字段 (推送 目标端传送数据，默认为空，填充{},以KEY：Value JSON形式传输)
	PlatformList   []*PushReqPlatform `json:"platformList"`   //平台列表
	PushTitle      string             `json:"pushTitle"`      //标题
	PushContent    string             `json:"pushContent"`    //内容
	RegistrationId []string           `json:"registrationId"` //注册ID
	Status         int                `json:"status"`         //发送类型 0.未发送（定时） 1，已发送（立刻） 默认1
	Tag            []string           `json:"tag"`            //并集
	TagAnd         []string           `json:"tagAnd"`         //交集
	TagNot         []string           `json:"tagNot"`         //补集
	StartDate      string             `json:"startDate"`      //定时时间
	Type           int                `json:"type"`           //类型（PushTypeEnum 枚举定义）默认是0，1 广播所有人 2 设备标签 3设备别名(Alias) 4 Registration ID
	TimeLive       int                `json:"timeLive"`       //离线时间
	BusinessTime   string             `json:"businessTime"`   //业务时间（yyyy-MM-dd HH:mm:ss,可为空）
}
type PushReq struct {
	SiteId       string             `json:"siteId"`
	MemberIds    []int64            `json:"memberIds"`
	Alias        []string           `json:"alias"`
	PushTitle    string             `json:"pushTitle"`
	PushContent  string             `json:"pushContent"`
	Status       int                `json:"status"`
	TimeLive     int                `json:"timeLive"`
	Type         int                `json:"type"`
	PlatformList []*PushReqPlatform `json:"platformList"`
	OtherSiteId  string             `json:"otherSiteId"`
}

type PushReqPlatform struct {
	Platform   string `json:"platform"`   // 平台编码
	DeviceType string `json:"deviceType"` //设备类型 多个设备支持逗号分隔，默认 all
}

// SendPushReqParam 参数结构体
type SendPushReqParam struct {
	Alias          []string `json:"alias"`          //设备别名
	Audience       string   `json:"audience"`       //广播所有人
	Extras         string   `json:"extras"`         //扩展字段 (推送 目标端传送数据，默认为空，填充{},以KEY：Value JSON形式传输)
	Platform       string   `json:"platform"`       //平台编码
	PushTitle      string   `json:"pushTitle"`      //标题
	PushContent    string   `json:"pushContent"`    //内容
	RegistrationId []string `json:"registrationId"` //注册ID
	Status         int      `json:"status"`         //发送类型 0.未发送（定时） 1，已发送（立刻） 默认1
	StartDate      string   `json:"startDate"`      //定时时间
	Type           int      `json:"type"`           //类型（PushTypeEnum 枚举定义）默认是0，1 广播所有人 2 设备标签 3设备别名(Alias) 4 Registration ID
	TimeLive       int      `json:"timeLive"`       //离线时间
	SysKey         string   `json:"sysKey"`         //系统的AppKey
	ApnsProduction bool     `json:"apnsProduction"`
	Tag            []string `json:"tag"`    //并集
	TagAnd         []string `json:"tagAnd"` //交集
	TagNot         []string `json:"tagNot"` //补集
}

type PushTypeEnum struct {
	Code     int
	CodeName string
}

// 推送设备类型枚举类
type DeviceTypeEnum struct {
	Code     string
	CodeName string
}

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type SitePushConfigQuery struct {
	SysCode      string   `json:"sysCode"`      // 系统名称
	SiteId       int      `json:"siteId"`       // 站点ID
	RelateSiteId int      `json:"relateSiteId"` // 推送配置关联的其他站点ID
	ClientTypes  []string `json:"clientTypes"`  // 客户端类型
	AppPlatform  string   `json:"appPlatform"`  // 平台 android ios
}

// 用户广播消费实体
type MsgDataToKafkaBrRead struct {
	MemberId     int    `json:"memberId"`   //会员Id
	SiteId       int    `json:"siteId"`     //站点Id
	Flag         int    `json:"flag"`       // 1:一键已读，2：批量已读
	Ids          []int  `json:"ids"`        //当批量已读时，传入的广播id
	CreateTime   string `json:"createTime"` //消息创建时间
	MsgType      int    `json:"msgType" form:"msgType" binding:"required"`
	UserSystem   string `json:"userSystem" form:"userSystem"`
	RegisterTime string `json:"registerTime"` //用户注册时间
	PlatformCode string `json:"platformCode"` //平台ID
}

// 广播消息表
type SiteLetterBroadcast struct {
	ID          int64  `gorm:"column:id" json:"id" form:"id"`
	Title       string `gorm:"column:title" json:"title" form:"title"`
	Content     string `gorm:"column:content" json:"content" form:"content"`
	Platform    int    `gorm:"column:platform" json:"platform" form:"platform"`
	Msgtype     int    `gorm:"column:msgtype" json:"msgType" form:"msgtype"`
	BatchId     string `gorm:"column:batch_id" json:"batchId" form:"batch_id"`
	SysStatus   int    `gorm:"column:sys_status" json:"sysStatus" form:"sys_status"`
	SiteId      int    `gorm:"column:site_id" json:"siteId" form:"site_id"`
	Top         int    `gorm:"column:top" json:"top" form:"top"`
	Icon        string `gorm:"column:icon" json:"icon" form:"icon"`
	UserSystem  string `gorm:"column:user_system" json:"useSystem" form:"user_system"`
	CreatedAt   string `gorm:"column:created_at" json:"createdAt" form:"created_at"`
	UpdatedAt   string `gorm:"column:updated_at" json:"updatedAt" form:"updated_at"`
	DelFlag     int    `gorm:"column:del_flag" json:"delFlag" form:"del_flag"`
	SendTime    string `gorm:"column:send_time" json:"sendTime" form:"send_time"`
	VipGradeNum int    `gorm:"column:vip_grade_num" json:"vipGradeNum" form:"vip_grade_num"`
}

// 消息已读
type SiteLetterRead struct {
	ID         int    `gorm:"column:id" json:"id" form:"id"`
	MsgId      int64  `gorm:"column:msg_id" json:"msg_id" form:"msg_id"`
	MemberId   int    `gorm:"column:member_id" json:"member_id" form:"member_id"`
	Category   int    `gorm:"column:category" json:"category" form:"category"`
	Msgtype    int    `gorm:"column:msgtype" json:"msgtype" form:"msgtype"`
	DelFlag    int    `gorm:"column:del_flag" json:"del_flag" form:"del_flag"`
	CreatedAt  string `gorm:"column:created_at" json:"created_at" form:"created_at"`
	UpdatedAt  string `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`
	UserSystem string `gorm:"column:user_system" json:"user_system" form:"user_system"`
	SiteId     int    `gorm:"column:site_id" json:"site_id" form:"site_id"`
}
