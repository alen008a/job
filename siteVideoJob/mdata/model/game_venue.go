package model

// 场馆表
type GameVenue struct {
	Id                int      `gorm:"column:id" json:"id"`                                    // 自增主键id
	EnName            string   `gorm:"column:en_name" json:"enName"`                           // 英文名称
	ZhName            string   `gorm:"column:zh_name" json:"zhName"`                           // 中文名称
	Category          int      `gorm:"column:category" json:"category"`                        // 场馆类型
	CommissionRate    float64  `gorm:"column:commission_rate" json:"commissionRate,omitempty"` // 场馆平台费率
	Sort              int      `gorm:"column:sort" json:"sort,omitempty"`                      // 排序
	Status            int      `gorm:"column:status" json:"status"`                            // 总控使用字段  场馆是否维护中  0 正常，1 维护中
	IsDisplay         int      `gorm:"column:is_display" json:"isDisplay"`                     // 站点使用字段  状态 0 展示,1 不展示
	Url               string   `gorm:"column:url" json:"url,omitempty"`                        // 图标
	ConfigItem        string   `gorm:"column:config_item" json:"configItem,omitempty"`         // 配置项
	WalletStatus      int      `gorm:"column:wallet_status" json:"walletStatus"`               // 钱包锁定状态1锁定0未锁定
	CreatedAt         string   `gorm:"column:created_at" json:"createdAt,omitempty"`           // 创建时间
	UpdatedAt         string   `gorm:"column:updated_at" json:"updatedAt,omitempty"`           // 最后更新时间
	GameType          string   `gorm:"column:game_type" json:"gameType,omitempty"`             // 游戏类型
	ChannelCode       string   `gorm:"column:channel_code" json:"channelCode,omitempty"`       // 游戏钱包code
	ChannelName       string   `gorm:"column:channel_name" json:"channelName,omitempty"`
	AdminWalletStatus int      `gorm:"column:admin_wallet_status" json:"adminWalletStatus"` // 总控操作钱包是否锁定(0-正常 1-锁定)
	RecycleId         int      `gorm:"column:recycle_id" json:"recycleId,omitempty"`        // 回收记录id
	IsRegister        bool     `json:"isRegister,omitempty"`                                // 返回字段  是否注册(数据查询忽略)
	ReachedStatus     string   `json:"reachedStatus,omitempty"`                             // 返回字段          (数据查询忽略)
	VenueTag          VenueTag `gorm:"-" json:"venueTag,omitempty"`                         // /   返回字段 配置标签（0:无;1:官方认证; 2:更多玩法; 3:新上线; 4:推荐:）
	StartAt           string   `gorm:"-" json:"startAt,omitempty"`                          // /   返回字段 维护开始时间
	EndAt             string   `gorm:"-" json:"endAt,omitempty"`                            // /   返回字段 维护结束时间
	IsDisplayMaintain int      `gorm:"-" json:"isDisplayMaintain"`                          // /   返回字段 是否维护
	ReasonRemark      string   `gorm:"-" json:"reasonRemark"`                               // /   返回字段 维护备注
	JumpVenueId       int      `gorm:"-" json:"jumpVenueId,omitempty"`                      // 场馆挂维护时指定跳转场馆ID
	Hint              string   `gorm:"-" json:"hint,omitempty"`                             // /   维护跳转场馆 提示语
	JumpChannelCode   string   `gorm:"-" json:"jumpChannelCode,omitempty"`                  // /   维护跳转场馆
	CategoryName      string   `gorm:"-" json:"categoryName,omitempty"`                     // / 场馆类型名
	RecycleResult     int      `gorm:"-" json:"recycleResult"`                              // /   场馆回收状态
	IsAutoJump        int      `gorm:"-" json:"isAutoJump"`                                 // 是否自动跳转 1:是 0:否
	IsUnknowEnd       int      `gorm:"-" json:"isUnknowEnd"`                                // 结束时间是否为待定 1:是 0:不是
	SortUpdatedAt     string   `gorm:"column:sort_updated_at" json:"sortUpdatedAt"`         // 排序字段更新时间
}

func (*GameVenue) TableName() string {
	return "game_venue"
}

type VenueTag struct {
	VenueTagOne string `json:"venueTagOne,omitempty"`
}
