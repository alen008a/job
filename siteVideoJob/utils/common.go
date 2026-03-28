package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/rs/xid"
	"github.com/shopspring/decimal"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func NumConvert(memberId int) int {
	return memberId & 63
}

func MaskRealName(realName string) string {
	if realName == "" {
		return ""
	} else {
		strings.Trim(realName, " ")
		return Overlay(realName, "**", 1, len(realName))
	}
}

// MaskIp
func MaskIp(ip string) string {
	if ip == "" {
		return ""
	} else {
		ipArray := strings.Split(ip, ".")
		if len(ipArray) >= 2 {
			return ipArray[0] + "." + ipArray[1] + ".*.*"
		} else {
			return ip
		}
	}
}
func MaskPhone(phone string) string {
	if phone == "" {
		return ""
	} else {
		strings.Trim(phone, " ")
		return Overlay(phone, "****", 3, 7)
	}
}

func MaskEmail(email string) string {
	if email == "" {
		return ""
	} else {
		strings.Trim(email, " ")
		at := "@"
		if !strings.Contains(email, at) {
			return email
		} else {
			options := strings.Split(email, at)
			if len(options) < 2 {
				return email
			} else {
				return Overlay(options[0], "****", 2, len(options[0])) + at + options[1]
			}
		}
	}
}

func MaskAddress(address string) string {
	if address == "" {
		return ""
	} else {
		addressArray := strings.Split(address, ",")
		sb := strings.Builder{}
		if len(addressArray) > 1 {
			sb.WriteString(addressArray[0])
			sb.WriteString(" ")
			sb.WriteString(addressArray[1])
			sb.WriteString(" ****")
			sb.WriteString(" ****")
		} else {
			sb.WriteString(addressArray[0])
			sb.WriteString(" ****")
			sb.WriteString(" ****")
		}
		return sb.String()
	}
}

func MaskBankNum(bankNum string) string {
	if bankNum == "" {
		return ""
	} else {
		return Overlay(bankNum, "**** **** **** ", 0, len(bankNum)-4)
	}
}

func MaskQq(qq string) string {
	if qq == "" {
		return ""
	} else {
		strings.Trim(qq, " ")
		return Overlay(qq, "****", 2, len(qq))
	}
}

func PageNUms(total, pageSize int) (pageNums int) {
	if total < 1 || pageSize < 1 {
		return
	}
	pageNums = int(math.Ceil(float64(total) / float64(pageSize)))
	return
}

//数组用
func PageOffsetAndEnd(total, pageSize, page int) (pageNums, offset, end int) {
	pageNums = PageNUms(total, pageSize)
	if pageNums < 1 || page > pageNums {
		return
	}
	offset = (page - 1) * pageSize
	end = offset + pageSize
	if end > total {
		end = total
	}
	return
}

// Scale 四舍五入
func Scale(in float64, num int32) float64 {
	res, _ := decimal.NewFromFloat(in).Round(num).Float64()
	return res
}

// RandFileName 生成文件名 并发模式下注意文件名会冲突
func RandFileName(fileType string) string {
	id := xid.New()
	fileName := fmt.Sprintf("%v_%d%s",
		id,
		RandNum(100000, 999999),
		strings.ToLower(fileType),
	)
	return fileName
}

/**
 * 解决四舍五入保留后2位精确到小数点后2位 返回decimal类型
 * 如果要转字符串+ .String() 转float .Float64()
 */
func MathRoundFloat(number float64) decimal.Decimal {
	value := decimal.NewFromFloat(number).Round(2)
	return value
}

func HostName() string {
	hostname, _ := os.Hostname()
	return hostname
}

func GetIMTimestamp(key string) (result string) {
	has := md5.Sum([]byte(key))
	key = fmt.Sprintf("%x", has)

	lens := len(key) / 2
	var md5raw string
	for i := 0; i < lens; i++ {
		hexByte, _ := hex.DecodeString(Substring(key, i*2, (i*2)+2))
		md5raw = md5raw + string(hexByte)
	}

	preDayTime := time.Now().Unix() - 12*3600
	nanoStr := strconv.FormatInt(time.Now().UnixNano(), 10)
	timeStr := time.Unix(preDayTime, 0).Format("2006-01-02 15:04:05")
	timeStamp := timeStr + "." + Substring(nanoStr, len(nanoStr)-3, len(nanoStr))

	return AesECBEncrypt(timeStamp, md5raw)
}
