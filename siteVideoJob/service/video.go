package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"siteVideoJob/config"
	"siteVideoJob/db/redisdb/core"
	"siteVideoJob/db/redisdb/game"
	"siteVideoJob/db/sqldb"
	"siteVideoJob/internal/context"
	"siteVideoJob/internal/glog"
	"siteVideoJob/lib/httpclient"
	"siteVideoJob/lib/rp"
	"siteVideoJob/mdata"
	"siteVideoJob/mdata/model"
	"siteVideoJob/mdata/redisKey"
	"siteVideoJob/service/metadata"
	"siteVideoJob/utils"
	"sort"
	"strconv"
	"strings"
	"sync"

	"siteVideoJob/mdata/video"
	"time"
)

type videoUrls []*video.VideoUrl

func (v videoUrls) MarshalBinary() ([]byte, error) {
	jsonBytes, err := mdata.Cjson.Marshal(v)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (v *videoUrls) UnmarshalBinary(b []byte) error {
	return mdata.Cjson.Unmarshal(b, v)
}

type matchDataRedis video.MatchDataVo

func (v matchDataRedis) MarshalBinary() ([]byte, error) {
	jsonBytes, err := mdata.Cjson.Marshal(v)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (v *matchDataRedis) UnmarshalBinary(b []byte) error {
	return mdata.Cjson.Unmarshal(b, v)
}

func DeleteLastDayMatch(c *context.Context) (count int64, err error) {
	var matchIds []int64
	list, err := sqldb.QueryFinishedInfos()
	if err != nil {
		return 0, err
	}
	for _, v := range list {
		matchIds = append(matchIds, v.MatchId)
		delPlatVo(c, v.VenueName, strconv.Itoa(int(v.MatchId)))
		delVideoSourceVo(c, v.VenueName, strconv.Itoa(int(v.MatchId)))
		delVideoMatchMap(c, v.VenueName, strconv.Itoa(int(v.MatchId)))
	}

	count, err = sqldb.BatchDeleteEvents(matchIds)
	if err != nil {
		return
	}
	return
}

func SiteAdminUpdateLiveVideo(c *context.Context) {
	videoSourceApi := config.GetVideoSourceApi()
	for i := range videoSourceApi {
		go func(i int) {
			apiConfig := videoSourceApi[i]
			matchClassArr := apiConfig.MatchClassArr
			for j := range matchClassArr {
				apiConfig.MatchClass = matchClassArr[j]
				liveData := getBkLiveData(c, apiConfig)
				info := createEventInfo(c, liveData, apiConfig.ChannelName)
				c.Infof("场馆:%v|分类:%v|视频源更新到表结果为%v", apiConfig.ChannelName, apiConfig.MatchClass, info)
			}
		}(i)
	}
}

// GetMatchList 获取播控赛事信息
func GetMatchList(c *context.Context, api, matchClass string) (resp []byte, err error) {
	var (
		LiveSourceUrl = config.GetConfig().VideoLiveSourceUrl
		LiveSiteAlias = config.GetConfig().VideoLiveSiteAlias
	)
	if len(api) == 0 {
		c.Errorf("siteVideoJob | GetMatchList err = api参数不能为空")
		return []byte(""), errors.New("api参数不能为空")
	}
	path := LiveSourceUrl + "/video/v1/" + LiveSiteAlias + "/" + api + "?matchClass=" + matchClass
	resp, err = httpclient.ProxyGet(path, map[string]string{mdata.HeaderSite: c.SiteId}, httpclient.GetVideoProxyClient(time.Minute*30))
	if err != nil {
		glog.EmergencyWithTimeout("url=%s |err=%v", 60, path, err)
		return nil, err
	}
	return
}

// GetMatchOddsList 获取赛事赔率信息
func GetMatchOddsList(c *context.Context, api, matchId, matchClass string) (resp []byte, err error) {
	var (
		LiveSourceUrl = config.GetConfig().VideoLiveSourceUrl
		LiveSiteAlias = config.GetConfig().VideoLiveSiteAlias
	)
	if len(api) == 0 {
		c.Errorf("siteVideoJob | GetMatchList err = api参数不能为空")
		return []byte(""), errors.New("api参数不能为空")
	}
	path := LiveSourceUrl + "/video/v1/" + LiveSiteAlias + "/" + api + "/" + matchId + "?matchClass=" + matchClass
	resp, err = httpclient.ProxyGet(path, map[string]string{mdata.HeaderSite: c.SiteId}, httpclient.GetVideoProxyClient(time.Minute*30))
	if err != nil {
		glog.EmergencyWithTimeout("url=%s |err=%v", 60, path, err)
		return nil, err
	}
	return
}

// GetVideoSourceList 实时获取赛事视频源
func GetVideoSourceList(c *context.Context, api, matchId, matchClass string) video.VideoSourceMap {
	var (
		LiveSourceUrl  = config.GetConfig().VideoLiveSourceUrl
		LiveSiteAlias  = config.GetConfig().VideoLiveSiteAlias
		videoSourceUrl = make(video.VideoSourceMap, 0)
		resp           video.VideoUrlSourceResp
	)
	if len(api) == 0 {
		c.Errorf("GetVideoSourceList err = api参数不能为空")
		return videoSourceUrl
	}
	path := LiveSourceUrl + "/video/v1/" + LiveSiteAlias + "/source/" + api + "/" + matchId + "?matchClass=" + matchClass
	data, err := httpclient.ProxyGet(path, map[string]string{mdata.HeaderSite: c.SiteId}, httpclient.GetVideoProxyClient(time.Minute*5))
	if err != nil {
		glog.EmergencyWithTimeout("url=%s |err=%v", 60, path, err)
		return videoSourceUrl
	}

	err = mdata.Cjson.Unmarshal(data, &resp)
	if err != nil {
		c.Errorf("GetVideoSourceList Unmarshal err = %v | api = %v", err, api)
		return videoSourceUrl
	}
	if resp.StatusCode != 200 || resp.Status != "success" {
		c.Errorf("GetVideoSourceList  err resp.StatusCode = %v | resp.Status = %v", resp.StatusCode, resp.Status)
		return videoSourceUrl
	}
	return resp.Data
}

func createEventInfo(c *context.Context, videoSourceVos []video.VideoSourceVo, channelName string) bool {
	var (
		eventsIds    []int64
		insertEvents []model.LiveEventsConfig
	)
	matchNameMap := make(map[int64]string)
	startAtMap := make(map[int64]string)
	homeNameMap := make(map[int64]string)
	visitNameMap := make(map[int64]string)
	homeLogoMap := make(map[int64]string)
	visitLogoMap := make(map[int64]string)
	leagueLogoMap := make(map[int64]string)
	videoSource := make(map[string]interface{})

	displayStartAt := time.Now().AddDate(0, 0, -3).Format(utils.TimeBarFormat)
	displayEndAt := time.Now().AddDate(0, 3, 0).Format(utils.TimeBarFormat)
	for _, vo := range videoSourceVos {
		eventsIds = append(eventsIds, vo.EventId)
		matchNameMap[vo.EventId] = vo.League
		startAtMap[vo.EventId] = vo.StartDate
		homeNameMap[vo.EventId] = vo.Team1
		visitNameMap[vo.EventId] = vo.Team2
		if len(vo.HomeTeamLogoUrl) > 0 {
			homeLogoMap[vo.EventId] = vo.HomeTeamLogoUrl
		}

		if len(vo.GuestTeamLogoUrl) > 0 {
			visitLogoMap[vo.EventId] = vo.GuestTeamLogoUrl
		}

		if len(vo.LeagueLogoUrl) > 0 {
			leagueLogoMap[vo.EventId] = vo.LeagueLogoUrl
		}

		if len(vo.VideoUrls) > 0 {
			videoSource[strconv.Itoa(int(vo.EventId))] = videoUrls(vo.VideoUrls)
		}
	}
	liveEventsV2List, err := sqldb.QueryEventsInfos(0, channelName, eventsIds)
	if err != nil {
		c.Errorf("SiteAdminUpdateLiveVideo QueryEventsInfos err = %v", err)
		return false
	}
	// videoSourceVos 全部比赛，liveEventsV2List 已有比赛
	updateList, insertList := removeLiveEvents(videoSourceVos, liveEventsV2List)
	c.Infof("SiteAdminUpdateLiveVideo updateList条数 = %v | insertList条数=%v", len(updateList), len(insertList))
	updateMap := make(map[int64]video.VideoSourceVo)
	for _, vo := range updateList {
		updateMap[vo.EventId] = vo
	}
	for _, vo := range insertList {
		if vo.EventId == 0 || vo.StartDate == "" {
			continue
		}
		//做二次重复检查
		if sqldb.IsExistEvents(channelName, 0, vo.EventId) {
			continue
		}
		var liveEvents model.LiveEventsConfig
		liveEvents.MatchId = vo.EventId
		liveEvents.MatchName = vo.League
		liveEvents.StartAt = vo.StartDate
		liveEvents.BallClassId = vo.BallClassId
		liveEvents.PlateId = vo.PlateId
		liveEvents.HomeName = vo.Team1
		liveEvents.HomeLogo = vo.HomeTeamLogoUrl
		liveEvents.VisitName = vo.Team2
		liveEvents.VisitLogo = vo.GuestTeamLogoUrl
		liveEvents.VenueName = channelName
		liveEvents.MatchClass = config.GetGameType(channelName, strings.ToUpper(vo.MatchClass))
		leagueId, _ := strconv.Atoi(vo.LeagueId)
		liveEvents.LeagueId = leagueId
		liveEvents.Cate = strconv.Itoa(vo.Cate)
		liveEvents.DisplayStartAt = displayStartAt
		liveEvents.DisplayEndAt = displayEndAt
		liveEvents.LiveStatus = vo.LiveStatus
		liveEvents.LeagueLogo = vo.LeagueLogoUrl
		// 欧洲杯专用
		//if liveEvents.MatchClass == "FT" && strings.Contains(liveEvents.MatchName, "欧洲足球锦标赛2024(在德国)") {
		//	liveEvents.MajorEvent = 1
		//}
		insertEvents = append(insertEvents, liveEvents)
	}
	if len(insertEvents) > 0 {
		affect, err1 := sqldb.BatchInsertEvents(insertEvents)
		if err1 != nil {
			c.Errorf("SiteAdminUpdateLiveVideo BatchInsertEvents 批量插入失败 err = %v", err1)
		} else if affect > 0 {
			c.Infof("SiteAdminUpdateLiveVideo InsertEvent 插入成功总条数=%v", affect)
		}
	}
	for _, eventsConfig := range liveEventsV2List {
		liveEvents := model.LiveEventsConfig{}
		liveEvents.Id = eventsConfig.Id
		liveEvents.PlateId = eventsConfig.PlateId
		updateLiveStatus := updateMap[eventsConfig.MatchId].LiveStatus
		if updateLiveStatus != eventsConfig.LiveStatus {
			liveEvents.LiveStatus = updateLiveStatus
		}
		if len(eventsConfig.MatchNameBy) == 0 && eventsConfig.MatchName != matchNameMap[eventsConfig.MatchId] {
			liveEvents.MatchName = matchNameMap[eventsConfig.MatchId]
		}
		if len(eventsConfig.StartAtBy) == 0 && eventsConfig.StartAt != startAtMap[eventsConfig.MatchId] {
			liveEvents.StartAt = startAtMap[eventsConfig.MatchId]
		}
		if len(eventsConfig.HomeNameBy) == 0 && eventsConfig.HomeName != homeNameMap[eventsConfig.MatchId] {
			liveEvents.HomeName = homeNameMap[eventsConfig.MatchId]
		}
		if len(eventsConfig.HomeLogoBy) == 0 && eventsConfig.HomeLogo != homeLogoMap[eventsConfig.MatchId] {
			liveEvents.HomeLogo = homeLogoMap[eventsConfig.MatchId]
		}
		if len(eventsConfig.VisitNameBy) == 0 && eventsConfig.VisitName != visitNameMap[eventsConfig.MatchId] {
			liveEvents.VisitName = visitNameMap[eventsConfig.MatchId]
		}
		if len(eventsConfig.VisitLogoBy) == 0 && eventsConfig.VisitLogo != visitLogoMap[eventsConfig.MatchId] {
			liveEvents.VisitLogo = visitLogoMap[eventsConfig.MatchId]
		}
		if eventsConfig.LeagueLogo != leagueLogoMap[eventsConfig.MatchId] {
			liveEvents.LeagueLogo = leagueLogoMap[eventsConfig.MatchId]
		}
		if len(liveEvents.MatchName) == 0 && len(liveEvents.StartAt) == 0 && len(liveEvents.HomeName) == 0 && len(liveEvents.HomeLogo) == 0 &&
			len(liveEvents.VisitName) == 0 && len(liveEvents.VisitLogo) == 0 && len(liveEvents.MatchName) == 0 && len(liveEvents.LeagueLogo) == 0 {
			continue
		}

		// 欧洲杯专用
		//if liveEvents.MatchClass == "FT" && strings.Contains(liveEvents.MatchName, "欧洲足球锦标赛2024(在德国)") && liveEvents.MajorEvent != 1 {
		//	liveEvents.MajorEvent = 1
		//}

		v2, err2 := sqldb.UpdateVideo(&liveEvents)
		if err2 != nil {
			c.Errorf("SiteAdminUpdateLiveVideo UpdateVideoV2 err = %v", err2)
			continue
		}
		c.Infof("SiteAdminUpdateLiveVideo 比赛id = %v ｜ 更新的数量为：%v", liveEvents.MatchId, v2)
	}
	batchSetVideoSourceVo(c, channelName, videoSource)
	return true
}

func removeLiveEvents(list []video.VideoSourceVo, exists []model.LiveEventsConfig) ([]video.VideoSourceVo, []video.VideoSourceVo) {
	updateList := make([]video.VideoSourceVo, 0)
	insertList := make([]video.VideoSourceVo, 0)
	existMap := make(map[int64]struct{})

	for _, lec := range exists {
		existMap[lec.MatchId] = struct{}{}
	}

	for _, vo := range list {
		if _, ok := existMap[vo.EventId]; ok {
			updateList = append(updateList, vo)
		} else {
			insertList = append(insertList, vo)
		}
	}
	return updateList, insertList
}

func PullLiveEvents(c *context.Context) {
	videoSourceApi := config.GetVideoSourceApi()
	for i := range videoSourceApi {
		go func(i int) {
			apiConfig := videoSourceApi[i]
			matchClassArr := apiConfig.MatchClassArr
			for j := range matchClassArr {
				apiConfig.MatchClass = matchClassArr[j]
				liveData := getBkLiveData(c, apiConfig)
				matchDataMap := make(map[string]interface{})
				for _, datum := range liveData {
					var matchDataVo video.MatchDataVo
					matchDataVo.MatchStatus = datum.LiveStatus     // 赛事状态
					matchDataVo.HomeLogo = datum.HomeTeamLogoUrl   // 赛事主队队标
					matchDataVo.VisitLogo = datum.GuestTeamLogoUrl // 赛事客队队标
					matchDataVo.StartTime = datum.StartDate        // 开赛时间
					matchDataVo.AnchorStatus = datum.AnchorStatus  // 主播状态
					matchDataVo.AnchorName = datum.AnchorName      // 主播名称
					matchDataVo.LeagueId = datum.LeagueId
					matchDataVo.LeagueLogo = datum.LeagueLogoUrl
					matchDataVo.PlateId = datum.PlateId
					matchDataVo.BallClassId = datum.BallClassId
					vUrls := queryVideoSourceVo(c, datum.VenueName, strconv.Itoa(int(datum.EventId)))
					if len(datum.VideoUrls) >= len(vUrls) || len(vUrls) == 0 {
						matchDataVo.VideoUrls = datum.VideoUrls
					} else {
						matchDataVo.VideoUrls = vUrls
					}
					if len(matchDataVo.VideoUrls) <= 1 && datum.VenueName == mdata.XMTY {
						source := GetVideoSourceList(c, datum.VenueName, strconv.Itoa(int(datum.EventId)), datum.MatchClass)
						if v, ok := source[strconv.Itoa(int(datum.EventId))]; ok {
							if len(v) > 0 {
								matchDataVo.VideoUrls = v
								setVideoSourceVo(c, datum.VenueName, strconv.Itoa(int(datum.EventId)), v)
							}
						}
					}
					matchDataMap[strconv.Itoa(int(datum.EventId))] = matchDataRedis(matchDataVo)
				}
				if len(matchDataMap) == 0 {
					return
				}

				err := core.HMSet(fmt.Sprintf(redisKey.VideoMatchMapKey, apiConfig.ChannelName), matchDataMap)
				if err != nil {
					c.Errorf("PullLiveEvents HMSet err = %v", err)
					continue
				}
				c.Infof("PullLiveEvents | channelName = %v | matchClass= %s | set redis length = %d", apiConfig.ChannelName, apiConfig.MatchClass, len(matchDataMap))
			}
		}(i)
	}
}

func PullLiveOdds(c *context.Context) {
	videoSourceApi := config.GetVideoSourceApi()
	for _, data := range videoSourceApi {
		tempData := data
		rp.Go(func() {
			DoSetVideoOdds(c, tempData)
		})
	}
}

func getBkLiveData(c *context.Context, venueConfig *mdata.VideoSourceApi) []video.VideoSourceVo {
	var (
		liveData      video.BkResp
		videoSourceVo []video.VideoSourceVo
	)
	jsonChar, err := GetMatchList(c, venueConfig.BkApi, venueConfig.MatchClass)
	if err != nil {
		c.Errorf("getBkLiveData GetMatchList err = %v | api = %v | matchClass = %v", err, venueConfig.BkApi, venueConfig.MatchClass)
	}
	err = mdata.Cjson.Unmarshal(jsonChar, &liveData)
	if err != nil {
		c.Errorf("getBkLiveData Unmarshal err = %v | api = %v | matchClass = %v", err, venueConfig.BkApi, venueConfig.MatchClass)
	}
	if liveData.StatusCode != 200 || liveData.Status != "success" {
		c.Errorf("getBkLiveData GetMatchList err liveData.StatusCode = %v | liveData.Status = %v", liveData.StatusCode, liveData.Status)
	}
	for _, datum := range liveData.Data {
		videoSource := video.VideoSourceVo{}
		videoSource.AnchorStatus = datum.AnchorStatus
		videoSource.AnchorName = datum.AnchorName
		videoSource.Cate = datum.Cate
		videoSource.PlateId = datum.PlateId
		videoSource.BallClassId = datum.BallClassId
		videoSource.MatchClass = datum.SportType
		videoSource.EventId = datum.Eid
		videoSource.HomeTeamLogoUrl = datum.Team1Logo
		videoSource.GuestTeamLogoUrl = datum.Team2Logo
		videoSource.League = datum.League
		videoSource.LeagueId = datum.LeagueId
		videoSource.LeagueLogoUrl = datum.LeagueLogo
		videoSource.StartDate = datum.StartDatetime
		videoSource.VenueName = datum.Channel
		videoSource.LiveStatus = datum.Status
		videoSource.StreamId = datum.StreamID
		videoSource.Team1 = datum.Team1
		videoSource.Team2 = datum.Team2
		videoSource.VideoUrls = datum.VideoUrls
		videoSourceVo = append(videoSourceVo, videoSource)
	}
	return videoSourceVo
}

func SetLiveVideoData(c *context.Context) {
	videoSourceApi := config.GetVideoSourceApi()
	for _, data := range videoSourceApi {
		tempData := data
		rp.Go(func() {
			DoSetVideoData(c, tempData.ChannelName)
		})
	}
}

func DoSetVideoOdds(c *context.Context, venueConfig *mdata.VideoSourceApi) {
	var (
		matchGameTypeMap = make(map[string][]string)
		batchSize        = 50
	)
	//获取公共赛事
	eventsList, err := queryLivesList(0, 0, venueConfig.ChannelName)
	if err != nil {
		c.Errorf("DoSetVideoOdds DoSetVideoOdds err = %v", err)
		return
	}
	for _, liveEvents := range eventsList {
		if matchIdArr, ok := matchGameTypeMap[liveEvents.MatchClass]; ok {
			matchIdArr = append(matchIdArr, strconv.Itoa(int(liveEvents.MatchId)))
			matchGameTypeMap[liveEvents.MatchClass] = matchIdArr
		} else {
			matchGameTypeMap[liveEvents.MatchClass] = []string{strconv.Itoa(int(liveEvents.MatchId))}
		}
	}
	c.Infof("DoSetVideoOdds 热门赛事赔率，场馆=%v|数量=%v｜分组结果=%+v", venueConfig.ChannelName, len(eventsList), matchGameTypeMap)
	for gameType, v := range matchGameTypeMap {
		num := int(math.Ceil(float64(len(v)) / float64(batchSize)))
		for i := 0; i < num; i++ {
			var items []string
			start := i * batchSize
			end := (i + 1) * batchSize
			if end > len(v) {
				end = len(v)
			}
			items = v[start:end]
			var oddsData video.BkOddsResp
			jsonChar, err := GetMatchOddsList(c, venueConfig.BkApi, strings.Join(items, ","), gameType)
			if err != nil {
				c.Errorf("getBkLiveData GetMatchList err = %v | api = %v", err, venueConfig.BkApi)
				continue
			}
			err = mdata.Cjson.Unmarshal(jsonChar, &oddsData)
			if err != nil {
				c.Errorf("getBkLiveData Unmarshal err = %v | api = %v", err, venueConfig.BkApi)
				continue
			}
			if oddsData.StatusCode != 200 || oddsData.Status != "success" {
				c.Errorf("getBkLiveData GetMatchList err liveData.StatusCode = %v | liveData.Status = %v", oddsData.StatusCode, oddsData.Status)
				continue
			}
			c.Infof("DoSetVideoOdds 热门赛事赔率，三方拉取数据成功！场馆=%v|分类=%v|数量=%v", venueConfig.ChannelName, gameType, len(items))
			for _, datum := range oddsData.Data {
				setPlatVo(c, venueConfig.ChannelName, datum.Gid, datum)
			}
		}
	}
}

func DoSetVideoData(c *context.Context, channelName string) {
	var (
		eventsForRedis video.LiveEventList
		dataMap        map[string]string
	)
	eventsList, err := queryLivesList(0, 0, channelName)
	if err != nil {
		c.Errorf("SetLiveVideoData DoSetVideoData err = %v", err)
		return
	}
	c.Infof("热门赛事，场馆=%v|数量=%v", channelName, len(eventsList))
	dataMap, err = core.HGetAlL(fmt.Sprintf(redisKey.VideoMatchMapKey, channelName))
	if err != nil {
		c.Errorf("SetLiveVideoData DoSetVideoData GetKey err = %v", err)
		return
	}
	token := queryXMTYUserToken(c, channelName)
	c.Infof("SetLiveVideoData queryXMTYUserToken channelName=%s token=%s", channelName, token)
	for _, liveEvents := range eventsList {
		downTime, err1 := utils.CountDownTime(liveEvents.StartAt)
		if err1 != nil {
			c.Errorf("SetLiveVideoData DoSetVideoData CountDownTime err = %v", err)
			continue
		}
		liveEvents.OpenTime = downTime // 开赛倒计时
		if liveEvents.OpenTime == "0" {
			liveEvents.IsStart = true // 直播是否已开始
		}
		var matchDataVo video.MatchDataVo
		if dataVal, ok := dataMap[strconv.Itoa(int(liveEvents.MatchId))]; ok {
			err = mdata.Cjson.UnmarshalFromString(dataVal, &matchDataVo)
			if err != nil {
				c.Errorf("SetLiveVideoData DoSetVideoData  UnmarshalFromString video.MatchDataVo err = %v,val=%s", err, dataVal)
			}
			liveEvents.PlateId = matchDataVo.PlateId
			liveEvents.LiveStatus = matchDataVo.MatchStatus
			if len(matchDataVo.VideoUrls) != 0 {
				liveEvents.VideoUrls = matchDataVo.VideoUrls
			}
		}
		if len(liveEvents.VideoUrls) == 0 {
			vUrls := queryVideoSourceVo(c, channelName, strconv.Itoa(int(liveEvents.MatchId)))
			if len(vUrls) != 0 {
				liveEvents.VideoUrls = vUrls
			}
		}
		plateVo := queryPlatVo(c, channelName, liveEvents.MatchId)
		liveEvents.GroupBy = utils.Substring(liveEvents.StartAt, 0, 10) // 小时分类
		if plateVo.Gid == "" || plateVo.Kid == "" {
			c.Errorf("SetLiveVideoData DoSetVideoData not plat data channel=%s|matchClass=%s|matchId=%d", channelName, liveEvents.MatchClass, liveEvents.MatchId)
			if liveEvents.PlateId == 0 {
				plateId, liveStatus := convertStartDateToGameTime(liveEvents.StartAt, liveEvents.LiveStatus)
				liveEvents.PlateId = plateId
				liveEvents.LiveStatus = liveStatus
			}
			err = sqldb.UpdateVideoPlatStatus(int(liveEvents.Id), liveEvents.PlateId, liveEvents.LiveStatus)
			if err != nil {
				c.Errorf("SetLiveVideoData DoSetVideoData UpdateVideoPlatStatus  err = %v", err)
			}
			continue
		}
		c.Infof("SetLiveVideoData DoSetVideoData | 成功获取赛事数据 场馆名=%v | 赛事id = %v | 数据值 = %+v", channelName, liveEvents.MatchId, plateVo)
		liveStatus, _ := strconv.Atoi(plateVo.LiveStatus)
		liveEvents.LiveStatus = liveStatus
		liveEvents.MatchName = strings.ReplaceAll(liveEvents.MatchName, "\\*", "")
		liveEvents.HomeName = strings.ReplaceAll(liveEvents.HomeName, "\\*", "")
		liveEvents.VisitName = strings.ReplaceAll(liveEvents.VisitName, "\\*", "")
		p, err1 := strconv.Atoi(plateVo.P)
		if err1 != nil {
			continue
		}
		err = sqldb.UpdateVideoPlatStatus(int(liveEvents.Id), p, liveStatus)
		if err != nil {
			c.Errorf("SetLiveVideoData DoSetVideoData UpdateVideoPlatStatus  err = %v", err)
		}
		liveEvents.PlateId = p
		liveEvents.HomeScore = plateVo.HomeScore
		liveEvents.AwayScore = plateVo.AwayScore
		liveEvents.PeriodStatus = plateVo.PeriodStatus
		liveEvents.GameTime = plateVo.GameTime
		week, err2 := utils.GetWeek(liveEvents.StartAt)
		if err2 != nil {
			liveEvents.Week = "周日"
			c.Errorf("SetLiveVideoData DoSetVideoData GetWeek err = %v", err)
		}
		liveEvents.Week = week

		kid, err3 := strconv.Atoi(plateVo.Kid)
		if err3 != nil {
			c.Errorf("SetLiveVideoData DoSetVideoData Atoi err = %v", err)
			continue
		}
		liveEvents.LeagueId = int64(kid)
		ballClassId, err4 := strconv.Atoi(plateVo.CateId)
		if err4 != nil {
			c.Errorf("SetLiveVideoData DoSetVideoData ballClassId Atoi err = %v", err)
			continue
		}
		liveEvents.BallClassId = ballClassId
		//liveEvents.AnchorStatus = dataMap[liveEvents.MatchId].AnchorStatus
		//liveEvents.AnchorName = dataMap[liveEvents.MatchId].AnchorName

		liveEvents.OddsInfo = plateVo.Odds
		liveEvents.VenueName = channelName
		if matchDataVo.MatchStatus == 2 || plateVo.LiveStatus == "2" { // 直播已结束
			continue
		}
		if roughMatchLiveFinished(liveEvents.MatchClass, liveEvents.StartAt, plateVo) {
			continue
		}
		// 替换播放链接中的token
		liveEvents.VideoUrls = replaceVideoLinkToken(channelName, token, liveEvents.VideoUrls)
		eventsForRedis = append(eventsForRedis, liveEvents)
	}

	if len(eventsForRedis) == 0 {
		return
	}

	siteIds, err := sqldb.QueryALLSite()
	c.Infof("SetLiveVideoData QueryALLSite siteIds:%+v", siteIds)
	if err != nil || len(siteIds) == 0 {
		c.Infof("SetLiveVideoData siteIds length is 0 or QueryALLSite error:%+v", err)
		return
	}
	//使用站点赛事替换公共赛事
	for _, v := range siteIds {
		var (
			basketball video.LiveEventList
			football   video.LiveEventList
			esports    video.LiveEventList
			tennis     video.LiveEventList
			important  video.LiveEventList
		)
		newEventsForRedis := make(video.LiveEventList, 0)
		siteEventsList, err := queryLivesList(v, -1, channelName)
		if err != nil {
			c.Errorf("SetLiveVideoData DoSetVideoData queryLivesListBySiteId  err = %v", err)
		}
		if len(siteEventsList) > 0 {
			for i := range eventsForRedis {
				for j := range siteEventsList {
					//更新站点赛事特有的数据
					if eventsForRedis[i].MatchId == siteEventsList[j].MatchId {
						eventsForRedis[i].HomeName = siteEventsList[j].HomeName
						eventsForRedis[i].HomeLogo = siteEventsList[j].HomeLogo
						eventsForRedis[i].VisitName = siteEventsList[j].VisitName
						eventsForRedis[i].VisitLogo = siteEventsList[j].VisitLogo
						eventsForRedis[i].MatchClass = siteEventsList[j].MatchClass
						eventsForRedis[i].MatchLogo = siteEventsList[j].MatchLogo
						eventsForRedis[i].MatchName = siteEventsList[j].MatchName
						eventsForRedis[i].MajorEvent = siteEventsList[j].MajorEvent
						eventsForRedis[i].Status = siteEventsList[j].Status
						eventsForRedis[i].DelFlag = siteEventsList[j].DelFlag
					}
				}
				if eventsForRedis[i].DelFlag == 0 && eventsForRedis[i].Status == 0 {
					newEventsForRedis = append(newEventsForRedis, eventsForRedis[i])
				}
			}
			sort.Sort(newEventsForRedis)
		} else {
			newEventsForRedis = eventsForRedis
		}

		c.Infof("场馆%s匹配播控的比赛有%d场", channelName, len(newEventsForRedis))
		for _, data := range newEventsForRedis {
			if data.MatchClass == "BK" {
				basketball = append(basketball, data)
			} else if data.MatchClass == "FT" {
				football = append(football, data)
			} else if data.MatchClass == "ESPORTS" {
				esports = append(esports, data)
			} else if data.MatchClass == "TN" {
				tennis = append(tennis, data)
			} else if data.MatchClass == "IMPORT" {
				important = append(important, data)
			}
		}

		err = insertRedisLiveSports(v, channelName, "BK", basketball)
		if err != nil {
			c.Errorf("SetLiveVideoData DoSetVideoData insertRedisLiveSports basketball err = %v", err)
		}
		c.Infof("站点=%v |场馆=%v | 处理篮球数据=%v条", v, channelName, len(basketball))
		err = insertRedisLiveSports(v, channelName, "FT", football)
		if err != nil {
			c.Errorf("SetLiveVideoData DoSetVideoData insertRedisLiveSports football err = %v", err)
		}
		c.Infof("站点=%v |场馆=%v | 处理足球数据=%v条", v, channelName, len(football))
		err = insertRedisLiveSports(v, channelName, "ESPORTS", esports)
		if err != nil {
			c.Errorf("SetLiveVideoData DoSetVideoData insertRedisLiveSports esports err = %v", err)
		}
		c.Infof("站点=%v |场馆=%v | 处理电竞数据=%v条", v, channelName, len(esports))
		err = insertRedisLiveSports(v, channelName, "TN", tennis)
		if err != nil {
			c.Errorf("SetLiveVideoData DoSetVideoData insertRedisLiveSports tennis err = %v", err)
		}
		c.Infof("站点=%v |场馆=%v | 处理网球数据=%v条", v, channelName, len(tennis))
		err = insertRedisLiveSports(v, channelName, "IMPORT", important)
		if err != nil {
			c.Errorf("SetLiveVideoData DoSetVideoData insertRedisLiveSports important err = %v", err)
		}
		c.Infof("站点=%v | 场馆=%v | 处理重要赛事数据=%v条", v, channelName, len(important))
	}
}

// 目前无需查询场馆维护状态。
func _(c *context.Context, channelName string) (*video.VenueStatusInfo, error) {
	// 检查场馆状态
	var (
		venueStatusInfo *video.VenueStatusInfo
	)
	redisData, err := core.GetOrSet(fmt.Sprintf(redisKey.LiveVenueStatus, channelName, metadata.GetSiteId(c)), func() (interface{}, error) {
		Info, err := sqldb.QueryVenueStatus(channelName)
		if err != nil || Info == nil {
			return nil, err
		}
		return Info, err
	}, 30000)
	if err != nil {
		c.Errorf("SetLiveVideoEvents | DoSetVideoData | getVenueStatusByEnName | GetOrSet err = %v | redisData = %v", err, string(redisData))
		return nil, err
	}
	err = json.Unmarshal(redisData, &venueStatusInfo)
	if err != nil {
		c.Errorf("SetLiveVideoEvents | DoSetVideoData | getVenueStatusByEnName | Unmarshal err = %v | redisData = %v", err, string(redisData))
		return nil, err
	}
	return venueStatusInfo, err
}

func convertStartDateToGameTime(startDate string, liveStatus int) (plateId, newLiveStatus int) {
	startTime, _ := utils.BjTBarFmtTime(startDate)
	now := utils.GetBjNowTime()
	if startTime.Before(now) && startTime.Format(utils.TimeBarYYMMDD) == utils.GetBjNowTime().Format(utils.TimeBarYYMMDD) {
		return 1, 0
	}
	if liveStatus == 0 && startTime.After(utils.EndOfDay(now)) {
		return 2, 0
	}
	if startTime.Before(now) && now.Before(startTime.Add(3*time.Hour)) && liveStatus == 1 {
		return 4, 1
	}
	if liveStatus == 2 {
		return 5, 2
	}
	if now.After(startTime.Add(3 * time.Hour)) {
		return 5, 2
	}
	return 1, liveStatus
}

// 大概判断比赛是否结束,预计一般足球比赛155分钟,篮球比赛100分钟
func roughMatchLiveFinished(matchClass, startTime string, plateVo video.PlateVo) bool {
	startDate, _ := utils.BjTBarFmtTime(startTime)
	if matchClass == "FT" {
		if strings.Contains(plateVo.League, "独家") || strings.Contains(plateVo.League, "FIFA") {
			startDate = startDate.Add(17 * time.Minute)
		} else {
			startDate = startDate.Add(120 * time.Minute)
		}
	}
	if matchClass == "BK" {
		startDate = startDate.Add(100 * time.Minute)
	}
	// 是否已过比赛时间
	state := utils.GetBjNowTime().After(startDate)
	if !state {
		return false
	}
	// 已经结束比赛
	if plateVo.LiveStatus == "2" {
		return true
	}
	// 盘口赔率是否都为0
	var (
		isO  bool
		odds map[string]map[string]string
	)
	oddsByte, _ := mdata.Cjson.Marshal(plateVo.Odds)
	_ = mdata.Cjson.Unmarshal(oddsByte, &odds)
	for _, v := range odds {
		for _, v1 := range v {
			vv0, _ := strconv.Atoi(v1)
			if vv0 == 0 {
				isO = true
			} else {
				isO = false
			}
		}
	}
	if state && isO {
		return true
	}
	return false
}

// 获取用户token
func queryXMTYUserToken(c *context.Context, venueName string) string {
	if venueName != mdata.XMTY {
		return ""
	}
	tokenStr, err := game.GetKey(redisKey.XMTYVideoTokenKey)
	if err != nil && !errors.Is(err, game.RedisNil) {
		c.Errorf("queryXMTYUserToken key=%s err=%+v", redisKey.XMTYVideoTokenKey, venueName)
		return ""
	}
	var loginData video.LoginData
	err = mdata.Cjson.UnmarshalFromString(tokenStr, &loginData)
	if err != nil {
		c.Errorf("queryXMTYUserToken UnmarshalFromString err=%+v", err)
		return ""
	}
	return loginData.Token
}

func replaceVideoLinkToken(venueName, token string, videoUrls []*video.VideoUrl) []*video.VideoUrl {
	if venueName != mdata.XMTY || token == "" {
		return videoUrls
	}
	if len(videoUrls) == 0 {
		return videoUrls
	}
	newVideoUrls := make([]*video.VideoUrl, 0)
	for _, v := range videoUrls {
		if v.VideoType == "1" && v.PlayType == "p" {
			purePath, err := url.Parse(v.Path)
			if err != nil {
				newVideoUrls = append(newVideoUrls, v)
				continue
			}
			value := purePath.Query()
			value.Set("token", token)
			path := url.URL{
				Scheme:   purePath.Scheme,
				Host:     purePath.Host,
				Path:     purePath.Path,
				RawQuery: value.Encode(),
			}
			vi := &video.VideoUrl{
				VideoType: v.VideoType,
				Path:      path.String(),
				PlayType:  v.PlayType,
			}
			newVideoUrls = append(newVideoUrls, vi)
		} else {
			newVideoUrls = append(newVideoUrls, v)
		}
	}
	return newVideoUrls
}

func queryLivesList(siteId, delStatus int, venueName string) (list video.LiveEventList, err error) {
	var req video.LiveEventsDto
	req.SiteId = siteId
	req.Status = 0
	req.DelStatus = delStatus
	req.VenueName = venueName
	req.StartAt = time.Now().Add(-time.Hour * 2).Format(utils.TimeBarFormat)
	list, err = sqldb.LiveEventsList(req)
	if err != nil {
		return
	}
	return
}

func queryPlatVo(c *context.Context, venueName string, matchId int64) (p video.PlateVo) {
	mapKey := "sport_event:" + venueName
	valueKey := fmt.Sprintf("%s_%d", venueName, matchId)
	str, err := core.HGet(mapKey, valueKey)
	if err != nil || str == "" {
		return
	}
	err = mdata.Cjson.Unmarshal([]byte(str), &p)
	if err != nil {
		c.Errorf("SetLiveVideoData DoSetVideoData queryPlatVo err = %v", err)
	}
	return
}

func setPlatVo(c *context.Context, venueName, matchId string, plateVo video.PlateVo) {
	mapKey := "sport_event:" + venueName
	valueKey := fmt.Sprintf("%s_%s", venueName, matchId)
	if roughMatchLiveFinished(plateVo.Category, plateVo.OpenTime, plateVo) {
		plateVo.LiveStatus = "2"
	}
	plateVoBytes, _ := mdata.Cjson.Marshal(plateVo)
	err := core.HSet(mapKey, valueKey, plateVoBytes)
	if err != nil {
		c.Errorf("setPlatVo err = %v", err)
		return
	}
	return
}

func delPlatVo(c *context.Context, venueName, matchId string) {
	mapKey := "sport_event:" + venueName
	valueKey := fmt.Sprintf("%s_%s", venueName, matchId)
	err := core.HDel(mapKey, valueKey)
	if err != nil {
		c.Errorf("setPlatVo err = %v", err)
		return
	}
	return
}

func queryVideoSourceVo(c *context.Context, venue, matchId string) (videoSourceVo []*video.VideoUrl) {
	mapKey := fmt.Sprintf(redisKey.VenueSourceKey, venue)
	str, err := core.HGet(mapKey, matchId)
	if err != nil || str == "" {
		return
	}
	err = mdata.Cjson.UnmarshalFromString(str, &videoSourceVo)
	if err != nil {
		c.Errorf("queryVideoSourceVo err = %v", err)
	}
	return
}

func setVideoSourceVo(c *context.Context, venue, matchId string, videoSourceVo []*video.VideoUrl) {
	mapKey := fmt.Sprintf(redisKey.VenueSourceKey, venue)
	plateVoBytes, _ := mdata.Cjson.Marshal(videoSourceVo)
	err := core.HSet(mapKey, matchId, plateVoBytes)
	if err != nil {
		c.Errorf("setVideoSourceVo err = %v", err)
		return
	}
	return
}

func batchSetVideoSourceVo(c *context.Context, venue string, videoSource map[string]interface{}) {
	mapKey := fmt.Sprintf(redisKey.VenueSourceKey, venue)
	err := core.HMSet(mapKey, videoSource)
	if err != nil {
		c.Errorf("batchSetVideoSourceVo err = %v", err)
		return
	}
	return
}

func delVideoSourceVo(c *context.Context, venue string, matchId ...string) {
	mapKey := fmt.Sprintf(redisKey.VenueSourceKey, venue)
	err := core.HDel(mapKey, matchId...)
	if err != nil {
		c.Errorf("setVideoSourceVo err = %v", err)
		return
	}
	return
}

func delVideoMatchMap(c *context.Context, venue string, matchId ...string) {
	mapKey := fmt.Sprintf(redisKey.VideoMatchMapKey, venue)
	err := core.HDel(mapKey, matchId...)
	if err != nil {
		c.Errorf("delVideoMatchMap err = %v", err)
		return
	}
	return
}

func insertRedisLiveSports(siteId int, venueName string, sport string, data video.LiveEventList) (err error) {
	key := fmt.Sprintf(redisKey.VenueLiveSportsKey, venueName, siteId, sport)
	err = core.DelKey(key)
	if err != nil {
		return
	}
	if len(data) > 0 {
		marshalList := make([]interface{}, 0)
		for i, event := range data {
			// 最多保留300条数据
			if i > 300 {
				break
			}
			marshal, _ := mdata.Cjson.MarshalToString(event)
			marshalList = append(marshalList, utils.Compress([]byte(marshal)))
		}
		// 改为分批插入
		batchSize := 50
		end := 0
		for i := 0; end < len(marshalList); i++ {
			start := i * batchSize
			end = (i + 1) * batchSize
			if end >= len(marshalList) {
				end = len(marshalList)
			}
			err = core.RPush(key, marshalList[start:end]...)
			if err != nil {
				return
			}
		}
	}
	return err
}

// 从播控中获取主播场次信息
func PullAnchorEvents(c *context.Context) {
	var (
		anchorSourceApi []mdata.VideoSourceApi
	)

	err := mdata.Cjson.UnmarshalFromString(config.GetConfig().AnchorSourceApi, &anchorSourceApi)
	if err != nil {
		c.Errorf("PullLiveEvents Unmarshal err = %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(anchorSourceApi))
	for _, data := range anchorSourceApi {
		// 必须要使用一个临时变量将data拷贝过来再传递给协程
		tempData := data
		rp.Go(func() {
			defer wg.Done()
			eventIds := getBkAnchorData(c, tempData)
			key := fmt.Sprintf(redisKey.AnchorLiveVenue, metadata.GetSiteId(c))

			// 初始化默认场馆 默认wm + im
			if venueScore, err := core.ZSCORE(key, tempData.ChannelName); err != nil {
				var redisError error
				if tempData.ChannelName == "wm" || tempData.ChannelName == "im" {
					redisError = core.ZAdd(key, tempData.ChannelName, 2)
				} else {
					redisError = core.ZAdd(key, tempData.ChannelName, 1)
				}
				if redisError != nil {
					c.Infof("set LIVE_VENUE redis error:%s", err)
				}
			} else if venueScore == 2 {
				// 设置缺省数据
				if len(eventIds) > 0 {
					err := core.HMSet(fmt.Sprintf(redisKey.AnchorShowDefault, metadata.GetSiteId(c)), eventIds)
					if err != nil {
						c.Infof("set ANCHOR_SHOW_DEFAULT error:%s", err)
					}
				} else {
					// 博控为空
					oldKey, err := core.HKeys(fmt.Sprintf(redisKey.AnchorShowDefault, metadata.GetSiteId(c)))
					if err != nil {
						c.Errorf("ANCHOR_SHOW_DEFAULT not have err:%v", err)
					}
					if err == nil && len(oldKey) > 0 {
						err = core.HDel(fmt.Sprintf(redisKey.AnchorShowDefault, metadata.GetSiteId(c)), oldKey...)
						if err != nil {
							c.Errorf("reset ANCHOR_SHOW_DEFAULT  err:%v", err)
						}
					}
				}
			}
		})
	}
	wg.Wait()
	err = updateShowAnchorData(metadata.GetSiteId(c))
	if err != nil {
		c.Errorf("更新前台主播接口数据失败,error:%s", err)
	}
	c.Infof("生成前台主播接口数据成功")
}

func getBkAnchorData(c *context.Context, venueConfig mdata.VideoSourceApi) map[string]interface{} {
	var (
		anchorData video.BkAnchorResp
		// anchorListMap  = make(map[string]video.AnchorListVo) //给后台管理列表使用
		anchorEventMap = make(map[string]interface{}) // 给前台展示 结构同播控接口
		eventIDs       = make(map[string]interface{})
	)

	jsonChar, err := GetMatchList(c, venueConfig.BkApi, venueConfig.MatchClass)
	if err != nil {
		c.Errorf("getBkAnchorData GetMatchList err = %v | api = %v", err, venueConfig.BkApi)
	}
	err = mdata.Cjson.Unmarshal(jsonChar, &anchorData)
	if err != nil {
		c.Errorf("getBkAnchorData Unmarshal err = %v | api = %v", err, venueConfig.BkApi)
	}
	if anchorData.StatusCode != 200 || anchorData.Status != "success" {
		c.Errorf("getBkAnchorData GetMatchList err anchorData.StatusCode = %v | anchorData.Status = %v", anchorData.StatusCode, anchorData.Status)
		return nil
	}

	allEventIDs, _ := core.HKeys(fmt.Sprintf(redisKey.AnchorEventsList, metadata.GetSiteId(c)))
	defaultIDs, _ := core.HKeys(fmt.Sprintf(redisKey.AnchorShowDefault, metadata.GetSiteId(c)))
	allEventIDsMap := make(map[string]int)
	for _, id := range append(allEventIDs, defaultIDs...) {
		allEventIDsMap[id] = 1
	}
	// 两种数据 not_started和started
	// 主播数据路径：data.not_started[0].events[0].anchorVideo[0].anchor.nickname
	for _, events := range anchorData.Data { // not_started or started
		for _, eventDatum := range events {
			for _, eventInfo := range eventDatum.Events { // events[0]
				eventInfo.VenueName = venueConfig.ChannelName
				eventTime, _ := time.ParseInLocation("2006-01-02 15:04:05.000", eventInfo.StartDate, utils.GetBjTimeLoc())
				eventInfo.TimeStamp = eventTime.Unix()

				for _, anchorData := range eventInfo.AnchorVideo {
					anchorInfo := video.AnchorListVo{}
					anchorInfo.VID = venueConfig.ChannelName + "-" + eventInfo.Eid
					anchorInfo.EventID = eventInfo.Eid
					anchorInfo.VenueName = venueConfig.ChannelName
					anchorInfo.League = eventInfo.League
					anchorInfo.StartDate = eventInfo.StartDate
					anchorInfo.Team1 = eventInfo.Team1
					anchorInfo.Team2 = eventInfo.Team2
					anchorInfo.Cate = strings.ToLower(eventInfo.Cate)
					anchorInfo.AnchorsID = anchorData.Anchors.ID
					anchorInfo.AnchorStatus = anchorData.AnchorStatus
					anchorInfo.Nickname = anchorData.Anchors.Nickname
					anchorInfo.LogoSquare = anchorData.Anchors.LogoSquare

					// 将主播信息切片拆分成单条记录
					var tempEventInfo video.AnchorEvents
					// 用反序列化实现深拷贝
					data, _ := json.Marshal(eventInfo)
					err := json.Unmarshal(data, &tempEventInfo)
					if err != nil {
						return nil
					}
					tempEventInfo.AnchorVideo = tempEventInfo.AnchorVideo[:0]
					tempEventInfo.AnchorVideo = append(tempEventInfo.AnchorVideo, anchorData)
					anchorEventMap[eventInfo.VenueName+"-"+eventInfo.Eid+"-"+anchorData.Anchors.ID] = tempEventInfo
					eventIDs[eventInfo.VenueName+"-"+eventInfo.Eid+"-"+anchorData.Anchors.ID] = eventInfo.TimeStamp
					delete(allEventIDsMap, eventInfo.VenueName+"-"+eventInfo.Eid+"-"+anchorData.Anchors.ID)
				}
			}
		}
	}

	if len(anchorEventMap) == 0 {
		c.Infof("getBkAnchorData 视频源=%v返回数据为空", venueConfig.ChannelName)
	} else {
		err = saveAnchorEventList(metadata.GetSiteId(c), anchorEventMap)
		if err != nil {
			c.Errorf("update Anchor redis err = %v", err)
		}
		// 根据vid删除不在最新播控数据中的场次
		if len(allEventIDsMap) > 0 {
			deleteKeys := make([]string, 0, 0)
			for k := range allEventIDsMap {
				if strings.HasPrefix(k, venueConfig.ChannelName) {
					deleteKeys = append(deleteKeys, k)
				}
			}
			if len(deleteKeys) > 0 {
				DeleteAnchorLive(c, deleteKeys...)
			}
			c.Infof("delete expired event :%d", len(deleteKeys))
		}
	}
	c.Infof("PullAnchorLiveEvents | channelName = %v | set redis length = %d", venueConfig.ChannelName, len(anchorEventMap))
	return eventIDs
}

func saveAnchorEventList(siteId int, data map[string]interface{}) (err error) {
	key := fmt.Sprintf(redisKey.AnchorEventsList, siteId)
	if len(data) > 0 {
		err = core.HMSet(key, data)
		if err != nil {
			return
		}
	}
	return err
}

func updateShowAnchorData(siteId int) error {
	var anchorEvent video.AnchorEvents
	resp := make(map[string]video.AnchorList)
	showAnchorIDs, err := core.HKeys(fmt.Sprintf(redisKey.AnchorShowIds, siteId))
	if err != nil {
		return err
	}
	if len(showAnchorIDs) == 0 {
		showAnchorIDs, err = core.HKeys(fmt.Sprintf(redisKey.AnchorShowDefault, siteId))
	}

	if err != nil {
		return err
	}

	// 兼容播控数据为空 且 没有展示中的场次
	if len(showAnchorIDs) == 0 {
		err = core.SetNotExpireKV(fmt.Sprintf(redisKey.AnchorShowList, siteId), "{}")
		if err != nil {
			return err
		}
		return nil
	}
	key := fmt.Sprintf(redisKey.AnchorEventsList, siteId)
	anchorListRaw, err := core.HMGet(key, showAnchorIDs...)
	if err != nil || anchorListRaw == nil {
		// 重试一次
		time.Sleep(300 * time.Millisecond)
		anchorListRaw, err = core.HMGet(key, showAnchorIDs...)
		if err != nil {
			return err
		}
	}

	for _, anchorDatum := range anchorListRaw {
		if anchorDatum == nil {
			continue
		}
		err := mdata.Cjson.Unmarshal([]byte(anchorDatum.(string)), &anchorEvent)
		if err != nil {
			continue
		}
		if anchorEvent.IsStarted == "" {
			resp["not_started"] = append(resp["not_started"], anchorEvent)
		} else {
			resp["started"] = append(resp["started"], anchorEvent)
		}
		anchorEvent.AnchorVideo = nil
	}

	// 尝试从 Redis 或 Sql 中获取冠名场馆的清单
	venuesSortMap, err := fetchVenuesSortScore(siteId)
	if err != nil {
		return err
	}

	// 处理排序规则：以时间正序排序（Timestamp越小越靠前），若时间一致则优先冠名场馆在前，其余场馆以后台 sort 分数进行排序，小的优先。
	// 冠名场馆定义：冠上站点名称的场馆，如：爱游戏体育 就属于爱游戏的冠名场馆
	// 判定方式： 每个场馆类型中序号最小的那个。
	if len(resp["not_started"]) > 1 {
		sortAnchorEventList(resp["not_started"], venuesSortMap)
	}
	if len(resp["started"]) > 1 {
		sortAnchorEventList(resp["started"], venuesSortMap)
	}

	marshal, err := mdata.Cjson.MarshalToString(resp)
	if err != nil {
		return err
	}
	if len(resp["not_started"])+len(resp["started"]) > 0 {
		err = core.SetNotExpireKV(fmt.Sprintf(redisKey.AnchorShowList, siteId), marshal)
		if err != nil {
			return err
		}
	}
	return nil
}

// 取得场馆排序分数 Map
func fetchVenuesSortScore(siteId int) (venueScoreMap map[string]int, err error) {
	venuesSortMapKey := fmt.Sprintf(redisKey.VenueSortList, siteId)

	// 尝试从 redis 中捞取
	venuesBytes, err := core.GetKeyBytes(venuesSortMapKey)
	if err == nil && len(venuesBytes) > 0 {
		err = mdata.Cjson.Unmarshal(venuesBytes, &venueScoreMap)

		if err != nil {
			return nil, err
		}
		// 取得并解码成功，直接返回。
		return venueScoreMap, err
	}

	// Redis 中无资料，尝试从 Database 中读取并存入快取。
	venueScoreMap, err = sqldb.QueryVenuesSortScore()
	if err != nil {
		return nil, err
	}

	// 取得冠名场馆资料
	specificVenues, err := sqldb.QuerySpecificVenues()
	if err != nil {
		return nil, err
	}

	// 合并两者将冠名场馆分数设定为负值，优先级最高。
	for venueName := range specificVenues {
		venueScoreMap[venueName] = -1000
	}

	// 转换成 string 存入 redis
	cacheStr, err := mdata.Cjson.MarshalToString(venueScoreMap)
	if err != nil {
		return nil, err
	}
	// 当前过期条件：5分钟
	err = core.SetExpireKV(venuesSortMapKey, cacheStr, time.Minute*5)

	if err != nil {
		return nil, err
	}

	return venueScoreMap, nil
}

// 依据传入的排序分数 Map 对场馆依照时间正序及分数进行排序。
func sortAnchorEventList(list video.AnchorList, venuesScoreMap map[string]int) {
	// sort.Slice() 可直接用于空 slice 故不需要额外判断。
	sort.Slice(list, func(i, j int) bool {
		if list[i].TimeStamp < list[j].TimeStamp {
			return true
		} else if list[i].TimeStamp > list[j].TimeStamp {
			return false
		}
		// 越小优先级越高
		return getVenueScore(list[i].VenueName, venuesScoreMap) < getVenueScore(list[j].VenueName, venuesScoreMap)
	})
}

// 从传入的 Map 中取得场馆分数，若不存在则返回不存在的分数。
func getVenueScore(venueName string, venuesMap map[string]int) int {
	venueTrueName := getVenueTrueName(venueName)

	if val, ok := venuesMap[venueTrueName]; ok {
		return val
	}
	return 9999
}

// 转换场馆正式名称。
func getVenueTrueName(venueName string) string {
	upperVenueName := strings.ToUpper(venueName)

	switch upperVenueName {
	case "IM":
		return "IMTY"
	case "WM":
		return "WMTY"
	case "FB":
		return "FBTY"
	case "V188":
		return "XJTY"
	default:
		return upperVenueName
	}
}

func UpdateAnchorEventList(c *context.Context) {
	var (
		anchorListMap      = make(map[string]string)
		anchorListInterMap = make(map[string]interface{})
		expiredKey         []string
	)
	anchorRedisKey := fmt.Sprintf(redisKey.AnchorEventsList, metadata.GetSiteId(c))
	anchorListMap, err := core.HGetAlL(anchorRedisKey)
	if err != nil {
		c.Errorf("get anchorRedisKey fail")
		return
	}
	for key, anchorDatum := range anchorListMap {
		if anchorDatum == "" {
			expiredKey = append(expiredKey, key)
			continue
		}
		var anchorEvent video.AnchorEvents
		err := mdata.Cjson.Unmarshal([]byte(anchorDatum), &anchorEvent)
		if err != nil {
			return
		}
		if anchorEvent.Expired() {
			expiredKey = append(expiredKey, key)
			continue
		}
		downTime, err1 := utils.CountDownTime(anchorEvent.StartDate)
		if err1 != nil {
			c.Errorf("SetAnchorLive CountDownTime err = %v", err)
			continue
		}
		if downTime == "0" {
			anchorEvent.IsStarted = "1" // 直播是否已开始
		}
		anchorListInterMap[key] = anchorEvent
	}
	if len(expiredKey) > 0 {
		c.Infof("delete expired anchor key:%d", len(expiredKey))
		DeleteAnchorLive(c, expiredKey...)
	}
	err = core.HMSet(anchorRedisKey, anchorListInterMap)
	if err != nil {
		c.Errorf("UpdateAnchorEventList error:%s", err)
	}
	err = updateShowAnchorData(metadata.GetSiteId(c))
	if err != nil {
		c.Errorf("更新主播接口数据失败:%s", err)
	}
	return
}

func DeleteAnchorLive(c *context.Context, key ...string) {
	err := core.HDel(fmt.Sprintf(redisKey.AnchorEventsList, metadata.GetSiteId(c)), key...)
	if err != nil {
		c.Infof("del ANCHOR_EVENTS_LIST:%s error:%s", key, err)
	}
	err = core.HDel(fmt.Sprintf(redisKey.AnchorShowDefault, metadata.GetSiteId(c)), key...)
	if err != nil {
		c.Infof("del ANCHOR_SHOW_DEFAULT:%s error:%s", key, err)
	}
	err = core.HDel(fmt.Sprintf(redisKey.AnchorShowIds, metadata.GetSiteId(c)), key...)
	if err != nil {
		c.Infof("del ANCHOR_SHOW_IDS:%s error:%s", key, err)
	}
	err = core.HDel(fmt.Sprintf(redisKey.AnchorShowUpdate, metadata.GetSiteId(c)), key...)
	if err != nil {
		c.Infof("del ANCHOR_SHOW_UPDATE:%s error:%s", key, err)
	}
}

func GetActivityContestThemeVideo(c *context.Context) {
	var (
		activityContestThemeLiveList []model.ActivityContestThemeLive
		err                          error
		liveData                     video.BkResponseForSourceUrl
	)
	value, err := core.GetKey(fmt.Sprintf("activity_contest_theme_video_sid_%d", metadata.GetSiteId(c)))
	c.Infof("GetActivityContestThemeVideo value=%v", value)
	if err != nil {
		c.Errorf("GetActivityContestThemeVideo GetKey err = %v", err)
		return
	}
	if value == "" {
		c.Errorf("GetActivityContestThemeVideo redis value is empty")
		return
	}
	err = mdata.Cjson.Unmarshal([]byte(value), &activityContestThemeLiveList)
	if err != nil {
		c.Errorf("GetActivityContestThemeVideo Unmarshal err = %v", err)
		return
	}
	if len(activityContestThemeLiveList) == 0 {
		c.Errorf("GetActivityContestThemeVideo activityContestThemeLiveList 长度为0")
		return
	}
	for _, live := range activityContestThemeLiveList {
		jsonChar, err := getTheVenueMatch(c, live.VenueCode, live.ContestId)
		if err != nil {
			c.Errorf("GetActivityContestThemeVideo | getTheVenueMatchList err = %v,场馆code=%v,比赛id=%v", err, live.VenueCode, live.ContestId)
			continue
		}
		err = mdata.Cjson.Unmarshal(jsonChar, &liveData)
		if err != nil {
			c.Errorf("GetActivityContestThemeVideo getBkLiveData Unmarshal err = %v,场馆code=%v,比赛id=%v = %v", err, live.VenueCode, live.ContestId)
			continue
		}
		if liveData.StatusCode != 200 {
			c.Errorf("GetActivityContestThemeVideo getBkLiveData 播控返回状态码 StatusCode= %v,Message = %v,场馆code=%v,比赛id=%v = %v", liveData.StatusCode, liveData.Message, live.VenueCode, live.ContestId)
			continue
		}
		key := "sport_event:" + live.VenueCode
		fieldKey := live.VenueCode + "_" + live.ContestId
		oddsJson, err := core.HGet(key, fieldKey)
		if err != nil || oddsJson == "" {
			c.Warnf("GetActivityContestThemeVideo err = %v,oddsJson=%v,fieldKey=%v", err, oddsJson, fieldKey)
			continue
		}
		err = mdata.Cjson.Unmarshal([]byte(oddsJson), &liveData.Data.PlateVo)
		if err != nil {
			c.Errorf("GetActivityContestThemeVideo Unmarshal err = %v,liveData.PlateVo=%+v", err, liveData.Data.PlateVo)
			continue
		}
		liveKey := "VIDEO_ACTIVITY_CONTEST_THEME"
		marshal, err := mdata.Cjson.Marshal(liveData.Data)
		if err != nil {
			c.Errorf("GetActivityContestThemeVideo Marshal err = %v, liveData = %+v", err, liveData)
			continue
		}
		err = core.HSet(liveKey, fieldKey, marshal)
		if err != nil {
			c.Errorf("GetActivityContestThemeVideo HSet err = %v, liveData = %+v", err, liveData)
			continue
		}
		c.Infof("GetActivityContestThemeVideo HSet success,fieldKey = %v", fieldKey)
	}
	return
}

// getTheVenueMatch 获取指定比赛播控视频信息
func getTheVenueMatch(c *context.Context, liveSiteAlias, eid string) (resp []byte, err error) {
	var (
		LiveSourceUrl = config.GetConfig().VideoLiveSourceUrl
		alias         string
		api           string
	)
	if liveSiteAlias == "IMTY" {
		alias = "im"
	} else if liveSiteAlias == "WMTY" {
		alias = "wm"
	} else if liveSiteAlias == "FBTY" {
		alias = "fb"
	} else {
		alias = ""
	}
	if alias == "" {
		c.Errorf("getTheVenueMatchList 未知场馆 liveSiteAlias = %s", liveSiteAlias)
		err = errors.New("未知场馆:" + liveSiteAlias)
		return
	}
	api = alias + "_source_url_list.txt?eid=" + eid
	path := LiveSourceUrl + "/video/v1/" + config.GetConfig().VideoLiveSiteAlias + "/" + api
	resp, err = httpclient.ProxyGet(path, map[string]string{mdata.HeaderSite: c.SiteId}, httpclient.GetVideoProxyClient(time.Minute*5))
	if err != nil {
		c.Errorf("getTheVenueMatch err = %v | url = %v", err, path)
		return nil, err
	}
	return
}
