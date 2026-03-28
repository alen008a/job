package service

import (
	"fmt"
	"gorm.io/gorm"
	"siteVideoJob/lib/httpclient"
	"siteVideoJob/utils"
	"strconv"
	"strings"

	"siteVideoJob/config"
	"siteVideoJob/db/sqldb"
	"siteVideoJob/internal/context"
	"siteVideoJob/internal/glog"
	"siteVideoJob/mdata"
	"siteVideoJob/mdata/model"
	"siteVideoJob/xxl"
	"time"
)

const (
	getLanguageList    = "/ls/language/%s/list"    //获取各类型数据多语言版本
	getCategoryList    = "/ls/category/%s/list"    //获取分类列表
	getCountryList     = "/ls/country/%s/list"     //获取国家列表
	getCompetitionList = "/ls/competition/%s/list" //获取联赛列表
	getTeamList        = "/ls/team/%s/list"        //获队伍列表
	getMatchList       = "/ls/match/%s/list"       //获队伍列表
	getVideoUrlList    = "/ls/video/list"          //获队直播源
)

func PullLSTYCountryList(c *context.Context) string {
	var msgBuild strings.Builder
	matchClass := config.GetLSTYMatchClassList()
	for i := range matchClass {
		msg := DoSetLSTYCountry(c, matchClass[i])
		msgBuild.WriteString(msg)
		msgBuild.WriteString("\r\n")
	}
	return msgBuild.String()
}

func PullLSTYCategoryList(c *context.Context) string {
	var msgBuild strings.Builder
	matchClass := config.GetLSTYMatchClassList()
	for i := range matchClass {
		msg := DoSetLSTYCategory(c, matchClass[i])
		msgBuild.WriteString(msg)
		msgBuild.WriteString("\r\n")
	}
	return msgBuild.String()
}

func PullLSTYTeamList(c *context.Context, executorParams *xxl.ExecutorParams) string {
	var msgBuild strings.Builder
	matchClass := config.GetLSTYMatchClassList()
	for i := range matchClass {
		msg := DoSetLSTYTeam(c, executorParams, matchClass[i])
		msgBuild.WriteString(msg)
		msgBuild.WriteString("\r\n")
	}
	return msgBuild.String()
}

func PullLSTYCompetitionList(c *context.Context, executorParams *xxl.ExecutorParams) string {
	var msgBuild strings.Builder
	matchClass := []mdata.GameType{mdata.FT, mdata.BK}
	for i := range matchClass {
		msg := DoSetLSTYCompetition(c, executorParams, string(matchClass[i]))
		msgBuild.WriteString(msg)
		msgBuild.WriteString("\r\n")
	}
	return msgBuild.String()
}

func PullLSTYMatchList(c *context.Context, executorParams *xxl.ExecutorParams) string {
	var msgBuild strings.Builder
	matchClass := config.GetLSTYMatchClassList()
	for i := range matchClass {
		msg := DoSetLSTYMatch(c, executorParams, matchClass[i])
		msgBuild.WriteString(msg)
		msgBuild.WriteString("\r\n")
	}
	return msgBuild.String()
}

func PullLSTYVideoUrlList(c *context.Context) string {
	return DoSetLSTYVideoUrl(c)
}

func PullLSTYLanguageList(c *context.Context, executorParams *xxl.ExecutorParams) string {
	var msgBuild strings.Builder
	matchClass := config.GetLSTYMatchClassList()
	dataTypes := []model.DataType{model.CategoryDataType, model.CountryDataType, model.CompetitionDataType, model.TeamDataType}
	for i := range matchClass {
		for j := range dataTypes {
			msg := DoSetLSTYLanguage(c, executorParams, matchClass[i], int(dataTypes[j]))
			msgBuild.WriteString(msg)
			msgBuild.WriteString("\r\n")
		}
	}
	return msgBuild.String()
}

// DoSetLSTYCountry 落地国家数据
func DoSetLSTYCountry(c *context.Context, matchClass string) (msg string) {
	var (
		countryResp       model.CountryResp
		matchCountryDatas []*model.MatchCountryData
		dataIds           = make([]string, 0)
	)
	data, err := GetLSTYDataList(c, fmt.Sprintf(getCountryList, matchClass), nil)
	if err != nil {
		msg = fmt.Sprintf("DoSetLSTYCountry matchClass=%s error:%+v", matchClass, err)
		c.Error(msg)
		return
	}
	if len(data) == 0 {
		return
	}
	err = mdata.Cjson.Unmarshal(data, &countryResp)
	if err != nil {
		msg = fmt.Sprintf("DoSetLSTYCountry matchClass=%s error:%+v", matchClass, err)
		c.Error(msg)
		return
	}
	if countryResp.StatusCode != 200 {
		msg = fmt.Sprintf("DoSetLSTYCountry GetLSTYDataList err matchClass=%s | countryResp.StatusCode = %v | countryResp.Status = %v", matchClass, countryResp.StatusCode, countryResp.Status)
		c.Error(msg)
		return
	}

	for _, v := range countryResp.Data {
		dataIds = append(dataIds, v.ID)
	}

	nameMap := getLanguageMap(c, dataIds, matchClass, model.CountryDataType)
	for _, v := range countryResp.Data {
		matchCountryData := &model.MatchCountryData{
			CategoryID:       v.CategoryID,
			CountryID:        v.ID,
			MatchClass:       matchClass,
			Name:             v.Name,
			Logo:             v.Logo,
			UpdatedTimestamp: v.UpdatedAt,
			CreatedBy:        mdata.System,
			UpdatedBy:        mdata.System,
		}
		if m, ok := nameMap[v.ID]; ok {
			matchCountryData.NameZh = m.NameZh
			matchCountryData.NameZht = m.NameZht
		}
		matchCountryDatas = append(matchCountryDatas, matchCountryData)
	}

	err = sqldb.Video().Transaction(func(tx *gorm.DB) error {
		affectRows, err := sqldb.DeleteLSCountry(tx, matchClass)
		if err != nil {
			c.Errorf("DoSetLSTYCountry DeleteLSCountry error:%+v matchClass:%s", err, matchClass)
			return err
		}
		c.Infof("DoSetLSTYCountry DeleteLSCountry affectRows:%d matchClass:%s", affectRows, matchClass)
		affectRows, err = sqldb.BatchInsertLSCountry(tx, matchCountryDatas)
		if err != nil {
			c.Errorf("DoSetLSTYCountry BatchInsertLSCountry error:%+v matchClass:%s", err, matchClass)
			return err
		}
		c.Infof("DoSetLSTYCountry BatchInsertLSCountry affectRows:%d matchClass:%s", affectRows, matchClass)
		return nil
	})
	if err != nil {
		msg = fmt.Sprintf("DoSetLSTYCountry error:%+v matchClass:%s", err, matchClass)
		c.Error(msg)
	} else {
		msg = fmt.Sprintf("DoSetLSTYCountry success num:%d matchClass:%s", len(matchCountryDatas), matchClass)
		c.Info(msg)
	}
	return
}

// DoSetLSTYCategory 落地类别数据
func DoSetLSTYCategory(c *context.Context, matchClass string) (msg string) {
	var (
		categoryResp       model.CategoryResp
		matchCategoryDatas []*model.MatchCategoryData
		dataIds            = make([]string, 0)
	)
	data, err := GetLSTYDataList(c, fmt.Sprintf(getCategoryList, matchClass), nil)
	if err != nil {
		msg = fmt.Sprintf("DoSetLSTYCategory matchClass:%s, error:%+v", matchClass, err)
		c.Error(msg)
		return c.GetConsoleLog()
	}
	if len(data) == 0 {
		return
	}
	err = mdata.Cjson.Unmarshal(data, &categoryResp)
	if err != nil {
		msg = fmt.Sprintf("DoSetLSTYCategory matchClass:%s, error:%+v", matchClass, err)
		c.Error(msg)
		return
	}
	if categoryResp.StatusCode != 200 {
		msg = fmt.Sprintf("DoSetLSTYCategory GetLSTYDataList err matchClass=%s | categoryResp.StatusCode = %v | categoryResp.Status = %v", matchClass, categoryResp.StatusCode, categoryResp.Status)
		c.Error(msg)
		return
	}

	for _, v := range categoryResp.Data {
		dataIds = append(dataIds, v.Id)
	}

	nameMap := getLanguageMap(c, dataIds, matchClass, model.CategoryDataType)
	for _, v := range categoryResp.Data {
		matchCategoryData := &model.MatchCategoryData{
			CategoryID:       v.Id,
			MatchClass:       matchClass,
			Name:             v.Name,
			UpdatedTimestamp: v.UpdatedAt,
			CreatedBy:        mdata.System,
			UpdatedBy:        mdata.System,
		}
		if m, ok := nameMap[v.Id]; ok {
			matchCategoryData.NameZh = m.NameZh
			matchCategoryData.NameZht = m.NameZht
		}
		matchCategoryDatas = append(matchCategoryDatas, matchCategoryData)
	}

	err = sqldb.Video().Transaction(func(tx *gorm.DB) error {
		affectRows, err := sqldb.DeleteLSCategory(tx, matchClass)
		if err != nil {
			c.Errorf("DoSetLSTYCategory DeleteLSCategory error:%+v matchClass:%s", err, matchClass)
			return err
		}
		c.Infof("DoSetLSTYCategory DeleteLSCategory affectRows:%d matchClass:%s", affectRows, matchClass)
		affectRows, err = sqldb.BatchInsertLSCategory(tx, matchCategoryDatas)
		if err != nil {
			c.Errorf("DoSetLSTYCategory BatchInsertLSCategory error:%+v matchClass:%s", err, matchClass)
			return err
		}
		c.Infof("DoSetLSTYCategory BatchInsertLSCategory affectRows:%d matchClass:%s", affectRows, matchClass)
		return nil
	})
	if err != nil {
		msg = fmt.Sprintf("DoSetLSTYCategory error:%+v matchClass:%s", err, matchClass)
		c.Error(msg)
	} else {
		msg = fmt.Sprintf("DoSetLSTYCategory success num:%d matchClass:%s", len(matchCategoryDatas), matchClass)
		c.Info(msg)
	}
	return
}

// DoSetLSTYTeam 落地队伍数据
func DoSetLSTYTeam(c *context.Context, executorParams *xxl.ExecutorParams, matchClass string) (msg string) {
	var (
		total, endDate int
		pullType       = executorParams.PullType
		pullReqType    = executorParams.PullReqType
		startTime      = executorParams.StartTime
		commonReq      = model.CommonReq{}
	)
	if executorParams.EndTime != "" {
		endTime, err := utils.ParseTime(startTime)
		if err == nil {
			endDate = int(endTime.Unix())
		}
	}
	//全量拉取
	if pullType == mdata.PullFull {
		if pullReqType == mdata.PagePullReqType {
			commonReq.Page = 1
		}
		if pullReqType == mdata.TimePullReqType {
			if startTime != "" {
				startDate, err := utils.ParseTime(startTime)
				if err != nil {
					c.ConsoleErr("DoSetLSTYTeam ParseTime error:%+v ,param:%+v", err, executorParams)
					return c.GetConsoleLog()
				}
				commonReq.Time = int(startDate.Unix())
			} else {
				commonReq.Time = int(time.Now().Add(-24 * time.Hour).Unix())
			}
		}
	}
	//增量拉取
	if pullType == mdata.PullIncr {
		pullReqType = mdata.TimePullReqType
		if startTime != "" {
			startDate, err := utils.ParseTime(startTime)
			if err != nil {
				c.ConsoleErr("DoSetLSTYTeam ParseTime error:%+v ,param:%+v", err, executorParams)
				return c.GetConsoleLog()
			}
			commonReq.Time = int(startDate.Unix())
		} else {
			//查找最近的更新时间
			updatedTime := sqldb.GetLastTeamUpdatedTime(matchClass)
			if updatedTime == 0 {
				updatedTime = int(time.Now().Add(-24 * time.Hour).Unix())
			}
			commonReq.Time = updatedTime
		}
	}

	for {
		teamResp, err := getLsTeamList(c, commonReq, matchClass)
		if err != nil {
			msg = fmt.Sprintf("DoSetLSTYTeam getLsTeamList matchClass:%s, error:%+v ,param:%+v", matchClass, err, executorParams)
			c.Error(msg)
			return
		}
		if teamResp.StatusCode != 200 {
			msg = fmt.Sprintf("DoSetLSTYTeam GetLSTYDataList err  matchClass:%s, teamResp.StatusCode = %v | teamResp.Status = %v", matchClass, teamResp.StatusCode, teamResp.Status)
			c.Error(msg)
			return
		}
		dataIds := make([]string, 0)
		for _, v := range teamResp.Data.Results {
			dataIds = append(dataIds, v.ID)
		}

		nameMap := getLanguageMap(c, dataIds, matchClass, model.TeamDataType)
		for _, v := range teamResp.Data.Results {
			matchTeamData := &model.MatchTeamData{
				TeamID:           v.ID,
				UID:              v.UID,
				CountryID:        v.CountryID,
				Name:             v.Name,
				MatchClass:       matchClass,
				ShortName:        v.ShortName,
				UpdatedTimestamp: v.UpdatedAt,
				Logo:             v.Logo,
				National:         v.National,
				CountryLogo:      v.CountryLogo,
				CreatedBy:        mdata.System,
				UpdatedBy:        mdata.System,
			}
			if m, ok := nameMap[v.ID]; ok {
				matchTeamData.NameZh = m.NameZh
				matchTeamData.NameZht = m.NameZht
			}
			affectRow, err := sqldb.UpsertMatchTeamData(matchTeamData)
			if err != nil {
				c.Errorf("DoSetLSTYTeam UpsertMatchTeamData error:%+v ,param:%+v", err, matchTeamData)
			} else {
				c.Infof("DoSetLSTYTeam UpsertMatchTeamData success affectRow:%+v ,param:%+v", affectRow, matchTeamData)
			}
		}
		if teamResp.Data.Query.Total > 0 {
			total += teamResp.Data.Query.Total
			if pullReqType == mdata.PagePullReqType {
				commonReq.Page += 1
			}
			if pullReqType == mdata.TimePullReqType {
				commonReq.Time = teamResp.Data.Query.MaxTime
				if endDate >= teamResp.Data.Query.MaxTime || teamResp.Data.Query.MaxTime == teamResp.Data.Query.MinTime {
					msg = fmt.Sprintf("DoSetLSTYTeam UpsertMatchTeamData matchClass:%s, success total:%+v ", matchClass, total)
					c.Info(msg)
					break
				}
			}
		} else {
			msg = fmt.Sprintf("DoSetLSTYTeam UpsertMatchTeamData matchClass:%s, success total:%+v ", matchClass, total)
			c.Info(msg)
			break
		}
	}
	return
}

// DoSetLSTYCompetition 落地联赛数据
func DoSetLSTYCompetition(c *context.Context, executorParams *xxl.ExecutorParams, matchClass string) (msg string) {
	var (
		total, endDate int
		pullType       = executorParams.PullType
		pullReqType    = executorParams.PullReqType
		startTime      = executorParams.StartTime
		commonReq      = model.CommonReq{}
	)

	if executorParams.EndTime != "" {
		endTime, err := utils.ParseTime(startTime)
		if err == nil {
			endDate = int(endTime.Unix())
		}
	}
	//全量拉取
	if pullType == mdata.PullFull {
		if pullReqType == mdata.PagePullReqType {
			commonReq.Page = 1
		}
		if pullReqType == mdata.TimePullReqType {
			if startTime != "" {
				startDate, err := utils.ParseTime(startTime)
				if err != nil {
					c.ConsoleErr("DoSetLSTYTeam ParseTime error:%+v ,param:%+v", err, executorParams)
					return c.GetConsoleLog()
				}
				commonReq.Time = int(startDate.Unix())
			} else {
				commonReq.Time = int(time.Now().Add(-24 * time.Hour).Unix())
			}
		}
	}

	//增量拉取
	if pullType == mdata.PullIncr {
		pullReqType = mdata.TimePullReqType
		if startTime != "" {
			startDate, err := utils.ParseTime(startTime)
			if err != nil {
				c.ConsoleErr("DoSetLSTYCompetition ParseTime matchClass=%s error:%+v ,param:%+v", matchClass, err, executorParams)
				return c.GetConsoleLog()
			}
			commonReq.Time = int(startDate.Unix())
		} else {
			//查找最近的更新时间
			updatedTime := sqldb.GetLastCompetitionUpdatedTime(matchClass)
			if updatedTime == 0 {
				updatedTime = int(time.Now().Add(-24 * time.Hour).Unix())
			}
			commonReq.Time = updatedTime
		}
	}

	for {
		competitionResp, err := getLsCompetitionList(c, commonReq, matchClass)
		if err != nil {
			msg = fmt.Sprintf("DoSetLSTYCompetition getLsCompetitionList matchClass=%s error:%+v,param:%+v", matchClass, err, executorParams)
			c.Error(msg)
			return
		}

		if competitionResp.StatusCode != 200 {
			msg = fmt.Sprintf("DoSetLSTYCompetition GetLSTYDataList err matchClass=%s | competitionResp.StatusCode = %v | competitionResp.Status = %v", matchClass, competitionResp.StatusCode, competitionResp.Status)
			c.Error(msg)
			return
		}
		dataIds := make([]string, 0)
		for _, v := range competitionResp.Data.Results {
			dataIds = append(dataIds, v.ID)
		}

		nameMap := getLanguageMap(c, dataIds, matchClass, model.CompetitionDataType)
		for _, v := range competitionResp.Data.Results {
			matchCompetitionData := &model.MatchCompetitionData{
				CompetitionID:    v.ID,
				CategoryID:       v.CategoryID,
				MatchClass:       matchClass,
				CountryID:        v.CountryID,
				Name:             v.Name,
				ShortName:        v.ShortName,
				Logo:             v.Logo,
				Type:             v.Type,
				CurSeasonID:      v.CurSeasonID,
				CurStageID:       v.CurStageID,
				CurRound:         v.CurRound,
				RoundCount:       v.RoundCount,
				UpdatedTimestamp: v.UpdatedAt,
				CreatedBy:        mdata.System,
				UpdatedBy:        mdata.System,
			}
			if m, ok := nameMap[v.ID]; ok {
				matchCompetitionData.NameZh = m.NameZh
				matchCompetitionData.NameZht = m.NameZht
			}
			affectRow, err := sqldb.UpsertMatchCompetitionData(matchCompetitionData)
			if err != nil {
				c.Errorf("DoSetLSTYCompetition UpsertMatchCompetitionData error:%+v ,param:%+v", err, matchCompetitionData)
			} else {
				c.Infof("DoSetLSTYCompetition UpsertMatchCompetitionData success affectRow:%+v ,param:%+v", affectRow, matchCompetitionData)
			}
		}
		if competitionResp.Data.Query.Total > 0 {
			total += competitionResp.Data.Query.Total
			if pullReqType == mdata.PagePullReqType {
				commonReq.Page += 1
			}
			if pullReqType == mdata.TimePullReqType {
				commonReq.Time = competitionResp.Data.Query.MaxTime
				if endDate >= competitionResp.Data.Query.MaxTime || competitionResp.Data.Query.MaxTime == competitionResp.Data.Query.MinTime {
					msg = fmt.Sprintf("DoSetLSTYCompetition UpsertMatchCompetitionData matchClass:%s, success total:%+v ", matchClass, total)
					c.Info(msg)
					break
				}
			}
		} else {
			msg = fmt.Sprintf("DoSetLSTYCompetition UpsertMatchCompetitionData matchClass:%s, success total:%+v ", matchClass, total)
			c.Info(msg)
			break
		}
	}
	return
}

// DoSetLSTYMatch 落地比赛数据
func DoSetLSTYMatch(c *context.Context, executorParams *xxl.ExecutorParams, matchClass string) (msg string) {
	var (
		total, endDate int
		pullType       = executorParams.PullType
		pullReqType    = executorParams.PullReqType
		startTime      = executorParams.StartTime
		commonReq      = model.CommonReq{}
	)

	if executorParams.EndTime != "" {
		endTime, err := utils.ParseTime(startTime)
		if err == nil {
			endDate = int(endTime.Unix())
		}
	}
	//全量拉取
	if pullType == mdata.PullFull {
		if pullReqType == mdata.PagePullReqType {
			commonReq.Page = 1
		}
		if pullReqType == mdata.TimePullReqType {
			if startTime != "" {
				startDate, err := utils.ParseTime(startTime)
				if err != nil {
					c.ConsoleErr("DoSetLSTYTeam ParseTime error:%+v ,param:%+v", err, executorParams)
					return c.GetConsoleLog()
				}
				commonReq.Time = int(startDate.Unix())
			} else {
				commonReq.Time = int(time.Now().Add(-24 * time.Hour).Unix())
			}
		}
	}

	//增量拉取
	if pullType == mdata.PullIncr {
		pullReqType = mdata.TimePullReqType
		if startTime != "" {
			startDate, err := utils.ParseTime(startTime)
			if err != nil {
				c.ConsoleErr("DoSetLSTYMatch ParseTime matchClass:%s error:%+v ,param:%+v", matchClass, err, executorParams)
				return c.GetConsoleLog()
			}
			commonReq.Time = int(startDate.Unix())
		} else {
			//查找最近的更新时间
			updatedTime := sqldb.GetLastMatchUpdatedTime(matchClass)
			if updatedTime == 0 {
				updatedTime = int(time.Now().Add(-24 * time.Hour).Unix())
			}
			commonReq.Time = updatedTime
		}
	}

	for {
		matchResp, err := getLsMatchList(c, commonReq, matchClass)
		if err != nil {
			msg = fmt.Sprintf("DoSetLSTYMatch getLsMatchList matchClass:%s error:%+v ,param:%+v", matchClass, err, executorParams)
			c.Error(msg)
			return
		}

		if matchResp.StatusCode != 200 {
			msg = fmt.Sprintf("DoSetLSTYMatch GetLSTYDataList err matchClass:%s | matchResp.StatusCode = %v | matchResp.Status = %v", matchClass, matchResp.StatusCode, matchResp.Status)
			c.Error(msg)
			return
		}
		competitionIds, teamIds := make([]string, 0), make([]string, 0)
		data := matchResp.Data.Results
		for _, v := range data {
			if v.CompetitionID != "" {
				competitionIds = append(competitionIds, v.CompetitionID)
			}
			if v.HomeTeamID != "" {
				teamIds = append(teamIds, v.HomeTeamID)
			}
			if v.AwayTeamID != "" {
				teamIds = append(teamIds, v.AwayTeamID)
			}
		}
		competitionNameMap := getLanguageMap(c, competitionIds, matchClass, model.CompetitionDataType)
		teamNameMap := getLanguageMap(c, teamIds, matchClass, model.TeamDataType)
		cLogoMap := getCompetitionMap(c, competitionIds)
		tLogoMap := getTeamMap(c, teamIds)
		for _, v := range data {
			startAt := time.Unix(int64(v.MatchTime), 0).In(utils.GetLoctionBJ()).Format(utils.TimeBarFormat)
			matchOriginData := &model.MatchOriginData{
				MatchID:          v.ID,
				MatchClass:       matchClass,
				StatusCode:       v.StatusID,
				HomeId:           v.HomeTeamID,
				HomeLogo:         tLogoMap[v.HomeTeamID],
				VisitId:          v.AwayTeamID,
				VisitLogo:        tLogoMap[v.AwayTeamID],
				VenueName:        mdata.LSTY,
				LeagueID:         v.CompetitionID,
				LeagueLogo:       cLogoMap[v.CompetitionID],
				LiveStatus:       v.LiveStatus,
				MLive:            v.Coverage.MLive,
				Lineup:           v.Coverage.Lineup,
				StageCode:        v.Round.StageID,
				GroupNum:         v.Round.GroupNum,
				RoundNum:         v.Round.RoundNum,
				UpdatedTimestamp: v.UpdatedAt,
				StartAt:          startAt,
				CreatedBy:        mdata.System,
				UpdatedBy:        mdata.System,
			}
			if m, ok := competitionNameMap[v.CompetitionID]; ok {
				matchOriginData.MatchName = m.NameZh
			}
			if m, ok := teamNameMap[v.HomeTeamID]; ok {
				matchOriginData.HomeName = m.NameZh
			}
			if m, ok := teamNameMap[v.AwayTeamID]; ok {
				matchOriginData.VisitName = m.NameZh
			}
			affectRow, err := sqldb.UpsertMatchOriginData(matchOriginData)
			if err != nil {
				c.Errorf("DoSetLSTYMatch UpsertMatchOriginData error:%+v ,param:%+v", err, matchOriginData)
			} else {
				c.Infof("DoSetLSTYMatch UpsertMatchOriginData success affectRow:%+v ,param:%+v", affectRow, matchOriginData)
			}
		}
		if matchResp.Data.Query.Total > 0 {
			total += matchResp.Data.Query.Total
			if pullReqType == mdata.PagePullReqType {
				commonReq.Page += 1
			}
			if pullReqType == mdata.TimePullReqType {
				commonReq.Time = matchResp.Data.Query.MaxTime
				if endDate >= matchResp.Data.Query.MaxTime || matchResp.Data.Query.MaxTime == matchResp.Data.Query.MinTime {
					msg = fmt.Sprintf("DoSetLSTYMatch UpsertMatchOriginData matchClass:%s, success total:%+v ", matchClass, total)
					c.Info(msg)
					break
				}
			}
		} else {
			msg = fmt.Sprintf("DoSetLSTYMatch UpsertMatchOriginData matchClass:%s, success total:%+v ", matchClass, total)
			c.Info(msg)
			break
		}
	}
	return
}

// DoSetLSTYLanguage 落地语言数据
func DoSetLSTYLanguage(c *context.Context, executorParams *xxl.ExecutorParams, matchClass string, dataType int) (msg string) {
	var (
		total, endDate int
		pullType       = executorParams.PullType
		startTime      = executorParams.StartTime
		pullReqType    = executorParams.PullReqType
		commonReq      = model.CommonReq{Type: dataType}
	)
	if executorParams.EndTime != "" {
		endTime, err := utils.ParseTime(startTime)
		if err == nil {
			endDate = int(endTime.Unix())
		}
	}
	//全量拉取
	if pullType == mdata.PullFull {
		if pullReqType == mdata.PagePullReqType {
			commonReq.Page = 1
		}
		if pullReqType == mdata.TimePullReqType {
			if startTime != "" {
				startDate, err := utils.ParseTime(startTime)
				if err != nil {
					c.ConsoleErr("DoSetLSTYTeam ParseTime error:%+v ,param:%+v", err, executorParams)
					return c.GetConsoleLog()
				}
				commonReq.Time = int(startDate.Unix())
			} else {
				commonReq.Time = int(time.Now().Add(-24 * time.Hour).Unix())
			}
		}
	}

	//增量拉取
	if pullType == mdata.PullIncr {
		pullReqType = mdata.TimePullReqType
		if startTime != "" {
			startDate, err := utils.ParseTime(startTime)
			if err != nil {
				c.ConsoleErr("DoSetLSTYLanguage ParseTime error:%+v ,param:%+v", err, executorParams)
				return c.GetConsoleLog()
			}
			commonReq.Time = int(startDate.Unix())
		} else {
			//查找最近的更新时间
			updatedTime := sqldb.GetLastLanguageUpdatedTime(matchClass, dataType)
			if updatedTime == 0 {
				updatedTime = int(time.Now().Add(-24 * time.Hour).Unix())
			}
			commonReq.Time = updatedTime
		}
	}

	for {
		languageResp, err := getLsLanguageList(c, commonReq, matchClass)
		if err != nil {
			msg = fmt.Sprintf("DoSetLSTYLanguage getLsLanguageList matchClass:%s, dataType:%d, error:%+v ,param:%+v", matchClass, dataType, err, executorParams)
			c.Error(msg)
			return
		}
		if languageResp.StatusCode != 200 {
			msg = fmt.Sprintf("DoSetLSTYLanguage GetLSTYDataList err matchClass=%s, dataType=%d, languageResp.StatusCode = %v | languageResp.Status = %v", matchClass, dataType, languageResp.StatusCode, languageResp.Status)
			c.Error(msg)
			return
		}

		for _, v := range languageResp.Data.Results {
			matchLanguageData := &model.MatchLanguageData{
				DataId:           v.Id,
				MatchClass:       matchClass,
				DataType:         dataType,
				NameEn:           v.NameEn,
				NameZh:           v.NameZh,
				NameZht:          v.NameZht,
				UpdatedTimestamp: v.UpdatedAt,
				CreatedBy:        mdata.System,
				UpdatedBy:        mdata.System,
			}
			affectRow, err := sqldb.UpsertMatchLanguageData(matchLanguageData)
			if err != nil {
				c.Errorf("DoSetLSTYLanguage UpsertMatchLanguageData error:%+v ,param:%+v", err, matchLanguageData)
			} else {
				c.Infof("DoSetLSTYLanguage UpsertMatchLanguageData success affectRow:%+v ,param:%+v", affectRow, matchLanguageData)
			}
			switch dataType {
			case model.CategoryDataType:
				affectRow, err = sqldb.UpdateLsCategory(&model.MatchCategoryData{
					CategoryID: v.Id,
					NameZh:     v.NameZh,
					NameZht:    v.NameZht,
				})
				if err != nil {
					c.Errorf("DoSetLSTYLanguage UpdateLsCategory error:%+v,affectRow:$+v,param:%+v", err, affectRow, matchLanguageData)
				}
			case model.CountryDataType:
				affectRow, err = sqldb.UpdateLsCountry(&model.MatchCountryData{
					CountryID: v.Id,
					NameZh:    v.NameZh,
					NameZht:   v.NameZht,
				})
				if err != nil {
					c.Errorf("DoSetLSTYLanguage UpdateLsCountry error:%+v,affectRow:$+v,param:%+v", err, affectRow, matchLanguageData)
				}
			case model.CompetitionDataType:
				affectRow, err = sqldb.UpdateLsCompetition(&model.MatchCompetitionData{
					CompetitionID: v.Id,
					NameZh:        v.NameZh,
					NameZht:       v.NameZht,
				})
				if err != nil {
					c.Errorf("DoSetLSTYLanguage UpdateLsCompetition error:%+v,affectRow:$+v,param:%+v", err, affectRow, matchLanguageData)
				}
			case model.TeamDataType:
				affectRow, err = sqldb.UpdateLsTeam(&model.MatchTeamData{
					TeamID:  v.Id,
					NameZh:  v.NameZh,
					NameZht: v.NameZht,
				})
				if err != nil {
					c.Errorf("DoSetLSTYLanguage UpdateLsTeam error:%+v,affectRow:$+v,param:%+v", err, affectRow, matchLanguageData)
				}
			}
		}
		if languageResp.Data.Query.Total > 0 {
			total += languageResp.Data.Query.Total
			if pullReqType == mdata.PagePullReqType {
				commonReq.Page += 1
			}
			if pullReqType == mdata.TimePullReqType {
				commonReq.Time = languageResp.Data.Query.MaxTime
				if endDate >= languageResp.Data.Query.MaxTime || languageResp.Data.Query.MaxTime == languageResp.Data.Query.MinTime {
					msg = fmt.Sprintf("DoSetLSTYLanguage UpsertMatchLanguageData matchClass:%s,dataType:%d,success total:%+v ", matchClass, dataType, total)
					c.Info(msg)
					break
				}
			}
		} else {
			msg = fmt.Sprintf("DoSetLSTYLanguage UpsertMatchLanguageData matchClass:%s,dataType:%d,success total:%+v ", matchClass, dataType, total)
			c.Info(msg)
			break
		}
	}
	return
}

// DoSetLSTYVideoUrl  落地视频直播源
func DoSetLSTYVideoUrl(c *context.Context) (msg string) {
	var (
		videoUrlResp     model.VideoUrlResp
		updatedTimestamp = int(time.Now().Unix())
		expireTime       = 24 //视频链接过期时间
		pullHost         = config.GetConfig().Application.LSTYMatchPullVideoHost
		pushHost         = config.GetConfig().Application.LSTYMatchPushVideoHost
		cdnKey           = config.GetConfig().Application.LiveStreamCdnKey
	)
	data, err := GetLSTYDataList(c, getVideoUrlList, nil)
	if err != nil {
		msg = fmt.Sprintf("DoSetLSTYVideoUrl  error:%+v", err)
		c.Error(msg)
		return
	}
	if len(data) == 0 {
		return
	}
	err = mdata.Cjson.Unmarshal(data, &videoUrlResp)
	if err != nil {
		msg = fmt.Sprintf("DoSetLSTYVideoUrl  error:%+v", err)
		c.Error(msg)
		return
	}
	if videoUrlResp.StatusCode != 200 {
		msg = fmt.Sprintf("DoSetLSTYVideoUrl GetLSTYDataList err  videoUrlResp.StatusCode = %v | videoUrlResp.Status = %v", videoUrlResp.StatusCode, videoUrlResp.Status)
		c.Error(msg)
		return
	}
	for _, v := range videoUrlResp.Data {
		if v.Pushurl1 == "" && v.Pushurl2 == "" {
			continue
		}
		startAt := time.Unix(int64(v.MatchTime), 0).In(utils.GetLoctionBJ()).Format(utils.TimeBarFormat)
		matchClass := config.GetGameType("LSTYCode", strconv.Itoa(v.SportId))
		matchVideoData := &model.MatchVideoData{
			MatchID:          v.MatchId,
			MatchClass:       matchClass,
			UpdatedTimestamp: updatedTimestamp,
			StartAt:          startAt,
			CreatedBy:        mdata.System,
			UpdatedBy:        mdata.System,
		}
		if v.Pushurl1 != "" {
			txSecret, txTime := getVideoUrlAuthKey(v.Pushurl1, cdnKey, expireTime)
			matchVideoData.PushUrl1 = strings.Replace(v.Pushurl1, pushHost, pullHost, 1) + ".m3u8?txSecret=" + txSecret + "&txTime=" + txTime
		}
		if v.Pushurl2 != "" {
			txSecret, txTime := getVideoUrlAuthKey(v.Pushurl2, cdnKey, expireTime)
			matchVideoData.PushUrl2 = strings.Replace(v.Pushurl2, pushHost, pullHost, 1) + ".m3u8?txSecret=" + txSecret + "&txTime=" + txTime
		}
		affectRow, err := sqldb.UpsertMatchVideoUrlData(matchVideoData)
		if err != nil {
			c.Errorf("DoSetLSTYVideoUrl UpsertMatchVideoUrlData error:%+v ,param:%+v", err, matchVideoData)
		} else {
			c.Infof("DoSetLSTYVideoUrl UpsertMatchVideoUrlData success affectRow:%+v ,param:%+v", affectRow, matchVideoData)
		}
	}
	return
}

func getVideoUrlAuthKey(path, key string, expire int) (txSecret string, txTime string) {
	streamName := path[strings.LastIndex(path, "/")+1:]
	txTime = fmt.Sprintf("%X", utils.GetBjNowTime().Add(time.Duration(expire)*time.Hour).Unix())
	txSecret = utils.Md5EncodeToString(key + streamName + txTime)
	return
}

func getCompetitionMap(c *context.Context, ids []string) map[string]string {
	logoMap := make(map[string]string)
	data, err := sqldb.QueryLSCompetitionByIds(ids)
	if err != nil {
		c.Errorf("getCompetitionMap err:%+v", err)
		return logoMap
	}
	for _, v := range data {
		logoMap[v.CompetitionID] = v.Logo
	}
	return logoMap
}

func getTeamMap(c *context.Context, ids []string) map[string]string {
	logoMap := make(map[string]string)
	data, err := sqldb.QueryLSTeamByIds(ids)
	if err != nil {
		c.Errorf("getTeamMap err:%+v", err)
		return logoMap
	}
	for _, v := range data {
		logoMap[v.TeamID] = v.Logo
	}
	return logoMap
}

func getLanguageMap(c *context.Context, ids []string, matchClass string, dataType int) map[string]*model.MatchLanguageData {
	languageMap := make(map[string]*model.MatchLanguageData)
	data, err := sqldb.QueryLSLanguageByIds(ids, matchClass, dataType)
	if err != nil {
		c.Errorf("getLanguageMap err:%+v", err)
		return languageMap
	}
	c.Infof("getLanguageMap result:%+v len:%d", data, len(data))
	for _, v := range data {
		languageMap[v.DataId] = v
	}
	return languageMap
}

func getLsTeamList(c *context.Context, commonReq model.CommonReq, matchClass string) (model.TeamResp, error) {
	var (
		teamResp model.TeamResp
	)
	data, err := GetLSTYDataList(c, fmt.Sprintf(getTeamList, matchClass), commonReq)
	if err != nil {
		return teamResp, err
	}
	err = mdata.Cjson.Unmarshal(data, &teamResp)
	if err != nil {
		return teamResp, err
	}
	return teamResp, nil
}

func getLsCompetitionList(c *context.Context, commonReq model.CommonReq, matchClass string) (model.CompetitionResp, error) {
	var (
		competitionResp model.CompetitionResp
	)
	data, err := GetLSTYDataList(c, fmt.Sprintf(getCompetitionList, matchClass), commonReq)
	if err != nil {
		return competitionResp, err
	}
	err = mdata.Cjson.Unmarshal(data, &competitionResp)
	if err != nil {
		return competitionResp, err
	}
	return competitionResp, nil
}

func getLsMatchList(c *context.Context, commonReq model.CommonReq, matchClass string) (model.MatchResp, error) {
	var (
		matchResp model.MatchResp
	)
	data, err := GetLSTYDataList(c, fmt.Sprintf(getMatchList, matchClass), commonReq)
	if err != nil {
		return matchResp, err
	}
	err = mdata.Cjson.Unmarshal(data, &matchResp)
	if err != nil {
		return matchResp, err
	}
	return matchResp, nil
}

func getLsLanguageList(c *context.Context, commonReq model.CommonReq, matchClass string) (model.LanguageResp, error) {
	var (
		languageResp model.LanguageResp
	)
	data, err := GetLSTYDataList(c, fmt.Sprintf(getLanguageList, matchClass), commonReq)
	if err != nil {
		return languageResp, err
	}
	err = mdata.Cjson.Unmarshal(data, &languageResp)
	if err != nil {
		return languageResp, err
	}
	return languageResp, nil
}

// GetLSTYDataList  获取雷速体育数据列表
func GetLSTYDataList(c *context.Context, path string, param interface{}) (resp []byte, err error) {
	var (
		LiveSourceUrl = config.GetConfig().VideoLiveSourceUrl
		LiveSiteAlias = config.GetConfig().VideoLiveSiteAlias
	)

	path = LiveSourceUrl + "/video/v1/" + LiveSiteAlias + path
	paramByte, _ := mdata.Cjson.Marshal(param)
	resp, err = httpclient.POSTJson(path, paramByte, map[string]string{mdata.HeaderSite: c.SiteId}, httpclient.GetVideoProxyClient(time.Minute*20))
	if err != nil {
		glog.EmergencyWithTimeout("GetLSTYDataList url=%s |err=%v", 60, path, err)
		return nil, err
	}
	return
}
