/*
*

	@note:

*
*/
package seaweedfs

import "strings"

// 上传参数
type UploadResp struct {
	FileName string `json:"name"`    // 文件名字
	FileUrl  string `json:"fileURL"` // 文件完整路径
	Size     int64  `json:"size"`    // 文件大小
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

const (
	HTTP_TIMEOUT_DEFAULT int = 10 // 单位：秒
)
