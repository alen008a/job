package sqldb

import (
	"siteVideoJob/mdata/model"
	"siteVideoJob/mdata/video"
)

func QueryVenueStatus(venueName string) (v *video.VenueStatusInfo, err error) {
	err = sqlDB.SiteSlave.Table("").Select("status,is_display as isDisplay,en_name as enName").Where("en_name = ? ", venueName).First(&v).Error
	return
}

// 查询所有场馆的排序分数（越小的越前面）
// 返回 Map: en_name => score
func QueryVenuesSortScore() (venueSortScore map[string]int, err error) {
	var venues []model.GameVenue
	venueSortScore = make(map[string]int)

	// 查询所有场馆的排序 sort (越小的越前面)
	err = SiteSlave().Table("game_venue").
		Select("en_name", "sort").
		Where("status <> 2").
		Where("game_type IN (?)", []string{"TY", "DJ"}).
		Order("sort asc").
		Find(&venues).Error

	// 转换成 en_name => sort 的 map
	for _, venue := range venues {
		venueSortScore[venue.EnName] = venue.Sort
	}

	return venueSortScore, err
}

// 查询体育/电竞赛事的冠名场馆
func QuerySpecificVenues() (specificVenue map[string]byte, err error) {
	var venues []model.GameVenue
	specificVenue = make(map[string]byte)

	// 建立子查询，这里只拿体育/电竞相关的资料
	// 参考资料： KY-40335
	subquery := SiteSlave().Table("game_venue").
		Select("substring_index(group_concat(id order by sort,sort_updated_at desc),',',1) as id").
		Where("status <> 2").
		Where("game_type IN (?)", []string{"TY", "DJ"})

	// 主查询
	err = SiteSlave().Table("game_venue").
		Select("en_name").
		Where("id IN (?)", subquery).
		Find(&venues).Error

	if err != nil {
		return nil, err
	}

	for _, venue := range venues {
		specificVenue[venue.EnName] = 1
	}

	return specificVenue, nil
}
