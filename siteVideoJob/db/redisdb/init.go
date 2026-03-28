package redisdb

import (
	"siteVideoJob/db/redisdb/core"
	"siteVideoJob/db/redisdb/game"
	"siteVideoJob/internal/glog"
)

func InitRedis() error {

	err := core.InitRedis()
	if err != nil {
		glog.Emergency("InitCoreRedis |err=%v", err)
		return err
	}
	err = game.InitRedis()
	if err != nil {
		glog.Emergency("InitGameRedis |err=%v", err)
		return err
	}

	return nil
}
