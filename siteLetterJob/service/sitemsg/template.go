package sitemsg

import (
	"fmt"
	"siteLetterJob/config"
	memberdb "siteLetterJob/db/member"
	"siteLetterJob/db/redisdb/core"
	"siteLetterJob/internal/context"
	"siteLetterJob/internal/glog"
	"siteLetterJob/lib/randid"
	"siteLetterJob/lib/snowid"
	mdata2 "siteLetterJob/mdata"
	"siteLetterJob/mdata/sitemsg"
	"siteLetterJob/utils"
	"strconv"
	"strings"
	"time"
)

// TemplateConsumer 站内模板消息消费发送,消费 site_inner_msg_template_topic
func TemplateConsumer(version string, data, serverInfo []byte) {
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

	msgTmplList := make([]sitemsg.MsgTemplateDto, 0)
	err := mdata2.Cjson.Unmarshal(data, &msgTmplList)
	if err != nil {
		c.Errorf("消费站内信模板消息处理解析json错误！ err=%v data=%s", err, string(data))
	}
	if msgTmplList == nil || len(msgTmplList) == 0 {
		c.Infof("消费站内信处理数据为空，处理结束")
		return
	}
	// 幂等性校验
	validate := IdempotentValidate(msgTmplList, c.SiteId)
	if !validate {
		c.Infof("幂等性校验不通过  data=%s", string(data))
		return
	}
	c.Infof("幂等性校验通过, 本次消费模板数据开始:%+v", string(data))
	msgHandle(c, msgTmplList)
}

func msgHandle(c *context.Context, tmplDto []sitemsg.MsgTemplateDto) {
	c.Infof("数据msgTemplateDto len:=%d", len(tmplDto))
	msgDtoList := make([]*sitemsg.MsgDto, 0)
	siteInnerTemplateVo := &memberdb.SiteInnerMessageTemplate{}
	for _, dto := range tmplDto {
		// 查询当前模板, 默认使用 0 系统模板
		if dto.TemplateType == 0 {
			//查询使用中的系统类型 模板
			siteInnerTemplate, err := memberdb.GetSiteInnerMessageTemplate(dto.TemplateNo)
			if err != nil {
				c.Errorf("msgHandle GetSiteInnerMessageTemplate query err=%v", err)
				continue
			}
			if siteInnerTemplate == nil {
				c.Infof("当前模板不存在,传入的模板编号=%v req=%+v", dto.TemplateNo, dto)
				continue
			}
			c.Infof("根据编号：%v 查询使用中的系统类型模板成功, data=%+v", dto.TemplateNo, dto)
			if siteInnerTemplate.SysStatus == 0 {
				c.Infof("当前模板已停用,传入的模板编号=%v req=%+v", dto.TemplateNo, dto)
				continue
			}
			siteInnerTemplateVo = siteInnerTemplate
		}
		//查询当前会员账号对应的账号
		memberId, _ := strconv.ParseInt(ValueConvert(dto.MemberId), 10, 64)
		memberInfo, err := memberdb.GetMemberInfo(memberId)
		if memberInfo == nil || err != nil {
			time.Sleep(time.Second * 3)
			memberInfo, err = memberdb.GetMemberInfo(memberId)
			if memberInfo == nil || err != nil {
				c.Errorf("msgHandle GetMemberInfo query is null or error. memberId=%d error=%v", memberId, err)
				continue
			}
		}
		if memberInfo.Name == "" {
			c.Errorf("msgHandle GetMemberInfo query is null or error. memberId=%d error=%v", memberId, err)
			continue
		}
		msgDto := buildMsgDto(c, siteInnerTemplateVo, &dto, memberInfo.Name)
		msgDtoList = append(msgDtoList, msgDto)
	}
	// 站内信新增修改，删除，启用，禁用
	if len(msgDtoList) > 0 {
		sendStationLetter(c, msgDtoList)
	}
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

func buildMsgDto(c *context.Context, vo *memberdb.SiteInnerMessageTemplate, msgTemplateDto *sitemsg.MsgTemplateDto, name string) *sitemsg.MsgDto {
	batchId := msgTemplateDto.BatchId
	if msgTemplateDto.Msgtype == 0 {
		msgTemplateDto.Msgtype = 1
	}

	if batchId == "" {
		batchId = strconv.FormatInt(snowid.SnowflakeId(), 10)
	}
	siteId, _ := strconv.Atoi(msgTemplateDto.SiteId)
	memberId, _ := strconv.Atoi(msgTemplateDto.MemberId)
	msgDto := sitemsg.MsgDto{
		GroupType:     1,
		OperationType: "1",
		SiteId:        siteId,
		IconType:      msgTemplateDto.IconType,
		BatchId:       batchId,
		MemberId:      memberId,
		Platform:      4094,
		Msgtype:       msgTemplateDto.Msgtype,
		UserSystem:    "1",
		Status:        1,
		PcPath:        msgTemplateDto.PcPath,
		PcUrl:         msgTemplateDto.PcUrl,
		H5Path:        msgTemplateDto.H5Path,
		H5Url:         msgTemplateDto.H5Url,
		JumpUrlType:   msgTemplateDto.JumpUrlType,
	}
	// 0系统模板
	if msgTemplateDto.TemplateType == 0 {
		if "11" == vo.SendObjectCode {
			msgDto.UserSystem = "2"
		}
		msgDto.IconType = vo.Icon

		msgTitle := convertTemplate(msgTemplateDto, name, vo.Title, vo.OriginalParams)
		c.Infof("0系统模板站内信[标题]转换 原内容:%s  新内容:%s", vo.Title, msgTitle)
		msgContent := convertTemplate(msgTemplateDto, name, vo.Content, vo.OriginalParams)
		c.Infof("0系统模板站内信[内容]转换 原内容:%s  新内容:%s", vo.Content, msgContent)
		msgDto.Title = msgTitle
		msgDto.Content = msgContent
		return &msgDto
	}
	// 1自定义模板
	if msgTemplateDto.TemplateType == 1 {
		msgTitle := convertTemplate(msgTemplateDto, name, msgTemplateDto.Title, "${用户名},${红利类型},${红利金额}")
		c.Infof("1自定义模板站内信[标题]转换 原内容:%s  新内容:%s", msgTemplateDto.Title, msgTitle)
		msgContent := convertTemplate(msgTemplateDto, name, msgTemplateDto.Content, "${用户名},${红利类型},${红利金额}")
		c.Infof("1自定义模板站内信[内容]转换 原内容:%s  新内容:%s", msgTemplateDto.Content, msgContent)
		msgDto.Title = msgTitle
		msgDto.Content = msgContent
		return &msgDto
	}
	return &msgDto
}

func convertTemplate(msgTemplateDto *sitemsg.MsgTemplateDto, name, content, sourceStr string) string {
	originalParams := strings.Split(sourceStr, ",")
	paramsValues := make([]string, 0)
	paramsValues = append(paramsValues, name)

	//添加自定义开发参数
	if msgTemplateDto.ParamsValues != nil {
		paramsValues = append(paramsValues, msgTemplateDto.ParamsValues...)
	}
	paramMap := make(map[string]string, 0)

	for i, param := range originalParams {
		tmpStr := strings.ReplaceAll(strings.TrimSpace(param), "[${]", "")
		paramVal := strings.ReplaceAll(tmpStr, "[}]", "")
		paramMap[paramVal] = paramsValues[i]
	}
	return substitute(content, paramMap)
}

func substitute(content string, params map[string]string) string {
	if params == nil || len(params) == 0 {
		return content
	}
	for k, v := range params {
		content = strings.ReplaceAll(content, k, v)
	}
	return content
}

// IdempotentValidate 幂等性校验 true:通过  false 不通过
func IdempotentValidate(data interface{}, siteId string) bool {
	dataBytes, err := mdata2.Cjson.Marshal(data)
	if err != nil {
		return false
	}
	md5Str := utils.Md5Encrypt(dataBytes)
	key := fmt.Sprintf("wm:%s:site:msg:new:%s", siteId, md5Str)
	exist, err := core.KeyExist(key)
	if err != nil {
		return false
	}
	if !exist {
		glog.Infof("幂等性校验通过")
		err = core.SetExpireKV(key, "1", 30*time.Second)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

// PushMsg 站内信极光推送
func PushMsg(c *context.Context, memberIds []int64, msgDto *sitemsg.MsgDto, flag int) {
	c.Infof("本次接受消息, 会员集合：%v 推送对象:%+v", memberIds, msgDto)
	if msgDto.PushFlag == 0 {
		c.Infof("推送标记 为关闭状态，处理结束")
		return
	}
	if !config.GetConfig().MsgPushFlag {
		c.Infof("消息推送关闭, 推送处理结束 会员集合=%v,推送对象=%+v", memberIds, msgDto)
		return
	}
	tm := time.Unix(msgDto.SendTime, 0)

	pushReq := sitemsg.SendPushReq{}
	pushReq.NoticeStr = "站内信推送"
	pushReq.MemberIds = memberIds
	pushReq.SiteId = strconv.Itoa(msgDto.SiteId)
	pushReq.MsgType = "notification"
	pushReq.PushTitle = msgDto.Title
	pushReq.PushContent = msgDto.Content
	pushReq.Type = 3
	pushReq.TimeLive = 86400
	pushReq.BusinessTime = tm.Format(utils.TimeBarFormat)
	if flag == 1 {
		// 广播所有人
		pushReq.Type = 1
	}
	platform := msgDto.PushPlatform
	pushReqPlatforms := make([]*sitemsg.PushReqPlatform, 0)
	stringBuffer := ""
	if msgDto.PushDevice != "" {
		devices := strings.Split(msgDto.PushDevice, ",")
		for i, device := range devices {
			de := "ios"
			if device == "1" {
				de = "android"
			}

			if i > 0 {
				stringBuffer += "," + de
			} else {
				stringBuffer += de
			}
			if platform != "" {
				if "-1" == platform {
					platform = "0,1"
				}

				platforms := strings.Split(platform, ",")
				for _, s := range platforms {
					if mdata2.PlatformTypeMap[s] != "" {
						reqPlatform := sitemsg.PushReqPlatform{Platform: mdata2.PlatformTypeMap[s], DeviceType: stringBuffer}
						pushReqPlatforms = append(pushReqPlatforms, &reqPlatform)
					}
				}
				pushReq.PlatformList = pushReqPlatforms
			}
		}
	}
	c.Infof("开始发送极光推送 start ====  参数：%+v", pushReq)
	//先注释，关掉极光推送防止推送到线上
	SendPush(c, &pushReq)
	c.Infof("发送极光推送结束 end 参数：%+v", pushReq)
}
