package mdata

import (
	jsoniter "github.com/json-iterator/go"
	extra "github.com/json-iterator/go/extra"
)

const (
	HeaderSite = "X-API-SITE" // 平台标志
	XMTY       = "XMTY"       //熊猫体育
	CRTY       = "CRTY"       //CR体育
	LSTY       = "LSTY"       //雷速体育
	System     = "SYSTEM"     //操作人默认系统
)

type PullType int //拉取方式

type PullReqType int //雷速体育参数类型

type GameType string //游戏类别

const (
	PullFull PullType = 1 //全量
	PullIncr PullType = 2 //增量
)

const (
	PagePullReqType PullReqType = 1 //page 按页码分页查询
	TimePullReqType PullReqType = 2 //time 按三方更新时间戳查询
)

const (
	FT GameType = "FT"
	BK GameType = "BK"
	TN GameType = "TN"
)

func init() {
	extra.RegisterFuzzyDecoders()
}

var Cjson = jsoniter.ConfigCompatibleWithStandardLibrary

func MustMarshal2Byte(v interface{}) []byte {
	if v == nil {
		return []byte{}
	}

	b, _ := Cjson.Marshal(v)
	return b
}

func MustMarshal2String(v interface{}) string {
	if v == nil {
		return ""
	}

	b, _ := Cjson.MarshalToString(v)
	return b
}

// IPLoc ip信息
type IPLoc struct {
	Country  string // 国家
	Province string // 省份
	City     string // 城市
}

type SendSlack struct {
	Code string `json:"code"`
	Key  string `json:"key"`
	Msg  string `json:"msg"`
}

type VideoSourceApi struct {
	ChannelName   string   `json:"channelName"`
	BkApi         string   `json:"bkApi"`
	MatchClassArr []string `json:"matchClassArr"`
	MatchClass    string   `json:"-"`
}
