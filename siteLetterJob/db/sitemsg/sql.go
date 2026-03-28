package sqldb

import (
	"fmt"
	"gorm.io/gorm"
	"siteLetterJob/db/sqldb"
	"siteLetterJob/internal/glog"
	"siteLetterJob/mdata/sitemsg"
	"strconv"
)

// SiteLetterBroadcast 广播消息表
type SiteLetterBroadcast struct {
	Id          int    `gorm:"column:id" json:"id"`                     // 主键id
	Title       string `gorm:"column:title" json:"title"`               // 标题
	Content     string `gorm:"column:content" json:"content"`           // 内容
	Platform    int    `gorm:"column:platform" json:"platform"`         // 平台
	Msgtype     int    `gorm:"column:msgtype" json:"msgtype"`           // 消息类型1通知，2活动
	BatchId     string `gorm:"column:batch_id" json:"batchId"`          // 批次号
	SysStatus   int    `gorm:"column:sys_status" json:"sysStatus"`      // 状态1启用，0停用
	SiteId      int    `gorm:"column:site_id" json:"siteId"`            // 站点id
	Top         int    `gorm:"column:top" json:"top"`                   // 置顶1置顶，0未置顶
	Icon        string `gorm:"column:icon" json:"icon"`                 // 图标id
	UserSystem  string `gorm:"column:user_system" json:"userSystem"`    // 用户类型1会员，2代理
	CreatedAt   string `gorm:"-" json:"createdAt"`                      // 数据创建时间
	UpdatedAt   string `gorm:"-" json:"updatedAt"`                      // 数据更新时间
	DelFlag     int    `gorm:"column:del_flag" json:"delFlag"`          // 删除标志1已删除，0未删除
	SendTime    string `gorm:"column:send_time" json:"sendTime"`        // 发送时间，兼容原es存储时间
	VipGradeNum int    `gorm:"column:vip_grade_num" json:"vipGradeNum"` // vip等级和默认8190代表所有等级
	PcPath      string `gorm:"column:pc_path" json:"pcPath"`            // pc图片路径
	H5Path      string `gorm:"column:H5_path" json:"h5Path"`            // h5图片路径
	PcUrl       string `gorm:"column:pc_url" json:"pcUrl"`              // pcUrl
	H5Url       string `gorm:"column:H5_url" json:"h5Url"`              // h5Url
	Sort        int    `gorm:"column:sort" json:"sort"`                 // url
	JumpUrlType int    `gorm:"column:jump_url_type" json:"jumpUrlType"` //跳转链接类型，1站内2站外
	ImgTop      int    `gorm:"column:img_top" json:"imgTop"`            //图片是否置顶 0=否，1=是
}

func (*SiteLetterBroadcast) TableName() string {
	return "site_letter_broadcast"
}

// SiteLetterOrientation 定向消息发送表
type SiteLetterOrientation struct {
	Id          int    `gorm:"column:id" json:"id"`                     // 主键id
	MemberId    int    `gorm:"column:member_id" json:"memberId"`        // 会员id
	SiteId      int    `gorm:"column:site_id" json:"siteId"`            // 站点id
	Title       string `gorm:"column:title" json:"title"`               // 标题
	Content     string `gorm:"column:content" json:"content"`           // 内容
	Platform    int    `gorm:"column:platform" json:"platform"`         // 平台编号
	Msgtype     int    `gorm:"column:msgtype" json:"msgtype"`           // 消息类型1通知，2活动
	BatchId     string `gorm:"column:batch_id" json:"batchId"`          // 消息批次号
	SysStatus   int    `gorm:"column:sys_status" json:"sysStatus"`      // 消息状态1启用，0禁用
	IsRead      int    `gorm:"column:is_read" json:"isRead"`            // 消息读取状态1已读，0未读
	UserSystem  string `gorm:"column:user_system" json:"userSystem"`    // 用户体系1会员，2代理
	Top         int    `gorm:"column:top" json:"top"`                   // 置顶状态1置顶，0未置顶
	Icon        string `gorm:"column:icon" json:"icon"`                 // 图标id
	CreatedAt   string `gorm:"-" json:"createdAt"`                      // 数据入库时间
	UpdatedAt   string `gorm:"-" json:"updatedAt"`                      // 数据更新时间
	DelFlag     int    `gorm:"column:del_flag" json:"delFlag"`          // 删除标志1已删除,0未删除
	SendTime    string `gorm:"column:send_time" json:"sendTime"`        // 发送时间，兼容原es存储时间
	PcPath      string `gorm:"column:pc_path" json:"pcPath"`            // pc图片路径
	H5Path      string `gorm:"column:H5_path" json:"h5Path"`            // h5图片路径
	PcUrl       string `gorm:"column:pc_url" json:"pcUrl"`              // pcUrl
	H5Url       string `gorm:"column:H5_url" json:"h5Url"`              // h5Url
	Sort        int    `gorm:"column:sort" json:"sort"`                 // url
	JumpUrlType int    `gorm:"column:jump_url_type" json:"jumpUrlType"` //跳转链接类型，1站内2站外
	ImgTop      int    `gorm:"column:img_top" json:"imgTop"`            //图片是否置顶 0=否，1=是
}

type MsgCountVo struct {
	NoticeCount   int `json:"noticeCount"`   //通知
	ActivityCount int `json:"activityCount"` //活动
	BulletinCount int `json:"bulletinCount"` //公告
	MatchCount    int `json:"matchCount"`    //赛事通告
}

// 站内信未读消息分类统计响应
type LetterUnreadForTypeRespV2 struct {
	BRNoticeCount   int `json:"brnotice_count"`
	BRActivityCount int `json:"bractivity_count"`
	BRFdcount       int `json:"brfd_count"`
	BRBulletinCount int `json:"brbulletin_count"`
	BRMatchCount    int `json:"brmatch_count"`
	NoticeCount     int `json:"notice_count"`   //通知
	ActivityCount   int `json:"activity_count"` //活动
	Fdcount         int `json:"fd_count"`       //财务
	BulletinCount   int `json:"bulletin_count"` //公告
	MatchCount      int `json:"match_count"`    //赛事通告
}

func (*SiteLetterOrientation) TableName(memberId int64) string {
	return fmt.Sprintf("site_letter_orientation_%d", memberId&63)
}

// SitePushLog 推送日志
type SitePushLog struct {
	Id           int    `gorm:"column:id" json:"id"`                                                  // 主键
	BatchNo      string `gorm:"column:batch_no" json:"batchNo"`                                       // 批次号 同一批次的批次号相同
	ConfigId     int    `gorm:"column:config_id;default:0" json:"configId"`                           // 配置ID
	SysCode      string `gorm:"column:sys_code" json:"sysCode"`                                       // 系统名称
	SiteId       int    `gorm:"column:site_id;default:0" json:"siteId"`                               // 站点ID
	ClientType   string `gorm:"column:client_type" json:"clientType"`                                 // 客户端类型
	AppPlatform  string `gorm:"column:app_platform" json:"appPlatform"`                               // 平台 android ios
	MsgType      string `gorm:"column:msg_type" json:"msgType"`                                       // 推送消息类型 notification-通知 message-消息
	PushTitle    string `gorm:"column:push_title" json:"pushTitle"`                                   // 推送标题
	PushContent  string `gorm:"column:push_content" json:"pushContent"`                               // 推送内容
	PushStatus   int    `gorm:"column:push_status;default:0" json:"pushStatus"`                       // 状态 1-成功，2-失败
	PushParam    string `gorm:"column:push_param;default:null" json:"pushParam"`                      // 推送参数
	ResponseText string `gorm:"column:response_text" json:"responseText"`                             // 返回文本值
	PushTime     string `gorm:"column:push_time;default:null" json:"pushTime"`                        // 推送时间
	ResponseTime string `gorm:"column:response_time;default:1971-01-01 00:00:00" json:"responseTime"` // 返回时间
	CreatedBy    string `gorm:"column:created_by" json:"createdBy"`                                   // 创建人
	CreatedAt    string `gorm:"column:created_at;default:null" json:"createdAt"`                      // 创建时间
	UpdatedBy    string `gorm:"column:updated_by" json:"updatedBy"`                                   // 更新者
	UpdatedAt    string `gorm:"column:updated_at;default:null" json:"updatedAt"`                      // 更新时间
}

func (*SitePushLog) TableName() string {
	return "site_push_log"
}

// InsertSitePushConfig 插入返回结果
func InsertSitePushConfig(log *SitePushLog) error {
	model := sqldb.Site().Model(log).Create(&log)
	err := model.Error
	if err != nil {
		glog.Errorf("InsertSitePushConfig Create err:%v | param: dto=%+v", err, log)
		return err
	}
	if model.RowsAffected == 0 {
		return fmt.Errorf("insert rowsAffected==0")
	}
	return nil
}

// InsertSiteLetterOrientation 定向消息
func InsertSiteLetterOrientation(tableName string, msgDto *sitemsg.MsgDto) (int64, int64, error) {
	data := newSiteLetterOrientation(msgDto)
	cmd := sqldb.EdgeDB().Table(tableName).Create(&data)
	return cmd.RowsAffected, int64(data.Id), cmd.Error
}

func newSiteLetterBroadcast(msg *sitemsg.MsgDto) *SiteLetterBroadcast {
	data := &SiteLetterBroadcast{
		Title:       msg.Title,
		Content:     msg.Content,
		Platform:    msg.Platform,
		Msgtype:     msg.Msgtype,
		BatchId:     msg.BatchId,
		SysStatus:   msg.SysStatus,
		SiteId:      msg.SiteId,
		Top:         msg.Top,
		Icon:        msg.IconType,
		UserSystem:  msg.UserSystem,
		DelFlag:     msg.DelFlag,
		SendTime:    fmt.Sprintf("%d", msg.SendTime),
		VipGradeNum: msg.VipGradeNum,
		PcPath:      msg.PcPath,
		H5Path:      msg.H5Path,
		PcUrl:       msg.PcUrl,
		H5Url:       msg.H5Url,
		JumpUrlType: msg.JumpUrlType,
		ImgTop:      msg.ImgTop,
	}
	if msg.Sort > 0 {
		data.Sort = msg.Sort
	}
	return data

}
func newSiteLetterOrientation(msg *sitemsg.MsgDto) *SiteLetterOrientation {
	id, _ := strconv.Atoi(msg.Id)
	data := &SiteLetterOrientation{
		Id:          id,
		MemberId:    msg.MemberId,
		SiteId:      msg.SiteId,
		Title:       msg.Title,
		Content:     msg.Content,
		Platform:    msg.Platform,
		Msgtype:     msg.Msgtype,
		BatchId:     msg.BatchId,
		SysStatus:   msg.SysStatus,
		IsRead:      msg.IsRead,
		UserSystem:  msg.UserSystem,
		Top:         msg.Top,
		Icon:        msg.IconType,
		UpdatedAt:   fmt.Sprintf("%d", msg.UpdateTime),
		DelFlag:     msg.DelFlag,
		SendTime:    fmt.Sprintf("%d", msg.SendTime),
		PcPath:      msg.PcPath,
		H5Path:      msg.H5Path,
		PcUrl:       msg.PcUrl,
		H5Url:       msg.H5Url,
		JumpUrlType: msg.JumpUrlType,
		ImgTop:      msg.ImgTop,
	}
	if msg.Sort > 0 {
		data.Sort = msg.Sort
	}
	return data
}

// InsertSiteLetterBroadcast 广播消息
func InsertSiteLetterBroadcast(tableName string, msgDto *sitemsg.MsgDto) (int64, error) {
	data := newSiteLetterBroadcast(msgDto)
	cmd := sqldb.EdgeDB().Table(tableName).Create(&data)
	return cmd.RowsAffected, cmd.Error
}

// UpdateSiteLetterBroadcast 修改广播消息表
func UpdateSiteLetterBroadcast(msgDto *sitemsg.MsgDto) (int64, error) {
	if msgDto == nil {
		return 0, nil
	}
	glog.Infof("修改广播消息表1：%+v", msgDto)
	data := newSiteLetterBroadcast(msgDto)
	glog.Infof("修改广播消息表2：%+v", data)
	switch msgDto.OperationType {
	case "2":
		updates := sqldb.EdgeDB().Table("site_letter_broadcast").Where("batch_id = ? and site_id = ?", msgDto.BatchId, msgDto.SiteId).Updates(&data)
		return updates.RowsAffected, updates.Error
	case "3":
		attrs := map[string]interface{}{}
		attrs["del_flag"] = 1
		updates := sqldb.EdgeDB().Table("site_letter_broadcast").Where("batch_id = ? and site_id = ?", msgDto.BatchId, msgDto.SiteId).Updates(attrs)
		return updates.RowsAffected, updates.Error
	case "4":
		DeleteLetterBroadcastOtherMsg(msgDto)
		attrs := map[string]interface{}{}
		attrs["sys_status"] = 1
		updates := sqldb.EdgeDB().Table("site_letter_broadcast").Where("batch_id = ? and site_id = ?", msgDto.BatchId, msgDto.SiteId).Updates(attrs)
		return updates.RowsAffected, updates.Error
	case "5":
		attrs := map[string]interface{}{}
		attrs["sys_status"] = 0
		updates := sqldb.EdgeDB().Table("site_letter_broadcast").Where("batch_id = ? and site_id = ?", msgDto.BatchId, msgDto.SiteId).Updates(attrs)
		return updates.RowsAffected, updates.Error
	case "6":
		attrs := map[string]interface{}{}
		attrs["top"] = 1
		updates := sqldb.EdgeDB().Table("site_letter_broadcast").Where("batch_id = ? and site_id = ?", msgDto.BatchId, msgDto.SiteId).Updates(attrs)
		return updates.RowsAffected, updates.Error
	case "7":
		attrs := map[string]interface{}{}
		attrs["top"] = 0
		updates := sqldb.EdgeDB().Table("site_letter_broadcast").Where("batch_id = ? and site_id = ?", msgDto.BatchId, msgDto.SiteId).Updates(attrs)
		return updates.RowsAffected, updates.Error
	}
	return 0, nil
}

func DeleteLetterBroadcastOtherMsg(msgDto *sitemsg.MsgDto) {
	var siteLetterBroadcast *SiteLetterBroadcast
	//可能会存在多条相同batch_id的消息，只保留更新时间最新的消息，其余同批次消息删掉
	err := sqldb.EdgeDB().Table("site_letter_broadcast").Where("batch_id = ? and site_id = ?", msgDto.BatchId, msgDto.SiteId).Order("created_at DESC").Take(&siteLetterBroadcast).Error
	if err == nil && siteLetterBroadcast != nil && siteLetterBroadcast.Id > 0 {
		//删除同批次其他ID的消息
		result := sqldb.EdgeDB().Table("site_letter_broadcast").Where("batch_id = ? and site_id = ? and id <> ?", msgDto.BatchId, msgDto.SiteId, siteLetterBroadcast.Id).
			Delete(&SiteLetterBroadcast{})
		if result.Error != nil {
			glog.Errorf("UpdateSiteLetterBroadcast OperationType=4 删除同批次 ID=%v 以外的消息 失败! params:batchId=%v siteId=%v error:%+v", siteLetterBroadcast.Id, msgDto.BatchId, msgDto.SiteId, result.Error)
		} else {
			glog.Infof("UpdateSiteLetterBroadcast OperationType=4 删除同批次 ID=%v 以外的消息 成功! params:batchId=%v siteId=%v 删除条数:%v", siteLetterBroadcast.Id, msgDto.BatchId, msgDto.SiteId, result.RowsAffected)
		}
	}
}

// UpdateSiteLetterOrientation 向消息修改
func UpdateSiteLetterOrientation(tableName string, memberId int64, msgDto *sitemsg.MsgDto) (int64, error) {
	data := newSiteLetterOrientation(msgDto)
	switch msgDto.OperationType {
	case "2":
		updates := sqldb.EdgeDB().Table(tableName).Where("id = ?", msgDto.Id).Updates(&data)
		return updates.RowsAffected, updates.Error
	case "3":
		attrs := map[string]interface{}{}
		attrs["del_flag"] = 1
		updates := sqldb.EdgeDB().Table(tableName).Where("id = ?", msgDto.Id).Updates(attrs)
		return updates.RowsAffected, updates.Error
	case "4":
		id, _ := strconv.Atoi(msgDto.Id)
		DeleteSiteLetterOrientationOtherMsg(tableName, memberId, int64(id), msgDto)
		attrs := map[string]interface{}{}
		attrs["sys_status"] = 1
		updates := sqldb.EdgeDB().Table(tableName).Where("id = ?", msgDto.Id).Updates(attrs)
		return updates.RowsAffected, updates.Error
	case "5":
		attrs := map[string]interface{}{}
		attrs["sys_status"] = 0
		updates := sqldb.EdgeDB().Table(tableName).Where("id = ?", msgDto.Id).Updates(attrs)
		return updates.RowsAffected, updates.Error
	case "6":
		attrs := map[string]interface{}{}
		attrs["top"] = 1
		updates := sqldb.EdgeDB().Table(tableName).Where("id = ?", msgDto.Id).Updates(attrs)
		return updates.RowsAffected, updates.Error
	case "7":
		attrs := map[string]interface{}{}
		attrs["top"] = 0
		updates := sqldb.EdgeDB().Table(tableName).Where("id = ?", msgDto.Id).Updates(attrs)
		return updates.RowsAffected, updates.Error
	}
	return 0, nil
}

func DeleteSiteLetterOrientationOtherMsg(tableName string, memberId, msgId int64, msgDto *sitemsg.MsgDto) {
	//删除同批次其他ID的消息
	result := sqldb.EdgeDB().Table(tableName).Where("batch_id = ? and site_id = ? and member_id = ? and id <> ?", msgDto.BatchId, msgDto.SiteId, memberId, msgId).
		Delete(&SiteLetterOrientation{})
	if result.Error != nil {
		glog.Errorf("UpdateSiteLetterOrientation OperationType=4 删除同批次 ID=%v 以外的消息 失败! params:batchId=%v siteId=%v memberId=%v error:%+v", msgId, msgDto.BatchId, msgDto.SiteId, memberId, result.Error)
	} else {
		glog.Infof("UpdateSiteLetterOrientation OperationType=4 删除同批次 ID=%v 以外的消息 成功! params:batchId=%v siteId=%v memberId=%v 删除条数:%v", msgId, msgDto.BatchId, msgDto.SiteId, memberId, result.RowsAffected)
	}
}

// GetOrientationByBatchIdAndSiteId 根据批次号和站点id查询
func GetOrientationByBatchIdAndSiteId(tableName string, msgDto *sitemsg.MsgDto) (siteLetterOrientation *SiteLetterOrientation, err error) {
	memberId, _ := strconv.ParseInt(ValueConvert(msgDto.MemberId), 10, 64)
	err = sqldb.EdgeDB().Table(tableName).Where("batch_id = ? and site_id = ? and member_id = ?",
		msgDto.BatchId, msgDto.SiteId, memberId).First(&siteLetterOrientation).Error
	return siteLetterOrientation, err
}

func ValueConvert(value interface{}) string {
	if value == nil {
		return ""
	}

	switch value.(type) {
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case string:
		return value.(string)
	case int:
		return strconv.Itoa(value.(int))
	default:
		return ""
	}
}

func GetBroadcastList(data *sitemsg.MsgDataToKafkaBrRead, table string, createTime string, ids []int) ([]*sitemsg.SiteLetterBroadcast, error) {
	var letters []*sitemsg.SiteLetterBroadcast
	var err error

	tx := sqldb.EdgeDB().
		Table("site_letter_broadcast t").
		Select("id,title,content,msgtype,user_system").
		Where("del_flag= 0 and sys_status= 1").
		Where("not exists(select 1 from "+table+" where msg_id=t.id and member_id= ? and site_id= ?)", data.MemberId, data.SiteId)
	if len(ids) > 0 {
		tx = tx.Where("id in (?)", ids)
	} else {
		tx = tx.Where("msgtype= ? and site_id= ? and user_system= ?", data.MsgType, data.SiteId, data.UserSystem)
	}

	if len(data.RegisterTime) > 0 {
		tx = tx.Where("created_at > ?", data.RegisterTime)
	}

	if len(createTime) > 0 {
		tx = tx.Where("created_at > ?", createTime)
	}

	err = tx.Find(&letters).Error
	return letters, err

}

func SaveReadMany(data []*sitemsg.SiteLetterRead, table string) (err error) {
	db := sqldb.EdgeDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		err = db.Table(table).Create(&data).Error
		if err != nil {
			return err
		}
		return nil
	},
	)
	return
}
