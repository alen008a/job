package db

import (
	"siteLetterJob/db/redisdb"
	"siteLetterJob/db/sqldb"
)

// InitDB 初始化 db
func InitDB() error {
	var err error

	// mysql 的初始化
	err = sqldb.InitMysql()
	if err != nil {
		return err
	}

	// 初始化 redis
	err = redisdb.InitRedis()
	if err != nil {
		return err
	}

	return nil
}

// Close close all db
func Close() error {
	return nil
}
