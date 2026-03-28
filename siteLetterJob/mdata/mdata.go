package mdata

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
	extra "github.com/json-iterator/go/extra"
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

const (
	HeaderSite = "X-API-SITE" // 平台标志
	WEB        = "web"
	AGENT_WEB  = "agent_web"
	H5         = "h5"
	Android    = "android"
	IOS        = "ios"
	PC         = "pc"
)

var (
	TokenEmptyErr = errors.New("token is empty")
)

// IPLoc ip信息
type IPLoc struct {
	Country  string // 国家
	Province string // 省份
	City     string // 城市
}

var (
	// ip/uuid校验
	IgnoreIp = map[string]int{
		WEB: 0,
		H5:  1,
	}
	IgnoreClientType = map[string]int{
		"agent_web": 0,
		"1":         1,
	}
	// 平台映射
	PlatformEnum = map[int]string{
		1:    "wm1",
		2:    "wm2",
		3:    "oubao",
		4:    "huohu",
		5:    "lol",
		6:    "yibo",
		7:    "coc",
		8:    "huanqiu",
		10:   "fenghuang",
		11:   "qiusu",
		1000: "完美测试",
		1001: "完美",
		2000: "完美2测试",
		2001: "完美2",
		3000: "完美3测试",
		3001: "完美3",
		4000: "完美4测试",
		4001: "完美4",
	}

	// 平台映射map
	PlatformTypeMap = map[string]string{
		"0": "web",
		"1": "site",
	}
)
