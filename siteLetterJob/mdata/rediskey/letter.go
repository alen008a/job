package rediskey

const (
	BatchReadBRV2                            = "site_msg_batch_read_br_%v_%v_%v"                          // 广播消息，会员批量读取，记录id集合,siteid,memberId,msgType
	UnreadMsgCacheMemberBroadcastUnreadCount = "site_msg_cache_member_broadcast_unread_count_v3_%v_%v_%v" //广播站内信未读数 站点id ,memberId,msgType
	UnreadMsgCacheMemberBroadcastListCount   = "site_msg_cache_member_broadcast_list_count_v3_%v_%v_%v"   //广播站内信列表缓存 站点id ,memberId,msgType
)
