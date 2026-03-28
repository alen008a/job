package sqldb

// QueryALLSite 获取所有开启站点的站点ID
func QueryALLSite() (siteIds []int, err error) {
	err = ControlSlave().Table("site_manage_info").Where("del_flag = 0").Order("site_id asc").Pluck("site_id", &siteIds).Error
	return
}
