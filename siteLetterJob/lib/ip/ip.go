// Package libip ip 信息
package libip

import (
	"github.com/ipipdotnet/ipdb-go"
	"siteLetterJob/config"
	"siteLetterJob/internal/context"
	"siteLetterJob/mdata"
)

var (
	ipdbInfo *ipdb.City
)

// 初始化 ip 库
func InitIP() error {
	var err error

	ipdbInfo, err = ipdb.NewCity(config.GetServiceAddr().IPConfAddr)
	if err != nil {
		return err
	}

	return nil
}

func GetIPLoc(c *context.Context, ip string) *mdata.IPLoc {
	var data mdata.IPLoc
	data.Country = "中国"

	if ipdbInfo == nil {
		c.Errorf("ipdbInfo is nil")
		return &data
	}

	loc, err := ipdbInfo.FindInfo(ip, "CN")
	if err != nil {
		c.Errorf("ip=%s query err: %v", ip, err)
		return &data
	}

	data.Country = loc.CountryName
	data.Province = loc.RegionName
	data.City = loc.CityName

	return &data
}
