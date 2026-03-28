package member

var (
	MemberConf Config
)

type Config struct {
	MsgAppCode    string `json:"msgAppCode"`    //告警群code
	MsgAppKey     string `json:"msgAppKey"`     //告警群key
	KfkTopic      string `json:"kfkTopic"`      //kafka推送站内信topic
	KfkTemplateNo string `json:"kfkTemplateNo"` //kafka推送模板编号 固定："20200604165111"
}

// 会员vip等级配置表
type MemberGradeConfig struct {
	Id             int     `gorm:"column:id" json:"id"`                           //主键
	Name           string  `gorm:"column:name" json:"name"`                       //会员等级
	UpgradeDeposit float64 `gorm:"column:upgrade_deposit" json:"upgrade_deposit"` //升级存款
	UpgradeFlows   float64 `gorm:"column:upgrade_flows" json:"upgrade_flows"`     //升级总流水
	RelegateFlows  float64 `gorm:"column:relegate_flows" json:"relegate_flows"`   //保级流水
	//UpgradeBonus        float64 `gorm:"column:upgrade_bonus" json:"upgrade_bonus"`                   //升级红利
	//BirthCash           float64 `gorm:"column:birth_cash" json:"birth_cash"`                         //生日礼金
	//MonthGift           float64 `gorm:"column:month_gift" json:"month_gift"`                         //上半月红包
	//DayWithdrawalCounts int64   `gorm:"column:day_withdrawal_counts" json:"day_withdrawal_counts"`   //单日提款限制次数
	//DayWithdrawalMoney  float64 `gorm:"column:day_withdrawal_money" json:"day_withdrawal_money"`     //单日提款限额
	//CreatedAt           int64   `gorm:"column:created_at" json:"created_at"`                         //创建时间
	//CreateBy            string  `gorm:"column:create_by" json:"create_by"`                           //创建人
	//UpdatedAt           int64   `gorm:"column:updated_at" json:"updated_at"`                         //数据修改时间
	//UpdateBy            string  `gorm:"column:update_by" json:"update_by"`                           //修改人
	//SecondHalfMonthGift float64 `gorm:"column:second_half_month_gift" json:"second_half_month_gift"` //下半月红包
	//WeekGift            float64 `gorm:"column:week_gift" json:"week_gift"`                           //周红包
	//SiteId                int64   `gorm:"column:stxx" json:"stxx"`                                     //stxx
}

func (b *MemberGradeConfig) TableName() string {
	return "member_grade_config"
}

// 用户基础信息中经常变更数据
type MemberInfoChange struct {
	Id                    int     `gorm:"column:id" json:"id"`                                          // 主键id
	MemberId              int     `gorm:"column:member_id" json:"memberId"`                             // 会员id
	TotalRecharge         float64 `gorm:"column:total_recharge" json:"totalRecharge"`                   // 总存款
	TotalRecord           float64 `gorm:"column:total_record" json:"totalRecord"`                       // 总流水
	VipGrade              int     `gorm:"column:vip_grade" json:"vipGrade"`                             // vip等级
	Svip                  int     `gorm:"column:svip" json:"svip"`                                      // svip等级
	CreditLevel           int     `gorm:"column:credit_level" json:"creditLevel"`                       // 信用层级
	CreditGrade           int     `gorm:"column:credit_grade" json:"creditGrade"`                       // 信用等级
	VipGradeChangeTime    string  `gorm:"column:vip_grade_change_time" json:"vipGradeChangeTime"`       // vip等级改变时间
	CreditLevelChangeTime string  `gorm:"column:credit_level_change_time" json:"creditLevelChangeTime"` // 信用层级变更时间
	CreditGradeChangeTime string  `gorm:"column:credit_grade_change_time" json:"creditGradeChangeTime"` // 信用等级变更时间
	DayWithdrawalCounts   int     `gorm:"column:day_withdrawal_counts" json:"dayWithdrawalCounts"`      // 单日提款限制次数
	DayWithdrawalMoney    float64 `gorm:"column:day_withdrawal_money" json:"dayWithdrawalMoney"`        // 单日提款限额
	UpdatedAt             string  `gorm:"column:updated_at" json:"updatedAt"`                           // 数据修改时间
	CreatedAt             string  `gorm:"column:created_at" json:"createdAt"`                           // 数据入库时间
	LostSignPopTime       string  `gorm:"column:lost_sign_pop_time" json:"lostSignPopTime"`             // 掉签弹窗时间
}

func (b *MemberInfoChange) TableName() string {
	return "member_info_change"
}

type MemberDayInfo struct {
	OmitChangeFeeRecord float64 `gorm:"column:omit_change_fee_record" json:"omit_change_fee_record"` // 主键id
}

func (b *MemberDayInfo) TableName() string {
	return "member_day_info"
}

// 会员等级记录表
type MemberGradeInfo struct {
	ID              int     `gorm:"column:id" json:"id"`                            // 主键id
	MemberAccount   string  `gorm:"column:member_account" json:"memberAccount"`     // 会员账号
	MemberId        int     `gorm:"column:member_id" json:"memberId"`               // 会员id
	TotalDeposit    float64 `gorm:"column:total_deposit" json:"totalDeposit"`       // 累计存款
	TotalJournal    float64 `gorm:"column:total_journal" json:"totalJournal"`       // 累计流水
	BeforeGrade     string  `gorm:"column:before_grade" json:"beforeGrade"`         // 调整前等级
	AfterGrade      string  `gorm:"column:after_grade" json:"afterGrade"`           // 调整后等级
	ChangeType      int     `gorm:"column:change_type" json:"changeType"`           // 调整类型0保级1升级2降级3初始化
	UpdatedAt       string  `gorm:"column:updated_at" json:"updatedAt"`             // 更新时间
	CreatedAt       string  `gorm:"column:created_at" json:"createdAt"`             // 创建人
	UpdatedBy       string  `gorm:"column:updated_by" json:"updatedBy"`             // 更新人
	CreatedBy       string  `gorm:"column:created_by" json:"createdBy"`             // 创建人
	GradeStartTime  string  `gorm:"column:grade_start_time" json:"gradeStartTime"`  // 升级时间段-开始时间
	GradeEndTime    string  `gorm:"column:grade_end_time" json:"gradeEndTime"`      // 升级时间段-结束时间
	UpdatedType     int     `gorm:"column:updated_type" json:"updatedType"`         // 操作类型0：手动，1：自动
	DemotionDeposit float64 `gorm:"column:demotion_deposit" json:"demotionDeposit"` //降级后累计存款',
	DemotionJournal float64 `gorm:"column:demotion_journal" json:"demotionJournal"` //降级后累计流水',
	IsKeep          int     `gorm:"column:is_keep" json:"isKeep"`                   //是否保级（0 否 1 是）',
	KeepGradeDate   string  `gorm:"column:keep_grade_date" json:"keepGradeDate"`    //审核保级日期,
	PeriodRecord    float64 `gorm:"column:period_record" json:"periodRecord"`       //保级或降级季流水
}

func (b *MemberGradeInfo) TableName() string {
	return "member_grade_info"
}

// kafka推送模板
type KfkMsgTemplateDto struct {
	MemberId     int      `json:"memberId"`
	SiteId       string   `json:"siteId"`
	TemplateNo   string   `json:"templateNo"`
	ParamsValues []string `json:"paramsValues"`
}
