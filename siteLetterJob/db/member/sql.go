package memberdb

import (
	"errors"
	"gorm.io/gorm"
	"siteLetterJob/db/sqldb"
)

type MemberInfoChange struct {
	Id                    int     `gorm:"id" json:"id"`                                          // 主键id
	MemberId              int     `gorm:"member_id" json:"memberId"`                             // 会员id
	TotalRecharge         float64 `gorm:"total_recharge" json:"totalRecharge"`                   // 总存款
	TotalRecord           float64 `gorm:"total_record" json:"totalRecord"`                       // 总流水
	VipGrade              int     `gorm:"vip_grade" json:"vipGrade"`                             // vip等级
	Svip                  int     `gorm:"svip" json:"svip"`                                      // svip等级
	SvipGradeChangeTime   string  `gorm:"svip_grade_change_time" json:"svipGradeChangeTime"`     // svip加入时间
	CreditLevel           int     `gorm:"credit_level" json:"creditLevel"`                       // 信用层级
	CreditGrade           int     `gorm:"credit_grade" json:"creditGrade"`                       // 信用等级
	VipGradeChangeTime    string  `gorm:"vip_grade_change_time" json:"vipGradeChangeTime"`       // vip等级改变时间
	CreditLevelChangeTime string  `gorm:"credit_level_change_time" json:"creditLevelChangeTime"` // 信用层级变更时间
	CreditGradeChangeTime string  `gorm:"credit_grade_change_time" json:"creditGradeChangeTime"` // 信用等级变更时间
	DayWithdrawalCounts   int     `gorm:"day_withdrawal_counts" json:"dayWithdrawalCounts"`      // 单日提款限制次数
	DayWithdrawalMoney    float64 `gorm:"day_withdrawal_money" json:"dayWithdrawalMoney"`        // 单日提款限额
	UpdatedAt             string  `gorm:"updated_at" json:"updatedAt"`                           // 数据修改时间
	CreatedAt             string  `gorm:"created_at" json:"createdAt"`                           // 数据入库时间
	LostSignPopTime       string  `gorm:"lost_sign_pop_time" json:"lostSignPopTime"`             // 掉签弹窗时间
	Stxx                  int     `gorm:"stxx" json:"stxx"`                                      // stxx
}

func (*MemberInfoChange) TableName() string {
	return "member_info_change"
}

type MemberInfo struct {
	Id                         int    `gorm:"id" json:"id"`                                                    // 自增主键id
	Name                       string `gorm:"name" json:"name"`                                                // 用户名（登录账号）
	Avatar                     string `gorm:"avatar" json:"avatar"`                                            // 头像地址
	Sex                        int    `gorm:"sex" json:"sex"`                                                  // 性别,0，未知，1，男；2，女
	Password                   string `gorm:"password" json:"password"`                                        // 加盐密码
	Salt                       string `gorm:"salt" json:"salt"`                                                // 盐值
	UpdatedAt                  string `gorm:"updated_at" json:"updatedAt"`                                     // 数据修改时间
	CreatedAt                  string `gorm:"created_at" json:"createdAt"`                                     // 数据入库时间
	Status                     int    `gorm:"status" json:"status"`                                            // 用户状态(0停用，1启用)
	TopId                      int    `gorm:"top_id" json:"topId"`                                             // 代理用户id
	Birthday                   string `gorm:"birthday" json:"birthday"`                                        // 出生日期
	SiteId                     int    `gorm:"site_id" json:"siteId"`                                           // 站点id
	NickName                   string `gorm:"nick_name" json:"nickName"`                                       // 昵称
	TagId                      string `gorm:"tag_id" json:"tagId"`                                             // 关联member_tag_info主键ID
	LastLoginIp                string `gorm:"last_login_ip" json:"lastLoginIp"`                                // 最后登录ip
	RegisterDeviceId           string `gorm:"register_device_id" json:"registerDeviceId"`                      // 注册设备号
	LastLoginDeviceId          string `gorm:"last_login_device_id" json:"lastLoginDeviceId"`                   // 最后登录设备号
	RegisterIp                 string `gorm:"register_ip" json:"registerIp"`                                   // 注册IP
	LastLoginTime              string `gorm:"last_login_time" json:"lastLoginTime"`                            // 上次登录时间
	SourceUrl                  string `gorm:"source_url" json:"sourceUrl"`                                     // 域名
	InviteCode                 string `gorm:"invite_code" json:"inviteCode"`                                   // icode
	CodeSource                 string `gorm:"code_source" json:"codeSource"`                                   // code资源
	RegisterDevice             string `gorm:"register_device" json:"registerDevice"`                           // 注册设备(端)
	LoginDevice                string `gorm:"login_device" json:"loginDevice"`                                 // 登陆设备(端)
	WithdrawPassword           string `gorm:"withdraw_password" json:"withdrawPassword"`                       // 取款密码
	AddressCipher              string `gorm:"address_cipher" json:"addressCipher"`                             // 用户详细地址密文
	MoveFlag                   int    `gorm:"move_flag" json:"moveFlag"`                                       // 0. 不是迁移用户 1.迁移中 2.迁移完成未登录 3.迁移完成已登录
	ProvincesCipher            string `gorm:"provinces_cipher" json:"provincesCipher"`                         // 省市区脱敏
	RegisterType               int    `gorm:"register_type" json:"registerType"`                               // 注册类型  0官方注册  1三方注册
	RegisterSourceName         string `gorm:"register_source_name" json:"registerSourceName"`                  // 注册来源
	RegisterSourceCode         string `gorm:"register_source_code" json:"registerSourceCode"`                  // 注册来源code
	AppId                      string `gorm:"app_id" json:"appId"`                                             // 三方注册appid
	XsS0                       string `gorm:"xs_s0" json:"xsS0"`                                               // xs_s0
	XsS1                       string `gorm:"xs_s1" json:"xsS1"`                                               // xs_s1
	XsS5                       string `gorm:"xs_s5" json:"xsS5"`                                               // xs_s5
	XsS6                       string `gorm:"xs_s6" json:"xsS6"`                                               // xs_s6
	XsS20                      string `gorm:"xs_s20" json:"xsS20"`                                             // xs_s20
	PhoneDesensitization       string `gorm:"phone_desensitization" json:"phoneDesensitization"`               // 脱敏手机号码
	RealNameDesensitization    string `gorm:"real_name_desensitization" json:"realNameDesensitization"`        // 脱敏真实姓名
	QqDesensitization          string `gorm:"qq_desensitization" json:"qqDesensitization"`                     // 脱敏qq号码
	WechatDesensitization      string `gorm:"wechat_desensitization" json:"wechatDesensitization"`             // 脱敏微信号
	EmailDesensitization       string `gorm:"email_desensitization" json:"emailDesensitization"`               // 脱敏邮箱地址
	RegisterIpCipher           string `gorm:"register_ip_cipher" json:"registerIpCipher"`                      // 注册ip密文
	RegisterIpDesensitization  string `gorm:"register_ip_desensitization" json:"registerIpDesensitization"`    // 注册ip脱敏
	LastLoginIpCipher          string `gorm:"last_login_ip_cipher" json:"lastLoginIpCipher"`                   // 最后登录ip密文
	LastLoginIpDesensitization string `gorm:"last_login_ip_desensitization" json:"lastLoginIpDesensitization"` // 最后登录ip脱敏
	PhoneHomeDesensitization   string `gorm:"phone_home_desensitization" json:"phoneHomeDesensitization"`      // 脱敏手机归属地
	XsS41                      string `gorm:"xs_s41" json:"xsS41"`                                             // xs_s41
	Stxx                       int    `gorm:"stxx" json:"stxx"`                                                // stxx
	VipGrad                    int    `gorm:"column:vip_grade" json:"vipGrade"`                                // vip等级
	Svip                       int    `gorm:"column:svip" json:"svip"`                                         // svip等级
}

func (*MemberInfo) TableName() string {
	return "member_info"
}

// GetMemberVipGrade 根据会员Id查询会员等级
func GetMemberVipGrade(memberId int) (int, error) {
	var memberData MemberInfoChange
	model := sqldb.SiteSlave().Table(memberData.TableName())
	err := model.Where("member_id = ? ", memberId).Take(&memberData).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return -1, err
	} else if err != nil {
		return 0, err
	}

	return memberData.VipGrade, nil
}

// GetMemberInfo 根据会员Id查询会员等级
func GetMemberInfo(memberId int64) (*MemberInfo, error) {
	var member *MemberInfo
	model := sqldb.SiteSlave().Table("member_info")
	err := model.Where("id = ? ", memberId).Take(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return member, nil
}

// SiteInnerMessageTemplate 站内信模板
type SiteInnerMessageTemplate struct {
	Id             int    `gorm:"id" json:"id"`                           // 主键id
	No             string `gorm:"no" json:"no"`                           // 模板编号
	Category       int    `gorm:"category" json:"category"`               // 1系统模板，2自定义模板
	Title          string `gorm:"title" json:"title"`                     // 标题
	Content        string `gorm:"content" json:"content"`                 // 内容
	SysStatus      int    `gorm:"sys_status" json:"sysStatus"`            // 0停用，1启用
	SendScenesCode string `gorm:"send_scenes_code" json:"sendScenesCode"` // 发送场景编码
	SendScenesName string `gorm:"send_scenes_name" json:"sendScenesName"` // 发送场景名称
	SendObjectCode string `gorm:"send_object_code" json:"sendObjectCode"` // 发送对象编码
	SendObjectName string `gorm:"send_object_name" json:"sendObjectName"` // 发送对象名称
	FunctionCode   string `gorm:"function_code" json:"functionCode"`      // 功能编码
	FunctionName   string `gorm:"function_name" json:"functionName"`      // 功能名称
	Params         string `gorm:"params" json:"params"`                   // 参数
	CreatedAt      string `gorm:"created_at" json:"createdAt"`            // 数据入库时间
	UpdatedAt      string `gorm:"updated_at" json:"updatedAt"`            // 数据更新时间
	CreatedBy      string `gorm:"created_by" json:"createdBy"`            // 数据创建账号
	UpdatedBy      string `gorm:"updated_by" json:"updatedBy"`            // 数据修改账号
	IsUse          int    `gorm:"is_use" json:"isUse"`                    // 0未使用，1已使用
	DelFlag        int    `gorm:"del_flag" json:"delFlag"`                // 0未删除，1已删除
	BelongNo       string `gorm:"belong_no" json:"belongNo"`              // 所属系统模板编码
	Icon           string `gorm:"icon" json:"icon"`                       // 图标id
	OriginalParams string `gorm:"original_params" json:"originalParams"`  // 原始开发参数
}

func (*SiteInnerMessageTemplate) TableName() string {
	return "site_inner_message_template"
}

// GetSiteInnerMessageTemplate 根据会员Id查询会员等级
func GetSiteInnerMessageTemplate(no string) (*SiteInnerMessageTemplate, error) {
	var innerMsgTmpl *SiteInnerMessageTemplate
	model := sqldb.SiteSlave().Table("site_inner_message_template")
	err := model.Where("del_flag = 0 and is_use = 1 and belong_no = ? ", no).First(&innerMsgTmpl).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return innerMsgTmpl, nil
}
