package sqldb

import (
	"errors"
	"siteVideoJob/internal/glog"
	"siteVideoJob/mdata/model"
	"siteVideoJob/mdata/video"
	"siteVideoJob/utils"
	"time"
)

func QueryFinishedInfos() (list []model.LiveEventsConfig, err error) {
	err = SiteSlave().Table("live_events_config").Where("start_at < ? ", time.Now().AddDate(0, 0, -1).Format(utils.TimeBarFormat)).Find(&list).Error
	return
}

func BatchDeleteEvents(matchIds []int64) (int64, error) {
	tx := Site().Table("live_events_config").Where("match_id in (?) ", matchIds).Delete(model.LiveEventsConfig{})
	return tx.RowsAffected, tx.Error
}

func DeleteLastDayMatch() (int64, error) {
	tx := Site().Table("live_events_config").Where("start_at < ? ", time.Now().AddDate(0, 0, -1).Format(utils.TimeBarFormat)).Delete(model.LiveEventsConfig{})
	return tx.RowsAffected, tx.Error
}

func IsExistEvents(channel string, siteId, eventsId int64) bool {
	var cnt int64
	err := SiteSlave().Table("live_events_config").Where("site_id = ? and match_id = ? and venue_name = ?", siteId, eventsId, channel).Count(&cnt).Error
	if err != nil {
		glog.Errorf("IsExistEvents count error. channel=%v eventsId=%v", channel, eventsId)
		return false
	}
	if cnt > 0 {
		return true
	}
	return false
}

func QueryEventsInfos(siteId int, channel string, eventsIds []int64) (list []model.LiveEventsConfig, err error) {
	err = SiteSlave().Table("live_events_config").Where("site_id = ? and match_id in (?) and venue_name = ? ", siteId, eventsIds, channel).Find(&list).Error
	return
}

func BatchInsertEvents(list []model.LiveEventsConfig) (int64, error) {
	db := Site().Table("live_events_config").CreateInBatches(list, len(list))
	return db.RowsAffected, db.Error
}

func UpdateVideo(req *model.LiveEventsConfig) (int64, error) {
	db := Site().Table("live_events_config").Updates(req)
	return db.RowsAffected, db.Error
}

func UpdateVideoPlatStatus(id, platId, liveStatus int) error {
	db := Site().Table("live_events_config").Where("id = ?", id).
		Updates(map[string]interface{}{"plate_id": platId, "live_status": liveStatus})
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return errors.New("update rowsAffected 0")
	}
	return nil
}

func LiveEventsList(req video.LiveEventsDto) (list video.LiveEventList, err error) {
	db := SiteSlave().Table("live_events_config").Where("site_id = ? and data_source = 0", req.SiteId)
	if req.DelStatus != -1 {
		db.Where("del_flag = ? and status = ?", req.DelStatus, req.Status)
	}
	if len(req.StartAt) > 0 {
		db.Where("start_at >= ?", req.StartAt)
	}
	if len(req.VenueName) > 0 {
		db.Where("venue_name = ?", req.VenueName)
	}
	db.Where("start_at < ?", time.Now().Add(time.Hour*8).Format(utils.TimeBarFormat))
	err = db.Order("major_event desc,start_at asc,sort asc").Find(&list).Error
	return
}
