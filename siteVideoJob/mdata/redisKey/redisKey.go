package redisKey

const (
	LiveVenueStatus    = "live_venue_status_%s_%d"
	VideoMatchMapKey   = "video_match_key_map_%s" //赛事数据更新 {channel}
	VenueSortList      = "venue_sort_list_%d"     // 从 Database 中捞取的场馆分数
	VenueLiveSportsKey = "venue_live_sports_%s_%d_%s"
	VenueSourceKey     = "venue_source_live_%s"   //视频源数据
	AnchorEventsList   = "anchor_events_list_%d"  //前台主播赛事源数据  hash：场次vid=》json场次信息 需要清理并更新状态
	AnchorLiveVenue    = "anchor_live_venue_%d"   //直播场馆集合 zset
	AnchorShowDefault  = "anchor_show_default_%d" //缺省展示主播 hash
	AnchorShowIds      = "anchor_show_ids_%d"     //状态为展示的主播id  hash：场次id=》开始时间
	AnchorShowUpdate   = "anchor_show_update_%d"  //运营管理-主播展示-展示状态的修改时间和操作者  hash:场次id-》操作数据
	AnchorShowList     = "anchor_show_list_%d"    //前台接口数据  string(json)
	XMTYVideoTokenKey  = "XMTY_video_token"       // 视频链接上的token
)
