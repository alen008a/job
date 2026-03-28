package utils

type StatusCode = int

const (
	StatusOK               StatusCode = 6000 // 正常返回
	ErrAccess                         = 6001 // 登录状态失效 类似401, 会跳到登录窗口
	ErrAccount                        = 6002 // 账号或密码错误
	ErrRefuse                         = 6003 // 访问被拒绝类似403
	ErrNotFound                       = 6004 // 找不到接口类似404
	ErrInternal                       = 6005 // 接口服务器产生错误 类似500
	ErrInvalidGateway                 = 6006 // 无效网关 类似502
	ErrAPIFailed                      = 6007 // 接口请求失败, 无法接受参数，会跳到登录窗口
	ErrInvalidParams                  = 6008 // 违法的参数
	ErrDataExistFailed                = 6009 // 数据已经存在
	ErrFileUploadFailed               = 6010 // 文件上传失败
	ErrLockedFailed                   = 6011 // 频繁提交拒绝
	ErrWarning                        = 6012 // 警告 非常用登录IP(web h5)或者非常用设备（移动端）
	ErrMaintenance                    = 6013 // 维护
	ErrCallError                      = 6014 // 调用错误
	ErrNotBindPhone                   = 6015 // 未绑定手机号的错误码
	ErrLoginVerify                    = 6016 // 用户登录需要手机验证
	ErrLoginVerify2                   = 6017 // 用户登录需要手机验证
	ErrGraphicVerification            = 6021 // 图形验证码
	ErrGeetestVerification            = 6022 // 极速验证
	ErrWalletAccess                   = 6099 // 钱包token失效
	Err60sLimit                       = 6066 // 60秒限制特殊状态码
	ErrCheckPassword                  = 6100 // 校验密码错误

	// 异地登陆状态码
	ErrIPNotAllowLogin        = 6025 // IP不通过
	ErrDeviceNotAllowLogin    = 6026 // 设备不通过
	ErrNotAllowLogin          = 6027 // 设备IP 均不通过校验
	ErrLongTimeNotLogin       = 6028 // 长时间未登陆
	ErrRiskLoginNotAllow      = 6029 // 登陆存在风险 不允许直接登陆 需要手机号验证
	ErrContactCustomerService = 6030 // 联系客服

	ErrLoginForbidAreaLimit = 6031 // 区域限制禁止登陆
	ErrAccessWarningLimit   = 6032 // 区域警告
	ErrIpLimitAccess        = 6033 // IP限制访问
	ErrLoginPhoneLimit      = 6034 // 手机号码输入错误次数上限 返回

	// 二级密码业务状态码
	ErrTwoPasswordUpdate = 6100 // 二级密码绑定失败
	// 提款保护业务状态码
	ErrNotFondGradeCantLock = 6110 //未获取当前会员的对应等级,无法加锁
	ErrNotCfgCantLock       = 6111 //提款保护配置更新异常

	ErrVenueMaintain = 6415 //场馆维护中
)

// 只能写通用的msg，非通用的，直接在c.WebRsp里面返回

type Message = string

const (
	MsgSuccess             Message = "成功"
	MsgInternalError               = "服务器错误"
	MsgAccessError                 = "请重新登录"
	MsgFileUploadError             = "文件上传失败"
	MsgInvalidParamsError          = "请求参数错误"
	MsgRefuseError                 = "请求拒绝"
	MsgLockedError                 = "请勿频繁提交"
	MsgMaintenanceError            = "维护中"
	MsgNotFoundError               = "未找到数据"
	MsgIllegalRequestError         = "非法请求" // 请求头未传加密串X-API-TIMESTAMP的时候使用
)
