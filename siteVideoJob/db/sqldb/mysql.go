package sqldb

import (
	"errors"
	"siteVideoJob/mdata"
	"strings"
	"time"

	"siteVideoJob/config"
	"siteVideoJob/internal/glog"
	"siteVideoJob/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var sqlDB *SqlDB
var dbPoolInit bool
var dbPoolMap map[string]DbPool

type SqlDB struct {
	Site          *gorm.DB
	SiteSlave     *gorm.DB
	ControlSlave  *gorm.DB
	Video         *gorm.DB
	VideoSlave    *gorm.DB
	XxlJob        *gorm.DB
	closeFunction []func()
}

type DbPool struct {
	MaxIdle int `json:"maxIdle"`
	MaxConn int `json:"maxConn"`
	MaxLife int `json:"maxLife"`
}

func Site() *gorm.DB {
	return sqlDB.Site
}

func SiteSlave() *gorm.DB {
	return sqlDB.SiteSlave
}

func Video() *gorm.DB {
	return sqlDB.Video
}

func VideoSlave() *gorm.DB {
	return sqlDB.VideoSlave
}

func ControlSlave() *gorm.DB {
	return sqlDB.ControlSlave
}

func XxlJobDB() *gorm.DB {
	return sqlDB.XxlJob
}

// getDbPool 根据db名称获取连接池配置
func getDbPool(dbName string, conf config.Mysql) DbPool {
	dbPools := config.GetConfig().Application.DbPools
	// 不会并发，用一个bool变量表示是否初始化过即可
	if !dbPoolInit {
		if dbPools == "" {
			dbPoolMap = make(map[string]DbPool, 0)
			glog.Info("未配置私有DB连接池")
		} else {
			err := mdata.Cjson.UnmarshalFromString(dbPools, &dbPoolMap)
			if err != nil {
				glog.Errorf("私有DB连接池配置格式错误：%s", dbPools)
				dbPoolMap = make(map[string]DbPool, 0)
			}
			glog.Infof("私有DB连接池配置为：%+v", dbPoolMap)
		}
		dbPoolInit = true
	}
	var (
		p  DbPool
		ok bool
	)
	if p, ok = dbPoolMap[dbName]; ok {
		//私有配置中的某些值为0或没配时，使用公共
		if p.MaxIdle == 0 {
			p.MaxIdle = conf.IdleConnect
			glog.Infof("%s连接池的MaxIdle参数使用公共配置值：%d", dbName, p.MaxIdle)
		}
		if p.MaxConn == 0 {
			p.MaxConn = conf.MaxConnect
			glog.Infof("%s连接池的MaxConn参数使用公共配置值：%d", dbName, p.MaxConn)
		}
		if p.MaxLife == 0 {
			p.MaxLife = conf.MaxLifeTime
			glog.Infof("%s连接池的MaxLife参数使用公共配置值：%d", dbName, p.MaxLife)
		}
	} else { //无私有配置，使用公共
		p.MaxIdle = conf.IdleConnect
		p.MaxConn = conf.MaxConnect
		p.MaxLife = conf.MaxLifeTime
		glog.Infof("没有找到%s的私有连接池配置，使用公共值：%+v", dbName, p)
	}
	return p
}

// InitMysql init mysql
func InitMysql() (err error) {
	sqlDB = new(SqlDB)

	sqlDB.Site, err = initSqlDB(config.GetConfig().Site, "Site")
	if err != nil {
		glog.Emergency("init Site db is err: %v", err)
		return err
	}

	sqlDB.SiteSlave, err = initSqlDB(config.GetConfig().SiteSlave, "SiteSlave")
	if err != nil {
		glog.Emergency("init SiteSlave db is err: %v", err)
		return err
	}

	sqlDB.Video, err = initSqlDB(config.GetConfig().Video, "Video")
	if err != nil {
		glog.Emergency("init Video db is err: %v", err)
		return err
	}

	sqlDB.VideoSlave, err = initSqlDB(config.GetConfig().VideoSlave, "VideoSlave")
	if err != nil {
		glog.Emergency("init VideoSlave db is err: %v", err)
		return err
	}

	sqlDB.ControlSlave, err = initSqlDB(config.GetConfig().ControlSlave, "ControlSlave")
	if err != nil {
		glog.Emergency("init ControlSlave db is err: %v", err)
		return err
	}

	sqlDB.XxlJob, err = initSqlDB(config.GetConfig().XxlJob, "XxlJob")
	if err != nil {
		glog.Emergency("init XxlJob db is err: %v", err)
		return err
	}

	return nil
}

// 初始化数据库
func initSqlDB(c config.Mysql, dbName string) (*gorm.DB, error) {
	var logLevel logger.LogLevel
	if c.LogEnable {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}
	if c.Address == "" {
		return nil, errors.New("数据库：" + dbName + "地址为空")
	}
	address := utils.GetRealString(config.GetConfig().DBSecretKey, c.Address)

	//调试，追加预编译相关参数
	if strings.Contains(address, "?") {
		address = address + "&"
	} else {
		address = address + "?"
	}
	address = address + "interpolateParams=true"

	glog.Infof("initSqlDB |address=%s", strings.Replace(c.Address, `\u0026`, "&", -1))
	db, err := gorm.Open(
		mysql.New(
			mysql.Config{
				DSN: address,
			},
		), &gorm.Config{
			PrepareStmt:            false,
			SkipDefaultTransaction: true,
			Logger: glog.NewDBLog(
				logger.Config{
					SlowThreshold:             time.Second,
					Colorful:                  true,
					IgnoreRecordNotFoundError: true,
					LogLevel:                  logLevel,
				},
			),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	)
	if err != nil {
		glog.Emergency("addr=%s |err=%v", address, err)
		return nil, err
	}

	dbp, err := db.DB()
	if err != nil {
		glog.Emergency("addr=%s |err=%v", address, err)
		return nil, err
	}

	p := getDbPool(dbName, c)
	dbp.SetMaxOpenConns(p.MaxConn)
	dbp.SetMaxIdleConns(p.MaxIdle)
	dbp.SetConnMaxLifetime(time.Duration(p.MaxLife) * time.Second)

	if err = dbp.Ping(); err != nil {
		glog.Emergency("addr=%s |err=%v", c.Address, err)
		return nil, err
	}

	sqlDB.closeFunction = append(
		sqlDB.closeFunction, func() {
			dbp.Close()
		},
	)
	glog.Infof("连接池配置，%s库，MaxIdle = %d，MaxConn = %d，MaxLife = %d秒", dbName, p.MaxIdle, p.MaxConn, p.MaxLife)
	return db, nil
}

func Close() {
	for i := range sqlDB.closeFunction {
		sqlDB.closeFunction[i]()
	}
}
