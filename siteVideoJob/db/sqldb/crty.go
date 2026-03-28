package sqldb

import (
	"gorm.io/gorm/clause"
	"siteVideoJob/mdata/model"
)

func QueryCRMatchData(startDate, endDate string) ([]*model.MatchCrtyOriginData, error) {
	matchCrtyOriginData := make([]*model.MatchCrtyOriginData, 0)
	err := VideoSlave().Model(&matchCrtyOriginData).Where("start_at >= ?", startDate).Where("start_at <= ?", endDate).Find(&matchCrtyOriginData).Error
	if err != nil {
		return matchCrtyOriginData, err
	}
	return matchCrtyOriginData, nil
}

// UpsertMatchCrtyOriginData 插入/更新记录
func UpsertMatchCrtyOriginData(matchCrtyOriginData *model.MatchCrtyOriginData) (int64, error) {
	db := Video().Model(&matchCrtyOriginData).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "match_id"}, {Name: "match_class"}}, //冲突的字段：match_id,match_class 为唯一键
		DoUpdates: clause.AssignmentColumns([]string{"league_id", "match_name", "venue_name",
			"start_at", "home_name", "home_logo", "visit_name", "visit_logo", "league_logo", "live_status"}),
	}).Create(&matchCrtyOriginData)
	return db.RowsAffected, db.Error
}
