package service

import (
	"fmt"
	"github.com/agnivade/levenshtein"
	"siteVideoJob/config"
	"siteVideoJob/db/redisdb/core"
	"siteVideoJob/db/sqldb"
	"siteVideoJob/internal/context"
	"siteVideoJob/lib/cache"
	"siteVideoJob/mdata"
	"siteVideoJob/mdata/model"
	"siteVideoJob/mdata/video"
	"siteVideoJob/utils"
	"strconv"
	"strings"
	"time"
)

const (
	venueMatchPath      = "/cond/%s/%s"
	lstyVideoUrlDataKey = "lsty_video_url_data_key"

	crtyMatchVideoUrlFromLsKey = "crty_match_video_lsty_key_%v_%v" //{matchClass}{matchId}
)

func PullCRTYMatchList(c *context.Context) string {
	var msgBuild strings.Builder
	matchClass := config.GetLSTYMatchClassList()
	//获取72小时以内赛事资料
	gameDates := getSpecifyFutureGameDate(5)
	for i := range matchClass {
		for j := range gameDates {
			msg := DoSetCRTYOriginData(c, matchClass[i], gameDates[j])
			msgBuild.WriteString(msg)
			msgBuild.WriteString("\r\n")
		}
	}
	return msgBuild.String()
}

// 获取未来指定天数的日期
func getSpecifyFutureGameDate(days int) []string {
	var gameDates []string
	startDate := utils.GetBjNowTime()
	gameDates = append(gameDates, startDate.Format(utils.TimeBarYYMMDD))
	for i := 1; i <= days; i++ {
		newDate := startDate.Add(time.Duration(i*24) * time.Hour)
		gameDates = append(gameDates, newDate.Format(utils.TimeBarYYMMDD))
	}
	return gameDates
}

// DoSetCRTYOriginData  落地CRTY比赛数据源
func DoSetCRTYOriginData(c *context.Context, matchClass, gameDate string) (msg string) {
	var (
		bkResp video.BkResp
	)
	data, err := GetLSTYDataList(c, fmt.Sprintf(venueMatchPath, mdata.CRTY, matchClass), model.PullAllGameReq{
		GameDate: gameDate,
	})
	if err != nil {
		msg = fmt.Sprintf("DoSetCRTYOriginData  error:%+v", err)
		c.Error(msg)
		return
	}
	if len(data) == 0 {
		return
	}
	err = mdata.Cjson.Unmarshal(data, &bkResp)
	if err != nil {
		msg = fmt.Sprintf("DoSetCRTYOriginData  error:%+v", err)
		c.Error(msg)
		return
	}
	if bkResp.StatusCode != 200 {
		msg = fmt.Sprintf("DoSetCRTYOriginData GetLSTYDataList err  bkResp.StatusCode = %v | bkResp.Status = %v", bkResp.StatusCode, bkResp.Status)
		c.Error(msg)
		return
	}
	for _, v := range bkResp.Data {
		matchCrtyOriginData := &model.MatchCrtyOriginData{
			MatchId:    strconv.Itoa(int(v.Eid)),
			LeagueId:   v.LeagueId,
			MatchName:  v.League,
			VenueName:  mdata.CRTY,
			MatchClass: matchClass,
			StartAt:    v.StartDatetime,
			HomeName:   v.Team1,
			HomeLogo:   v.Team1Logo,
			VisitName:  v.Team2,
			VisitLogo:  v.Team2Logo,
			LeagueLogo: v.LeagueLogo,
			LiveStatus: v.Status,
			CreatedBy:  mdata.System,
			UpdatedBy:  mdata.System,
		}
		affectRow, err := sqldb.UpsertMatchCrtyOriginData(matchCrtyOriginData)
		if err != nil {
			c.Errorf("DoSetCRTYOriginData UpsertMatchCrtyOriginData error:%+v ,param:%+v", err, matchCrtyOriginData)
		} else {
			c.Infof("DoSetCRTYOriginData UpsertMatchCrtyOriginData success affectRow:%+v ,param:%+v", affectRow, matchCrtyOriginData)
		}
	}
	return
}

// DoSetCrtyFromLsty 匹配雷速体育视频源
func DoSetCrtyFromLsty(c *context.Context) (msg string) {
	var (
		lstyVideoGroup       = make(map[string][]*model.MatchVideoUrlData)
		crtyMatchOriginGroup = make(map[string][]*model.MatchCrtyOriginData)
		threshold            = 3 // 设置相似度阈值,阈值越小表示越严格
	)
	startDate := utils.GetBjNowTime()
	startTime := utils.BeginOfDay(utils.GetBjNowTime()).Format(utils.TimeBarFormat)
	endTime := utils.EndOfDay(startDate.Add(3 * 24 * time.Hour)).Format(utils.TimeBarFormat)
	data, err := cache.GetOrSet(c, lstyVideoUrlDataKey, cache.RotatePeriod180, func() (interface{}, error) {
		return sqldb.QueryLSMatchVideoUrlData(startTime, endTime)
	})
	if err != nil {
		msg = fmt.Sprintf("DoSetCrtyFromLsty QueryLSMatchVideoUrlData error:%+v,startTime:%s,endTime:%s", err, startTime, endTime)
		c.Error(msg)
		return msg
	}

	if data == nil {
		return
	}

	matchVideoUrlDatas, ok := data.([]*model.MatchVideoUrlData)
	if !ok {
		msg = "DoSetCrtyFromLsty data not []*model.MatchVideoUrlData type"
		c.Error(msg)
		return msg
	}
	if len(matchVideoUrlDatas) == 0 {
		return
	}

	crOriginMatchDatas, err := sqldb.QueryCRMatchData(startTime, endTime)
	if err != nil {
		msg = fmt.Sprintf("DoSetCrtyFromLsty QueryCRMatchData error:%+v,startTime:%s,endTime:%s", err, startTime, endTime)
		c.Error(msg)
		return msg
	}
	if len(crOriginMatchDatas) == 0 {
		return
	}

	for _, v := range matchVideoUrlDatas {
		if lstyArr, ok := lstyVideoGroup[v.MatchClass]; ok {
			lstyArr = append(lstyArr, v)
			lstyVideoGroup[v.MatchClass] = lstyArr
		} else {
			lstyVideoGroup[v.MatchClass] = []*model.MatchVideoUrlData{v}
		}
	}

	for _, v := range crOriginMatchDatas {
		if crtyArr, ok := crtyMatchOriginGroup[v.MatchClass]; ok {
			crtyArr = append(crtyArr, v)
			crtyMatchOriginGroup[v.MatchClass] = crtyArr
		} else {
			crtyMatchOriginGroup[v.MatchClass] = []*model.MatchCrtyOriginData{v}
		}
	}

	//匹配相同类型的赛事为一组,,️以LSTY为基准
	for k, lstyArr := range lstyVideoGroup {
		if crtyArr, ok := crtyMatchOriginGroup[k]; ok {
			go func(lsty []*model.MatchVideoUrlData, crty []*model.MatchCrtyOriginData) {
				matchCrtyVideoData := findSimilarMatches(c, lsty, crty, threshold)
				if len(matchCrtyVideoData) > 0 {
					for _, v := range matchCrtyVideoData {
						key := fmt.Sprintf(crtyMatchVideoUrlFromLsKey, v.MatchClass, v.MatchId)
						startAt, _ := utils.BjTBarFmtTime(v.StartAt)
						expireDuration := startAt.Add(24 * time.Hour).Sub(utils.GetBjNowTime())
						matchCrtyVideoDataStr, _ := mdata.Cjson.MarshalToString(v)
						core.SetExpireKV(key, matchCrtyVideoDataStr, expireDuration)
						c.Infof("DoSetCrtyFromLsty findSimilarMatches 匹配成功的视频源:%s", matchCrtyVideoDataStr)
					}
					c.Infof("DoSetCrtyFromLsty findSimilarMatches 匹配成功的类型:%s,数量:%d", k, len(matchCrtyVideoData))
				}
			}(lstyArr, crtyArr)
		}
	}
	return
}

func isSimilar(c *context.Context, s1, s2 string, threshold int) bool {
	thres := levenshtein.ComputeDistance(s1, s2)
	c.Infof("isSimilar s1:%s,s2:%s,result:%d", s1, s2, thres)
	return thres <= threshold
}

func findSimilarMatches(c *context.Context, s1 []*model.MatchVideoUrlData, s2 []*model.MatchCrtyOriginData, threshold int) []*model.MatchCrtyVideoData {
	var (
		matches []*model.MatchCrtyVideoData
		status  bool
	)
	for _, item1 := range s1 {
		for _, item2 := range s2 {
			status = item1.StartAt == item2.StartAt
			status = isSimilar(c, item1.HomeName, item2.HomeName, threshold) && status
			status = isSimilar(c, item1.VisitName, item2.VisitName, threshold) && status
			if item1.MatchName != "" && item2.MatchName != "" {
				status = isSimilar(c, item1.MatchName, item2.MatchName, threshold) && status
			}
			if status {
				item := &model.MatchCrtyVideoData{
					MatchId:    item2.MatchId,
					LeagueId:   item2.LeagueId,
					MatchClass: item2.MatchClass,
					StartAt:    item2.StartAt,
					HomeLogo:   item2.HomeLogo,
					VisitLogo:  item2.VisitLogo,
					LeagueLogo: item2.LeagueLogo,
					PushUrl1:   item1.PushUrl1,
					PushUrl2:   item1.PushUrl2,
				}
				if item.HomeLogo == "" {
					item.HomeLogo = item1.HomeLogo
				}
				if item.VisitLogo == "" {
					item.VisitLogo = item1.VisitLogo
				}
				if item.LeagueLogo == "" {
					item.LeagueLogo = item1.LeagueLogo
				}
				matches = append(matches, item)
			}
		}
	}
	return matches
}
