package redisdb

import (
	"siteLetterJob/db/redisdb/core"
	"siteLetterJob/internal/glog"
)

func InitRedis() error {

	err := core.InitRedis()
	if err != nil {
		glog.Emergency("init redis err: ", err)
		return err
	}
	//todo 用到再打开
	//err = agent.InitRedis()
	//if err != nil {
	//	glog.Error(err)
	//	return err
	//}
	//
	//todo 用到再打开
	//err = game.InitRedis()
	//if err != nil {
	//	glog.Error(err)
	//	return err
	//}

	//todo 暂时没用到
	//err = sso.InitRedis()
	//if err != nil {
	//	glog.Error(err)
	//	return err
	//}

	return nil
}
