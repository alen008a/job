package sitemsg

import (
	"fmt"
	"siteLetterJob/db/redisdb/core"
	sqldb "siteLetterJob/db/sitemsg"
	"siteLetterJob/internal/context"
	"siteLetterJob/lib/randid"
	mdata2 "siteLetterJob/mdata"
	"siteLetterJob/mdata/rediskey"
	"siteLetterJob/mdata/sitemsg"
	utils2 "siteLetterJob/utils"
	"strings"
	"time"
)

// TemplateConsumer 站内模板消息消费发送,消费 site_msg_template_topic
func MsgBrConsumer(version string, data, serverInfo []byte) {
	var siteId string
	sInfo := string(serverInfo)
	if len(sInfo) > 0 {
		sInfos := strings.Split(sInfo, ":")
		if len(strings.Split(sInfo, ":")) > 0 {
			siteId = sInfos[0]
		}
	}
	c := &context.Context{
		Trace:        randid.GenerateId(),
		SiteId:       siteId,
		KafkaVersion: version,
		ServerInfo:   string(serverInfo),
	}

	defer recoverHandleC(c)

	msgBr := new(sitemsg.MsgDataToKafkaBrRead)
	err := mdata2.Cjson.Unmarshal(data, msgBr)
	if err != nil {
		c.Errorf("MsgBrConsumer 用户广播消息处理解析json错误！ err=%v data=%s", err, string(data))
		return
	}

	if msgBr.Flag != 1 && msgBr.Flag != 2 {
		c.Errorf(" MsgBrConsumer非法传参！ err=%v data=%s", err, string(data))
		return
	}

	c.Infof("MsgBrConsumer 数据MsgBrConsumero data:%+v", msgBr)

	if msgBr.Flag == 1 {
		MessageReadV2(c, msgBr, make([]int, 0))
	} else if msgBr.Flag == 2 {
		MessageReadV2(c, msgBr, msgBr.Ids)

		ids2 := make([]interface{}, 0)
		for _, v := range msgBr.Ids {
			ids2 = append(ids2, v)
		}

		ids3 := make([]interface{}, 0)
		ids3 = append(ids3, ids2...)
		//移除redis保存的集合
		key := fmt.Sprintf(rediskey.BatchReadBRV2, msgBr.SiteId, msgBr.MemberId, msgBr.MsgType)
		count, err := core.SRem(key, ids3...)
		if err != nil {
			c.Infof("MsgBrConsumer core.SRem count:%v err:%v", count, err)
		}
		refreshUnreadCountV2(c, msgBr)
	}
}

func MessageReadV2(c *context.Context, tmplDto *sitemsg.MsgDataToKafkaBrRead, ids []int) {
	//查询表是否有数据
	table := TableNameSiteLetterReadConvert(int64(tmplDto.MemberId))
	list, err := sqldb.GetBroadcastList(tmplDto, table, "", ids)
	if err != nil {
		c.Errorf("GetBroadcastList发生错误Err:%v,param:%v,tableSuffix:%v", err, tmplDto, table)
		return
	}

	numArr := segment(list, 300)
	for _, nums := range numArr {
		if len(nums) > 0 {
			var (
				letterReadList []*sitemsg.SiteLetterRead
				createdAt      = time.Now().Format(utils2.TimeBarFormat)
				updatedAt      = time.Now().Format(utils2.TimeBarFormat)
			)
			for _, v := range nums {
				letterRead := &sitemsg.SiteLetterRead{
					Category:   1,
					Msgtype:    v.Msgtype,
					MsgId:      v.ID,
					UserSystem: tmplDto.UserSystem,
					MemberId:   tmplDto.MemberId,
					DelFlag:    0,
					SiteId:     tmplDto.SiteId,
					CreatedAt:  createdAt,
					UpdatedAt:  updatedAt,
				}
				letterReadList = append(letterReadList, letterRead)
			}

			//如果没有放到已读表
			err = sqldb.SaveReadMany(letterReadList, table)
			if err != nil {
				c.Errorf("AddRead发生错误Err:%v,nums=%v", err, nums)
				return
			}
		}
	}
}

// 数据分段
func segment(list []*sitemsg.SiteLetterBroadcast, subGroupLength int64) [][]*sitemsg.SiteLetterBroadcast {
	max := int64(len(list))
	var segmens = make([][]*sitemsg.SiteLetterBroadcast, 0)
	quantity := max / subGroupLength
	remainder := max % subGroupLength
	i := int64(0)
	for i = int64(0); i < quantity; i++ {
		segmens = append(segmens, list[i*subGroupLength:(i+1)*subGroupLength])
	}
	if quantity == 0 || remainder != 0 {
		segmens = append(segmens, list[i*subGroupLength:i*subGroupLength+remainder])
	}
	return segmens
}

func refreshUnreadCountV2(c *context.Context, req *sitemsg.MsgDataToKafkaBrRead) {
	//站内信第二版，广播消息计数，新版本标识开关
	memberBroadcastUnreadCount := fmt.Sprintf(rediskey.UnreadMsgCacheMemberBroadcastUnreadCount, req.SiteId, req.MemberId, req.MsgType)
	_ = core.SetExpireKV(memberBroadcastUnreadCount, "-1", time.Hour*24)

	memberBroadcastListCount := fmt.Sprintf(rediskey.UnreadMsgCacheMemberBroadcastListCount, req.SiteId, req.MemberId, req.MsgType)
	_ = core.SetExpireKV(memberBroadcastListCount, "-1", time.Hour*24)
}
