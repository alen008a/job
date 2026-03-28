package sqldb

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"siteVideoJob/mdata/model"
)

func BatchInsertLSCountry(tx *gorm.DB, list []*model.MatchCountryData) (int64, error) {
	if tx == nil {
		tx = Video()
	}
	db := tx.Model(&model.MatchCountryData{}).CreateInBatches(list, len(list))
	return db.RowsAffected, db.Error
}

func DeleteLSCountry(tx *gorm.DB, matchClass string) (int64, error) {
	if tx == nil {
		tx = Video()
	}
	db := tx.Table((&model.MatchCountryData{}).TableName()).Where("match_class = ?", matchClass).Delete(&model.MatchCountryData{})
	err := db.Error
	if err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func UpdateLsCountry(matchCountryData *model.MatchCountryData) (int64, error) {
	var update = make(map[string]interface{})
	if matchCountryData.NameZh != "" {
		update["name_zh"] = matchCountryData.NameZh
	}
	if matchCountryData.NameZht != "" {
		update["name_zht"] = matchCountryData.NameZht
	}
	db := Video().Model(&model.MatchCountryData{}).Where("country_id = ?", matchCountryData.CountryID).Updates(update)
	err := db.Error
	if err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func UpdateLsTeam(matchTeamData *model.MatchTeamData) (int64, error) {
	var update = make(map[string]interface{})
	if matchTeamData.NameZh != "" {
		update["name_zh"] = matchTeamData.NameZh
	}
	if matchTeamData.NameZht != "" {
		update["name_zht"] = matchTeamData.NameZht
	}
	db := Video().Model(&model.MatchTeamData{}).Where("team_id = ?", matchTeamData.TeamID).Updates(update)
	err := db.Error
	if err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func GetLastTeamUpdatedTime(matchClass string) int {
	var matchTeamData *model.MatchTeamData
	err := Video().Model(&model.MatchTeamData{}).Select(" MAX(updated_timestamp) as updated_timestamp").Where("match_class = ?", matchClass).First(&matchTeamData).Error
	if err != nil {
		return 0
	}
	return matchTeamData.UpdatedTimestamp
}

// UpsertMatchTeamData 插入/更新记录
func UpsertMatchTeamData(matchTeamData *model.MatchTeamData) (int64, error) {
	db := Video().Model(&matchTeamData).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "team_id"}, {Name: "match_class"}}, //冲突的字段：team_id 为唯一键
		DoUpdates: clause.AssignmentColumns([]string{"uid", "country_id", "name", "short_name", "name_zh", "name_zht", "updated_timestamp", "logo", "national", "country_logo"}),
	}).Create(&matchTeamData)
	return db.RowsAffected, db.Error
}

// UpsertMatchCompetitionData 插入/更新记录
func UpsertMatchCompetitionData(matchCompetitionData *model.MatchCompetitionData) (int64, error) {
	db := Video().Model(&matchCompetitionData).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "competition_id"}, {Name: "match_class"}}, //冲突的字段：competition_id 为唯一键
		DoUpdates: clause.AssignmentColumns([]string{"category_id", "match_class", "country_id", "name",
			"short_name", "name_zh", "name_zht", "logo", "type", "cur_season_id", "cur_stage_id", "cur_round", "round_count", "updated_timestamp"}),
	}).Create(&matchCompetitionData)
	return db.RowsAffected, db.Error
}

// UpsertMatchOriginData 插入/更新记录
func UpsertMatchOriginData(matchOriginData *model.MatchOriginData) (int64, error) {
	db := Video().Model(&matchOriginData).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "match_id"}, {Name: "match_class"}}, //冲突的字段：match_id 为唯一键
		DoUpdates: clause.AssignmentColumns([]string{"match_name", "match_class", "home_id", "home_name",
			"home_logo", "visit_id", "visit_name", "visit_logo", "venue_name", "league_id", "league_logo", "live_status", "m_live",
			"line_up", "stage_code", "group_num", "round_num", "start_at", "updated_timestamp"}),
	}).Create(&matchOriginData)
	return db.RowsAffected, db.Error
}

// UpsertMatchLanguageData 插入/更新记录
func UpsertMatchLanguageData(matchOriginData *model.MatchLanguageData) (int64, error) {
	db := Video().Model(&matchOriginData).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "data_id"}, {Name: "match_class"}, {Name: "data_type"}}, //冲突的字段：data_id 为唯一键
		DoUpdates: clause.AssignmentColumns([]string{"name_en", "name_zh", "name_zht", "updated_timestamp"}),
	}).Create(&matchOriginData)
	return db.RowsAffected, db.Error
}

// UpsertMatchVideoUrlData 插入/更新记录
func UpsertMatchVideoUrlData(matchOriginData *model.MatchVideoData) (int64, error) {
	db := Video().Model(&matchOriginData).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "match_id"}, {Name: "match_class"}}, //冲突的字段： 为唯一键
		DoUpdates: clause.AssignmentColumns([]string{"start_at", "updated_timestamp", "push_url1", "push_url2"}),
	}).Create(&matchOriginData)
	return db.RowsAffected, db.Error
}

func BatchInsertLSCategory(tx *gorm.DB, list []*model.MatchCategoryData) (int64, error) {
	if tx == nil {
		tx = Video()
	}
	db := tx.Model(&model.MatchCategoryData{}).CreateInBatches(list, len(list))
	return db.RowsAffected, db.Error
}

func UpdateLsCategory(matchCategoryData *model.MatchCategoryData) (int64, error) {
	var update = make(map[string]interface{})
	if matchCategoryData.NameZh != "" {
		update["name_zh"] = matchCategoryData.NameZh
	}
	if matchCategoryData.NameZht != "" {
		update["name_zht"] = matchCategoryData.NameZht
	}
	db := Video().Model(&model.MatchCategoryData{}).Where("category_id = ?", matchCategoryData.CategoryID).Updates(update)
	err := db.Error
	if err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func DeleteLSCategory(tx *gorm.DB, matchClass string) (int64, error) {
	if tx == nil {
		tx = Video()
	}
	db := tx.Table((&model.MatchCategoryData{}).TableName()).Where("match_class = ?", matchClass).Delete(&model.MatchCategoryData{})
	err := db.Error
	if err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func UpdateLsCompetition(matchCompetitionData *model.MatchCompetitionData) (int64, error) {
	var update = make(map[string]interface{})
	if matchCompetitionData.NameZh != "" {
		update["name_zh"] = matchCompetitionData.NameZh
	}
	if matchCompetitionData.NameZht != "" {
		update["name_zht"] = matchCompetitionData.NameZht
	}
	db := Video().Model(&model.MatchCompetitionData{}).Where("competition_id = ?", matchCompetitionData.CompetitionID).Updates(update)
	err := db.Error
	if err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func QueryLSTeamByIds(ids []string) ([]*model.MatchTeamData, error) {
	var data []*model.MatchTeamData
	err := VideoSlave().Model(&model.MatchTeamData{}).Where("team_id in (?)", ids).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func QueryLSCompetitionByIds(ids []string) ([]*model.MatchCompetitionData, error) {
	var data []*model.MatchCompetitionData
	err := VideoSlave().Model(&model.MatchCompetitionData{}).Where("competition_id in (?)", ids).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetLastCompetitionUpdatedTime(matchClass string) int {
	var matchCompetitionData *model.MatchCompetitionData
	err := Video().Model(&model.MatchCompetitionData{}).Select(" MAX(updated_timestamp) as updated_timestamp").Where("match_class = ?", matchClass).First(&matchCompetitionData).Error
	if err != nil {
		return 0
	}
	return matchCompetitionData.UpdatedTimestamp
}

func GetLastMatchUpdatedTime(matchClass string) int {
	var matchOriginData *model.MatchOriginData
	err := Video().Model(&model.MatchOriginData{}).Select(" MAX(updated_timestamp) as updated_timestamp").Where("match_class = ?", matchClass).First(&matchOriginData).Error
	if err != nil {
		return 0
	}
	return matchOriginData.UpdatedTimestamp
}

func GetLastLanguageUpdatedTime(matchClass string, dataType int) int {
	var matchLanguageData *model.MatchLanguageData
	err := Video().Model(&model.MatchLanguageData{}).Select(" MAX(updated_timestamp) as updated_timestamp").Where("match_class = ? and data_type = ?", matchClass, dataType).First(&matchLanguageData).Error
	if err != nil {
		return 0
	}
	return matchLanguageData.UpdatedTimestamp
}

func QueryLSLanguageByIds(ids []string, matchClass string, dataType int) ([]*model.MatchLanguageData, error) {
	var data []*model.MatchLanguageData
	err := VideoSlave().Model(&model.MatchLanguageData{}).Where("data_id in (?) and match_class = ? and data_type= ?", ids, matchClass, dataType).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func QueryLSMatchVideoUrlData(startDate, endDate string) ([]*model.MatchVideoUrlData, error) {
	datas := make([]*model.MatchVideoUrlData, 0)
	err := VideoSlave().Table("match_video_data as v").Joins("inner join match_origin_data as o on v.match_id = o.match_id and v.match_class = o.match_class").
		Select("v.id,v.match_id,v.match_class,v.start_at,v.push_url1,"+
			"v.push_url2,o.match_name,o.home_name,o.home_logo,o.visit_name,o.visit_logo,"+
			"o.venue_name,o.league_logo,o.plate_id,o.ball_class_id").
		Where("v.start_at >= ?", startDate).Where("v.start_at <= ?", endDate).Find(&datas).Error
	if err != nil {
		return datas, err
	}
	return datas, nil
}
