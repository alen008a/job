package snowid

import (
	"strconv"
	"time"

	"siteVideoJob/utils"

	"github.com/sony/sonyflake"
)

const retry = 10
const idLength = 19

var gSnowflake *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	gSnowflake = sonyflake.NewSonyflake(st)
	if gSnowflake == nil {
		panic("snowflake not created")
	}
}

func SnowflakeId() int64 {
	var id uint64
	var err error
	for i := 0; i < retry; i++ {
		id, err = gSnowflake.NextID()
		if err != nil {
			continue
		}

		break
	}

	//判断长度
	idStr := strconv.FormatUint(id, 10)

	if id == 0 {
		//时间戳+6位随机数
		idStr = strconv.FormatInt(time.Now().Unix(), 10) + utils.RealRandNumber(6)
	} else {
		//长度补齐,如果生成雪花id长度位0，自动获取19位随机数
		delta := idLength - len(idStr)
		if delta > 0 {
			idStr += utils.RealRandNumber(delta)
		}
	}

	r, _ := strconv.ParseInt(idStr, 10, 64)

	return r
}
