package xxl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"siteVideoJob/xxl/xxldb"
)

/**
用来日志查询，显示到xxl-job-admin后台
*/

type LogHandler func(req *LogReq) *LogRes

// 默认返回
func defaultLogHandler(req *LogReq) *LogRes {
	return &LogRes{Code: SuccessCode, Msg: "", Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   2,
		LogContent:  "这是日志默认返回，说明没有设置LogHandler",
		IsEnd:       true,
	}}
}

// CustomLogHandle 自定义日志处理器
func CustomLogHandle(req *LogReq) *LogRes {
	return &LogRes{
		Code: SuccessCode,
		Msg:  "load successfully!",
		Content: LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  "自定义日志handler<br>",
			IsEnd:       true,
		}}
}

// GetDBLogHandle 从数据库获取执行日志
func GetDBLogHandle(req *LogReq) *LogRes {
	var result = LogRes{
		Code: SuccessCode,
		Msg:  "load successfully!",
		Content: LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			IsEnd:       true,
		}}
	if req.LogID <= 0 {
		msg := fmt.Sprintf("<p color='red'>查询执行记录失败，参数日志Id:%v 不合法!</p>", req.LogID)
		result.Content.LogContent = msg
		return &result
	}
	logMsg, err := xxldb.GetMsgById(req.LogID)
	if err != nil {
		logMsg = fmt.Sprintf("<p color='yellow'>查询执行记录不存在! 日志Id:%v </p>", req.LogID)
		result.Content.LogContent = logMsg
		return &result
	}
	result.Content.LogContent = logMsg
	return &result
}

// 请求错误
func reqErrLogHandler(w http.ResponseWriter, req *LogReq, err error) {
	logResp := &LogRes{Code: FailureCode, Msg: err.Error(), Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   0,
		LogContent:  err.Error(),
		IsEnd:       true,
	}}
	str, _ := json.Marshal(logResp)
	_, _ = w.Write(str)
}
