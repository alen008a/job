package metadata

import (
	"siteVideoJob/internal/context"
	"strconv"
)

func GetSiteId(c *context.Context) int {
	if c.SiteId == "" {
		return 0
	}
	siteIdInt, _ := strconv.Atoi(c.SiteId)
	return siteIdInt
}
