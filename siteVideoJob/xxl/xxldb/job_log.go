package xxldb

import "siteVideoJob/db/sqldb"

type JobLog struct {
	ID        int64  `gorm:"column:id" json:"id"`                // id
	HandleMsg string `gorm:"column:handle_msg" json:"handleMsg"` // 执行日志内容
}

func (j JobLog) GetTableName() string {
	return "xxl_job_log"
}

func GetMsgById(logId int64) (string, error) {
	var jobLog JobLog
	err := sqldb.XxlJobDB().Table(jobLog.GetTableName()).
		Select("handle_msg").
		Where("id = ?", logId).
		First(&jobLog).Error
	if err != nil || &jobLog == nil {
		return "", err
	}
	return jobLog.HandleMsg, nil
}
