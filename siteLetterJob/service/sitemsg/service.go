package sitemsg

import (
	"fmt"
	memberdb "siteLetterJob/db/member"
	"siteLetterJob/db/redisdb/core"
	sqldb "siteLetterJob/db/sitemsg"
	"siteLetterJob/internal/context"
	"siteLetterJob/lib/randid"
	"siteLetterJob/mdata"
	"siteLetterJob/mdata/sitemsg"
	"siteLetterJob/utils"
	"strconv"
	"strings"
	"time"
)

var (
	BinarySum      = 4094
	AllGrade       = 8190
	BinaryNum      = make(map[string]int, 0) // 平台计算map
	VipGrade       = make(map[string]int, 0) // vip等级计算map
	UnreadMsgCount = "site_msg_unread_ms_count_v5_"
	//UnreadMsgCountNew                    = "site_msg_unread_ms_count_v3_"
	UnreadMsgCountNew                    = "site_msg_unread_ms_count_v3_v2_"
	UnreadMsgCacheMember                 = "site_msg_unread_cache_member_v3_"
	UnreadMsgCacheMemberOrientationCount = "site_msg_cache_member_orientation_count_v4_%v_%v_%v" //定向站内信 站点id，memberId,msgType
	UnreadMsgCacheBroadcastCount         = "site_msg_cache_broadcast_count_v3_%v_%v"             //广播站内信 站点id,msgType
)

// 尝试获取锁
func TryGetlock(prefix string, keyText string) bool {
	key := fmt.Sprintf("Lock_%s_%s", prefix, keyText)
	valid, err := core.SetNX(key, "true", 10*time.Second)
	if err != nil {
		return false
	}
	return valid
}

// 释放锁
func UnLock(prefix string, keyText string) {
	key := fmt.Sprintf("Lock_%s_%s", prefix, keyText)
	err := core.DelKey(key)
	if err != nil {
		return
	}
}

func recoverHandleC(c *context.Context) {
	if err := recover(); err != nil {
		tmpStr := fmt.Sprintf("err=%v panic ==> %s\n", err, utils.PanicTrace(4)) //4KB
		c.Errorf(tmpStr)
	}
}

// SiteLetterConsumeHandle 消费站内信处理，消费message-site-msg
func SiteLetterConsumeHandle(version string, data, serverInfo []byte) {
	if len(BinaryNum) == 0 {
		BinaryNum["0"] = 2   //全站
		BinaryNum["1"] = 4   //体育
		BinaryNum["2"] = 8   //web
		BinaryNum["3"] = 16  //h5
		BinaryNum["9"] = 256 //全站体育
	}
	if len(VipGrade) == 0 {
		VipGrade["0"] = 2
		VipGrade["1"] = 4
		VipGrade["2"] = 8
		VipGrade["3"] = 16
		VipGrade["4"] = 32
		VipGrade["5"] = 64
		VipGrade["6"] = 128
		VipGrade["7"] = 256
		VipGrade["8"] = 512
		VipGrade["9"] = 1024
		VipGrade["10"] = 2048
		VipGrade["11"] = 4096
	}
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

	c.Infof("消费站内信处理, 获取到kafka的数据=%s", string(data))

	msgList := make([]*sitemsg.MsgDto, 0)
	err := mdata.Cjson.Unmarshal(data, &msgList)
	if err != nil {
		c.Errorf("消费站内信处理解析json错误！err=%v data=%s", err, string(data))
		return
	}
	if msgList == nil || len(msgList) == 0 {
		c.Infof("消费站内信处理数据为空，处理结束")
		return
	}
	// 过滤OperationType 为空的数据
	msgList = FilterSiteMsg(msgList)
	if len(msgList) == 0 {
		c.Infof("消费站内信数据过滤后size: %v 处理数据为空，处理结束", len(msgList))
		return
	}

	// 幂等性校验
	validate := IdempotentValidate(msgList, c.SiteId)
	if !validate {
		c.Infof("幂等性校验不通过  data=%s", string(data))
		return
	}
	c.Infof("幂等性校验通过，开始处理消息发送 dataSize:%d", len(msgList))

	sendStationLetter(c, msgList)
}

func FilterSiteMsg(msgList []*sitemsg.MsgDto) []*sitemsg.MsgDto {
	msgFilterList := make([]*sitemsg.MsgDto, 0)
	for _, dto := range msgList {
		if dto.OperationType != "" {
			msgFilterList = append(msgFilterList, dto)
		}
	}
	return msgFilterList
}

// 站内信新增修改，删除，启用，禁用
func sendStationLetter(c *context.Context, msgList []*sitemsg.MsgDto) {
	// 操作类型 1-新增，2-修改，3-删除，4-启用,5-停用，6-置顶，7-取消置顶
	if "1" == msgList[0].OperationType {
		msgAdd(c, msgList)
		return
	}

	for _, dto := range msgList {
		if dto.BatchId == "" {
			c.Infof("批次号为空，传入对象= %v", dto)
			continue
		}

		var msgUpdateDto sitemsg.MsgDto
		msgUpdateDto.BatchId = dto.BatchId
		msgUpdateDto.GroupType = dto.GroupType
		msgUpdateDto.MemberId = dto.MemberId
		msgUpdateDto.SiteId = dto.SiteId
		msgUpdateDto.OperationType = dto.OperationType
		switch dto.OperationType {
		case "2": // 修改
			msgUpdateDto.Title = dto.Title
			msgUpdateDto.Content = dto.Content
			msgUpdateDto.IconType = dto.IconType
			msgUpdateDto.Msgtype = dto.Msgtype
			msgUpdateDto.SysStatus = dto.SysStatus
			msgUpdateDto.PcPath = dto.PcPath
			msgUpdateDto.PcUrl = dto.PcUrl
			msgUpdateDto.H5Path = dto.H5Path
			msgUpdateDto.H5Url = dto.H5Url
			msgUpdateDto.JumpUrlType = dto.JumpUrlType
			msgUpdateDto.Sort = dto.Sort
			msgUpdateDto.ImgTop = dto.ImgTop

			c.Infof("2修改站内信：%+v", msgUpdateDto)
			msgUpdate(c, &msgUpdateDto)
			break
		case "3": // 删除
			msgUpdateDto.DelFlag = 1
			row := msgUpdate(c, &msgUpdateDto)
			if row > 0 {
				updateUnreadCount(c, dto, -1)
			}
			break
		case "4": // 启用
			msgUpdateDto.SysStatus = 1
			c.Infof("4启用站内信站内信：%+v", msgUpdateDto)
			row := msgUpdate(c, &msgUpdateDto)
			if row > 0 {
				updateUnreadCount(c, dto, 1)
			}
			break
		case "5": // 停用
			msgUpdateDto.SysStatus = 0
			c.Infof("5停用站内信站内信：%+v", msgUpdateDto)
			row := msgUpdate(c, &msgUpdateDto)
			if row > 0 {
				updateUnreadCount(c, dto, -1)
			}
			break
		case "6": // 置顶
			msgUpdateDto.Top = 1
			msgUpdate(c, &msgUpdateDto)
			break
		case "7": // 取消置顶
			msgUpdateDto.Top = 0
			msgUpdate(c, &msgUpdateDto)
			break
		}
	}
}

func msgAdd(c *context.Context, msgList []*sitemsg.MsgDto) {
	msgDto := msgList[0]
	//新增站内信，由于kafka没有sys_statsu这个字段，现在直接写死，为1
	msgDto.SysStatus = 1
	if msgDto.GroupType == 3 || msgDto.GroupType == 2 {
		// 平台转换
		if "-1" == msgDto.PushPlatform || msgDto.PushPlatform == "" {
			msgDto.Platform = BinarySum
		} else {
			pushPlatforms := strings.Split(msgDto.PushPlatform, ",")
			msgDto.Platform = 0
			for i := range pushPlatforms {
				msgDto.Platform = msgDto.Platform + BinaryNum[pushPlatforms[i]]
			}
		}
		// vip等级转换,只针对广播消息处理
		if msgDto.VipGrade == "" || msgDto.VipGrade == "99" {
			msgDto.VipGradeNum = AllGrade
		} else {
			vipGrades := strings.Split(msgDto.VipGrade, ",")
			for i := range vipGrades {
				msgDto.VipGradeNum = msgDto.VipGradeNum + VipGrade[vipGrades[i]]
			}
		}
		// 根据条件判断发送在线用户或者全部用户
		if msgDto.GroupType == 2 {
			c.Infof("发送在线用户站内信， msgDto=%v", msgDto)
			onlineKey := fmt.Sprintf("%d_member_sso_login_status", msgDto.SiteId)
			onlineStr, err := core.GetKey(onlineKey)
			if err != nil {
				c.Warnf("发送在线用户站内信, 获取redis错误  key:%s  err=%v", onlineKey, err)
				return
			}
			var onlineArr []int64
			err = mdata.Cjson.Unmarshal([]byte(onlineStr), &onlineArr)
			if err != nil {
				c.Warnf("msgAdd 反序列化json失败, str=%s", onlineStr)
				return
			}
			c.Infof("获取当前站点在线会员为：%v", onlineArr)
			if len(onlineArr) == 0 {
				c.Infof("获取当前在线会员为空 处理结束 data:%v", msgDto)
				return
			}

			memberIds := make([]int64, 0)
			for i := range onlineArr {
				currentMemberId := onlineArr[i]
				memberIds = append(memberIds, currentMemberId)
				msgDto.MemberId = int(currentMemberId)
				msgDto.GroupType = 1
				result, msgId, err := sqldb.InsertSiteLetterOrientation(TableNameConvert(currentMemberId), msgDto)
				if err != nil {
					c.Errorf("msgAdd InsertSiteLetterOrientation err:%v", msgDto)
					return
				}
				if result > 0 {
					// 极光推送消息
					go PushMsg(c, memberIds, msgDto, 0)
					c.Infof("站内消息:%v 入库成功，极光推送处理", msgDto)

					//修改缓存未读消息量
					updateUnreadCount(c, msgDto, 1)
					sqldb.DeleteSiteLetterOrientationOtherMsg(TableNameConvert(currentMemberId), currentMemberId, msgId, msgDto)
				}
				memberIds = make([]int64, 0)
			}
		}

		// 全部用户
		if msgDto.GroupType == 3 {
			if msgList[0].IconType == "" {
				msgList[0].IconType = "1"
			}
			c.Infof("发送全部用户站内信：%+v", msgDto)
			result, err := sqldb.InsertSiteLetterBroadcast("site_letter_broadcast", msgDto)
			if err != nil {
				c.Errorf("msgAdd InsertSiteLetter site_letter_broadcast err:%v", msgDto)
				return
			}
			if result > 0 {
				// 推送消息
				c.Infof("站内消息:%+v 入库成功，极光推送处理", msgDto)
				go PushMsg(c, make([]int64, 0), msgDto, 1)
				//修改缓存未读消息量
				updateUnreadCount(c, msgDto, 1)
				sqldb.DeleteLetterBroadcastOtherMsg(msgDto)
			}
		}
	} else {
		// 定向发送
		memberIds := make([]int64, 0)
		c.Infof("定向发送消息数据：%+v", msgList)
		for _, dto := range msgList {
			memberId, _ := strconv.ParseInt(ValueConvert(dto.MemberId), 10, 64)
			if len(msgDto.Content) != 0 && dto.Content == "" {
				dto.Content = msgDto.Content
			}
			// 平台转换
			if "-1" == msgDto.PushPlatform || msgDto.PushPlatform == "" {
				msgDto.Platform = BinarySum
			} else {
				pushPlatforms := strings.Split(msgDto.PushPlatform, ",")
				msgDto.Platform = 0
				for j := range pushPlatforms {
					msgDto.Platform = msgDto.Platform + BinaryNum[pushPlatforms[j]]
				}
			}

			if dto.Platform == 0 {
				dto.Platform = BinarySum
			}
			c.Infof("站内信优化上传图片数据推送数据:%+v", dto)
			if dto.IconType == "" {
				dto.IconType = "1"
			}
			dto.SysStatus = 1
			result, msgId, err := sqldb.InsertSiteLetterOrientation(TableNameConvert(memberId), dto)
			if err != nil {
				c.Errorf("msgAdd InsertSiteLetterOrientation err:%v", err)
				return
			}
			c.Infof("定向发送插入定向消息发送表 InsertSiteLetterOrientation RowsAffected = %d (>0 success!)", result)
			if result > 0 {
				memberIds = append(memberIds, memberId)
				go PushMsg(c, memberIds, dto, 0)
				//修改缓存未读消息量
				updateUnreadCount(c, dto, 1)
				sqldb.DeleteSiteLetterOrientationOtherMsg(TableNameConvert(memberId), memberId, msgId, msgDto)
				memberIds = make([]int64, 0)
			}
		}
	}
}

func msgUpdate(c *context.Context, updateDto *sitemsg.MsgDto) int64 {
	memberId, err := strconv.ParseInt(ValueConvert(updateDto.MemberId), 10, 64)
	if err != nil {
		c.Errorf("修改广播表解析用户名错误 error:%v data:%+v", err.Error(), updateDto)
		return 0
	}
	if updateDto.GroupType == 3 {
		// 修改广播表
		row, err := sqldb.UpdateSiteLetterBroadcast(updateDto)
		if err != nil {
			c.Errorf("修改广播表 UpdateSiteLetterBroadcast error:%v data:%+v", err.Error(), updateDto)
			return 0
		}
		return row
	}
	// 修改定向发送用户表
	if updateDto.MemberId == 0 {
		c.Infof("当前修改定向发送信息失败,会员id为空 %+v", updateDto)
		return 0
	}
	//先查后改
	siteLetterOrientation, err := sqldb.GetOrientationByBatchIdAndSiteId(TableNameConvert(memberId), updateDto)
	if err != nil {
		c.Errorf("GetOrientationByBatchIdAndSiteId query err=%v", err)
		return 0
	}
	if siteLetterOrientation != nil && siteLetterOrientation.Id > 0 {
		updateDto.Id = fmt.Sprintf("%d", siteLetterOrientation.Id)
		result, err := sqldb.UpdateSiteLetterOrientation(TableNameConvert(memberId), memberId, updateDto)
		if err != nil {
			c.Errorf("UpdateSiteLetterOrientation update err=%v", err)
			return 0
		}
		return result
	}
	return 0
}

// 异步修改未读消息量
func updateUnreadCount(c *context.Context, msgDto *sitemsg.MsgDto, add int) {
	c.Infof("异步修改未读消息量 msgDto=:%v, add=%v", msgDto, add)
	if msgDto.PushPlatform == "" {
		msgDto.PushPlatform = "-1"
	}
	// 指定用户发送
	if msgDto.GroupType == 1 {
		{
			//站内信第二版，新版本标识开关
			memberOrientationSwtichKey := fmt.Sprintf(UnreadMsgCacheMemberOrientationCount, msgDto.SiteId, msgDto.MemberId, msgDto.Msgtype)
			core.Incr(memberOrientationSwtichKey)
			_ = core.SetExpireKey(memberOrientationSwtichKey, time.Hour*24)
		}

		//判断当前消息是否为已读状态或者删除状态
		//启用，停用的时候查询当前用户是否已经读取或者删除该消息
		redisKey := fmt.Sprintf("%s%d_%s", UnreadMsgCountNew, msgDto.SiteId, strconv.Itoa(msgDto.MemberId))
		redisKeyOld := fmt.Sprintf("%s%d_%s", UnreadMsgCount, msgDto.SiteId, strconv.Itoa(msgDto.MemberId))

		if msgDto.OperationType == "4" || msgDto.OperationType == "5" {
			// 操作类型1-新增，2-修改，3-删除，4-启用,5-停用，6-置顶，7-取消置顶
			c.Infof("异步修改未读消息量,清除用户缓存 操作类型为 4-启用 or 5-停用，dto=%v", msgDto)
			_ = core.DelKey(redisKey)
			_ = core.DelKey(redisKeyOld)
			return
		}
		// 用户体系1-会员，2-代理
		if msgDto.UserSystem == "1" {
			c.Infof("会员 异步修改未读消息量 dto=%v", msgDto)
			// 推送平台 0全站，1体育,2Web,3 H5,7棋牌，6彩票，8真人，9全站体育App
			if msgDto.PushPlatform == "-1" {
				c.Infof("推送会员全平台，异步修改未读消息量 dto=%+v, ,memberId=%+v", msgDto, msgDto.MemberId)

				hash, err := core.HGetAlL(UnreadMsgCountNew)
				if err != nil {
					c.Errorf("异步修改未读消息量 获取Redis异常, UnreadMsgCountNew=%+v err=%v", UnreadMsgCountNew, err)
					return
				}
				if len(hash) <= 0 {
					c.Infof("单会员未读消息量key 获取为空")
					return
				}

				var hashMap = make(map[string]interface{})
				for key, val := range hash {
					if key == "register_time" || key == "agent_register_time" || key == "-1" {
						continue
					}
					msgCountVo := sqldb.LetterUnreadForTypeRespV2{}
					err = mdata.Cjson.Unmarshal([]byte(val), &msgCountVo)
					if err != nil {
						c.Errorf("hash map Unmarshal to msgCountVo err=%v", err)
						continue
					}

					if msgDto.Msgtype == 1 {
						msgCountVo.NoticeCount = 0
						if msgCountVo.NoticeCount+add >= 0 {
							msgCountVo.NoticeCount = msgCountVo.NoticeCount + add
						}
					}

					if msgDto.Msgtype == 2 {
						msgCountVo.ActivityCount = 0
						if msgCountVo.ActivityCount+add >= 0 {
							msgCountVo.ActivityCount = msgCountVo.ActivityCount + add
						}
					}

					if msgDto.Msgtype == 3 {
						msgCountVo.BulletinCount = 0
						if msgCountVo.BulletinCount+add >= 0 {
							msgCountVo.BulletinCount = msgCountVo.BulletinCount + add
						}
					}
					if msgDto.Msgtype == 4 {
						msgCountVo.MatchCount = 0
						if msgCountVo.MatchCount+add >= 0 {
							msgCountVo.MatchCount = msgCountVo.MatchCount + add
						}
					}
					if msgDto.Msgtype == 5 {
						msgCountVo.Fdcount = 0
						if msgCountVo.Fdcount+add >= 0 {
							msgCountVo.Fdcount = msgCountVo.Fdcount + add
						}
					}
					json, err := mdata.Cjson.MarshalToString(msgCountVo)
					if err != nil {
						c.Errorf("msgCountVo MarshalToString err")
						continue
					}
					hashMap[key] = json
				}

				c.Infof("重新塞回redis 值为：%+v", hashMap)
				err = core.HMSet(redisKey, hashMap)
				if err != nil {
					c.Errorf("调用重新塞入key错误。 redisKey=%+v err=%v", redisKey, err)
					return
				}
			} else {
				var platforms = strings.Split(msgDto.PushPlatform, ",")
				for _, platform := range platforms {
					msgCountVoStr, err := core.HGet(redisKey, platform)
					if err != nil {
						c.Errorf("core.HGetAlL error，err=%v", err)
						continue
					}
					c.Infof("单会员未读消息量key %v  result=%s", redisKey, msgCountVoStr)
					if msgCountVoStr == "" {
						c.Errorf("单会员未读消息量key 获取结果为空，redisKey=%+v ", redisKey)
						continue
					}

					msgCountVo := sqldb.LetterUnreadForTypeRespV2{}
					err = mdata.Cjson.Unmarshal([]byte(msgCountVoStr), &msgCountVo)
					if err != nil || &msgCountVo == nil {
						c.Errorf("hash map Unmarshal to msgCountVo err=%v", err)
						continue
					}

					if msgDto.Msgtype == 1 {
						msgCountVo.NoticeCount = 0
						if msgCountVo.NoticeCount+add >= 0 {
							msgCountVo.NoticeCount = msgCountVo.NoticeCount + add
						}
					}

					if msgDto.Msgtype == 2 {
						msgCountVo.ActivityCount = 0
						if msgCountVo.ActivityCount+add >= 0 {
							msgCountVo.ActivityCount = msgCountVo.ActivityCount + add
						}
					}

					if msgDto.Msgtype == 3 {
						msgCountVo.BulletinCount = 0
						if msgCountVo.BulletinCount+add >= 0 {
							msgCountVo.BulletinCount = msgCountVo.BulletinCount + add
						}
					}
					if msgDto.Msgtype == 4 {
						msgCountVo.MatchCount = 0
						if msgCountVo.MatchCount+add >= 0 {
							msgCountVo.MatchCount = msgCountVo.MatchCount + add
						}
					}
					if msgDto.Msgtype == 5 {
						msgCountVo.Fdcount = 0
						if msgCountVo.Fdcount+add >= 0 {
							msgCountVo.Fdcount = msgCountVo.Fdcount + add
						}
					}
					marshal, err := mdata.Cjson.Marshal(&msgCountVo)
					if err != nil {
						c.Errorf("updateUnreadCount Marshal msgCountVo data=%+v err=%v ", msgCountVo, err)
						continue
					}
					err = core.HSet(redisKey, platform, marshal)
					if err != nil {
						c.Errorf("updateUnreadCount SetHash reqData=%+v err=%v", redisKey, err)
					}
				}
			}
		} else if msgDto.UserSystem == "2" {
			// 代理会员发送
			msgCount, err := core.HGet(redisKey, "-1")
			if err != nil {
				c.Errorf("获取代理会员发送Redis中转服务错误 redisKey=%+v err=%v", redisKey, err)
				return
			}
			msgCountVo := sqldb.LetterUnreadForTypeRespV2{}
			err = mdata.Cjson.Unmarshal([]byte(msgCount), &msgCountVo)
			if err != nil {
				c.Errorf("msgCount Unmarshal err redisKey=%+v err=%v", redisKey, msgCount, err)
				return
			}
			if msgDto.Msgtype == 1 {
				msgCountVo.NoticeCount = 0
				if msgCountVo.NoticeCount+add >= 0 {
					msgCountVo.NoticeCount = msgCountVo.NoticeCount + add
				}
			}

			if msgDto.Msgtype == 2 {
				msgCountVo.ActivityCount = 0
				if msgCountVo.ActivityCount+add >= 0 {
					msgCountVo.ActivityCount = msgCountVo.ActivityCount + add
				}
			}

			if msgDto.Msgtype == 3 {
				msgCountVo.BulletinCount = 0
				if msgCountVo.BulletinCount+add >= 0 {
					msgCountVo.BulletinCount = msgCountVo.BulletinCount + add
				}
			}
			if msgDto.Msgtype == 4 {
				msgCountVo.MatchCount = 0
				if msgCountVo.MatchCount+add >= 0 {
					msgCountVo.MatchCount = msgCountVo.MatchCount + add
				}
			}
			if msgDto.Msgtype == 5 {
				msgCountVo.Fdcount = 0
				if msgCountVo.Fdcount+add >= 0 {
					msgCountVo.Fdcount = msgCountVo.Fdcount + add
				}
			}
			marshal, err := mdata.Cjson.Marshal(msgCountVo)
			if err != nil {
				c.Errorf("updateUnreadCount Marshal msgCountVo data=%v err=%v ", msgCountVo, err)
				return
			}
			err = core.HSet(redisKey, "-1", marshal)
			if err != nil {
				c.Errorf("updateUnreadCount Set Agent Type Hash redisKey=%+v err=%v", redisKey, err)
			}
		}
	} else {
		{
			//站内信第二版，广播消息计数，新版本标识开关
			memberBroadcastCount := fmt.Sprintf(UnreadMsgCacheBroadcastCount, msgDto.SiteId, msgDto.Msgtype)
			core.Incr(memberBroadcastCount)
		}

		// 广播消息处理
		broadcastUnreadCount(c, msgDto, add)
	}
}

// 广播消息处理
func broadcastUnreadCount(c *context.Context, msgDto *sitemsg.MsgDto, add int) {
	c.Infof("多线程处理广播消息未读消息量开始， 处理数据:%+v", msgDto)
	setRedisKey := UnreadMsgCacheMember + strconv.Itoa(msgDto.SiteId)
	// 未读消息redisKey
	redisKey := fmt.Sprintf("%s%d_%s", UnreadMsgCountNew, msgDto.SiteId, strconv.Itoa(msgDto.MemberId))
	redisKeyOld := fmt.Sprintf("%s%d_%s", UnreadMsgCount, msgDto.SiteId, strconv.Itoa(msgDto.MemberId))

	// 取出当前所有缓存会员
	memberIds, err := core.SMembers(setRedisKey)
	if len(memberIds) == 0 || err != nil {
		c.Errorf("取出当前所有缓存会员结果=%v 为空或错误，handle end. reqData=%+v err=%v", memberIds, setRedisKey, err)
		return
	}

	for _, id := range memberIds {
		idInt, err := strconv.Atoi(id)
		if err != nil || idInt <= 0 {
			continue
		}

		//启用，停用的时候重置用户未读已读消息
		value := id
		if msgDto.OperationType == "4" || msgDto.OperationType == "5" {
			//启用停用消息，直接清除用户缓存
			_ = core.DelKey(redisKey)
			_ = core.DelKey(redisKeyOld)

			_, err := core.SRem(setRedisKey, value)
			if err != nil {
				c.Errorf("remove set类型失败 key=%+v value=%s err=%v ", setRedisKey, value, err)
				continue
			}
			//判断当前用户是否在对应vip等级中
			if msgDto.VipGrade != "" && msgDto.VipGrade != "99" {
				//查询当前会员vip等级
				grade, err := memberdb.GetMemberVipGrade(idInt)
				if grade < 0 || err != nil {
					c.Errorf("broadcastUnreadCount GetMemberVipGrade err=%v", err)
					continue
				}
				//当前用户vip等级不在发送vip等级范围内,不做进一步计算
				gradeStr := strconv.Itoa(grade)
				if !((msgDto.VipGradeNum & VipGrade[gradeStr]) == VipGrade[gradeStr]) {
					continue
				}
			}

			//如果广播消息发送时间晚于会员注册时间，则不做缓存更新
			if msgDto.OperationType != "1" {
				//会员未读消息量处理
				if msgDto.UserSystem == "1" {
					regHash, err := core.HGet(redisKey, "register_time")
					c.Infof("GetHash redisKey=%v  field=register_time result=%+v", redisKey, regHash)
					if regHash == "" || err != nil {
						c.Errorf("会员】GetHash 结果为空或出错 redisKey=%v result=%v err=%v", redisKey, regHash, err)
						continue
					}
					registerTime, _ := time.ParseInLocation(utils.TimeBarFormat, regHash, utils.GetBjTimeLoc())
					regUnix := registerTime.Unix()
					sendUnix := msgDto.SendTime
					if sendUnix < regUnix {
						c.Infof("广播消息发送时间:%v 晚于会员:%v 注册时间:%v ", msgDto.SendTime, msgDto.MemberId, registerTime.Unix())
						continue
					}
				} else {
					//代理未读消息量处理
					regHash, err := core.HGet(redisKey, "agent_register_time")
					c.Infof("代理未读消息量处理 redisKey=%v result=%v", redisKey, regHash)
					if regHash == "" || err != nil {
						c.Errorf("【代理】GetHash 结果为空或出错 req=%+v result=%v err=%v", redisKey, regHash, err)
						continue
					}
					registerTime, _ := time.ParseInLocation(utils.TimeBarFormat, regHash, utils.GetBjTimeLoc())
					regUnix := registerTime.Unix()
					sendUnix := msgDto.SendTime
					if sendUnix < regUnix {
						c.Infof("广播消息发送时间:%v 晚于代理会员:%v 注册时间:%v ", msgDto.SendTime, msgDto.MemberId, registerTime.Unix())
						continue
					}
				}
			}
			//会员全部平台
			if "-1" == msgDto.PushPlatform && "1" == msgDto.UserSystem {
				hash, err := core.HGetAlL(redisKey)
				if err != nil {
					c.Errorf("异步修改未读消息量 获取Redis异常, req=%+v err=%v", redisKey, err)
					return
				}
				if len(hash) <= 0 {
					c.Infof("单会员未读消息量key 获取为空")
					return
				}

				var hashMap = make(map[string]interface{})
				for key, val := range hash {
					if key == "register_time" || key == "agent_register_time" || key == "-1" {
						continue
					}
					msgCountVo := sqldb.LetterUnreadForTypeRespV2{}
					err = mdata.Cjson.Unmarshal([]byte(val), &msgCountVo)
					if err != nil {
						c.Errorf("hash map Unmarshal to msgCountVo err=%v", err)
						continue
					}

					if msgDto.Msgtype == 1 {
						msgCountVo.BRNoticeCount = 0
						if msgCountVo.BRNoticeCount+add >= 0 {
							msgCountVo.BRNoticeCount = msgCountVo.BRNoticeCount + add
						}
					}

					if msgDto.Msgtype == 2 {
						msgCountVo.BRActivityCount = 0
						if msgCountVo.BRActivityCount+add >= 0 {
							msgCountVo.BRActivityCount = msgCountVo.BRActivityCount + add
						}
					}

					if msgDto.Msgtype == 3 {
						msgCountVo.BRBulletinCount = 0
						if msgCountVo.BRBulletinCount+add >= 0 {
							msgCountVo.BRBulletinCount = msgCountVo.BRBulletinCount + add
						}
					}
					if msgDto.Msgtype == 4 {
						msgCountVo.BRMatchCount = 0
						if msgCountVo.BRMatchCount+add >= 0 {
							msgCountVo.BRMatchCount = msgCountVo.BRMatchCount + add
						}
					}
					if msgDto.Msgtype == 5 {
						msgCountVo.BRFdcount = 0
						if msgCountVo.BRFdcount+add >= 0 {
							msgCountVo.BRFdcount = msgCountVo.BRFdcount + add
						}
					}
					json, err := mdata.Cjson.MarshalToString(&msgCountVo)
					if err != nil {
						c.Errorf("hash map Unmarshal to msgCountVo err=%v", err)
						continue
					}
					hashMap[key] = json
				}

				c.Infof("重新塞回redis 值为：%+v", hashMap)

				err = core.HMSet(redisKey, hashMap)
				if err != nil {
					c.Errorf("重新塞入key错误。 reqData=%+v err=%v", redisKey, err)
					return
				}
				c.Infof("塞入redis会员全部平台成功！发送站内信处理完毕")
			} else if msgDto.PushPlatform == "-1" && msgDto.UserSystem == "2" { //代理默认全部平台
				// 代理会员发送
				msgCount, err := core.HGet(redisKey, "-1")
				if err != nil {
					c.Errorf("获取代理会员发送 redisKey=%v err=%v", redisKey, err)
					return
				}
				msgCountVo := sqldb.LetterUnreadForTypeRespV2{}
				err = mdata.Cjson.Unmarshal([]byte(msgCount), &msgCountVo)
				if err != nil || &msgCountVo == nil {
					c.Errorf("hash map Unmarshal to msgCountVo err=%v", err)
					return
				}

				if msgDto.Msgtype == 1 {
					msgCountVo.BRNoticeCount = 0
					if msgCountVo.BRNoticeCount+add >= 0 {
						msgCountVo.BRNoticeCount = msgCountVo.BRNoticeCount + add
					}
				}

				if msgDto.Msgtype == 2 {
					msgCountVo.BRActivityCount = 0
					if msgCountVo.BRActivityCount+add >= 0 {
						msgCountVo.BRActivityCount = msgCountVo.BRActivityCount + add
					}
				}

				if msgDto.Msgtype == 3 {
					msgCountVo.BRBulletinCount = 0
					if msgCountVo.BRBulletinCount+add >= 0 {
						msgCountVo.BRBulletinCount = msgCountVo.BRBulletinCount + add
					}
				}
				if msgDto.Msgtype == 4 {
					msgCountVo.BRMatchCount = 0
					if msgCountVo.BRMatchCount+add >= 0 {
						msgCountVo.BRMatchCount = msgCountVo.BRMatchCount + add
					}
				}
				if msgDto.Msgtype == 5 {
					msgCountVo.BRFdcount = 0
					if msgCountVo.BRFdcount+add >= 0 {
						msgCountVo.BRFdcount = msgCountVo.BRFdcount + add
					}
				}
				marshal, err := mdata.Cjson.MarshalToString(msgCountVo)
				if err != nil {
					c.Errorf("updateUnreadCount Marshal msgCountVo data=%v err=%v ", msgCountVo, err)
					return
				}
				err = core.HSet(redisKey, "-1", marshal)
				if err != nil {
					c.Errorf("updateUnreadCount Set Agent Type Hash redisKey=%+v err=%v", redisKey, err)
				}
			} else {
				//会员指定平台
				var platforms = strings.Split(msgDto.PushPlatform, ",")
				for _, platform := range platforms {
					msgCountVoStr, err := core.HGet(redisKey, platform)
					c.Infof("获取redis key=%v platform=%v msgCountVoStr%+v", redisKey, platform, msgCountVoStr)
					if err != nil {
						c.Errorf("获取失败，redisKey=%v err=%v", redisKey, err)
						continue
					}
					if msgCountVoStr == "" {
						c.Errorf("获取结果为空，redisKey=%+v ", redisKey)
						continue
					}

					msgCountVo := sqldb.LetterUnreadForTypeRespV2{}
					err = mdata.Cjson.Unmarshal([]byte(msgCountVoStr), &msgCountVo)
					if err != nil || &msgCountVo == nil {
						c.Errorf("hash map Unmarshal to msgCountVo err=%v", err)
						continue
					}

					if msgDto.Msgtype == 1 {
						msgCountVo.BRNoticeCount = 0
						if msgCountVo.BRNoticeCount+add >= 0 {
							msgCountVo.BRNoticeCount = msgCountVo.BRNoticeCount + add
						}
					}

					if msgDto.Msgtype == 2 {
						msgCountVo.BRActivityCount = 0
						if msgCountVo.BRActivityCount+add >= 0 {
							msgCountVo.BRActivityCount = msgCountVo.BRActivityCount + add
						}
					}

					if msgDto.Msgtype == 3 {
						msgCountVo.BRBulletinCount = 0
						if msgCountVo.BRBulletinCount+add >= 0 {
							msgCountVo.BRBulletinCount = msgCountVo.BRBulletinCount + add
						}
					}
					if msgDto.Msgtype == 4 {
						msgCountVo.BRMatchCount = 0
						if msgCountVo.BRMatchCount+add >= 0 {
							msgCountVo.BRMatchCount = msgCountVo.BRMatchCount + add
						}
					}
					if msgDto.Msgtype == 5 {
						msgCountVo.BRFdcount = 0
						if msgCountVo.BRFdcount+add >= 0 {
							msgCountVo.BRFdcount = msgCountVo.BRFdcount + add
						}
					}
					marshal, err := mdata.Cjson.MarshalToString(msgCountVo)
					if err != nil {
						c.Errorf("updateUnreadCount Marshal msgCountVo data=%v err=%v ", msgCountVo, err)
						continue
					}
					err = core.HSet(setRedisKey, platform, marshal)
					if err != nil {
						c.Errorf("updateUnreadCount SetHash setRedisKey=%v err=%v :", setRedisKey, err)
					}
				}
			}
		}
	}
}
func TableNameConvert(memberId int64) string {
	return fmt.Sprintf("site_letter_orientation_%d", memberId&63)
}

func TableNameSiteLetterReadConvert(memberId int64) string {
	return fmt.Sprintf("site_letter_read_%d", memberId&63)
}
