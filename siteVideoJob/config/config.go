package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/ini.v1"
	"siteVideoJob/mdata"
	"siteVideoJob/mdata/namespace"
	"strings"
)

var (
	g   *Configurations
	PRO = "pro" //生产环境
	DEV = "dev" //测试环境
)

const GroupGlobal = "Global"

// Configurations 存放当前的配置
type Configurations struct {
	// 日志配置
	Logger Logger

	// siteVideoJob应用配置
	Application

	// 核心redis
	RedisCore Redis
	// 场馆Redis
	RedisGame Redis
	// 总控从库
	ControlSlave Mysql
	// 主站主库
	Site Mysql
	// 主站从库
	SiteSlave Mysql

	// 赛事数据源主库
	Video Mysql
	// 赛事数据源从库
	VideoSlave Mysql

	// xxl调度平台数据库
	XxlJob Mysql

	Kafka

	// 全局通用配置
	Common

	NameMapping NameMapping

	f map[namespace.NacosNamespace]interface{}

	ago    config_client.IConfigClient
	global config_client.IConfigClient

	app NacosConfig
}

type Kafka struct {
	KafkaAddr string
}

func GetConfig() *Configurations {
	return g
}

type NacosConfig struct {
	Nacos `ini:"Nacos"`
}

type Common struct {
	ProxyURL         string
	VerifyCodeDomain string // 短信、语音、email验证码，银行卡二要素服务域名
	VerifyProxyUrl   string
	WarningCode      string //应用程序告警
	DBSecretKey      string // 密钥
}

func GetApp() Nacos {
	return g.app.Nacos
}

// Nacos 从本地ini文件读取Nacos配置
type Nacos struct {
	Env         string `ini:"Env"`         //环境
	NamespaceId string `ini:"NamespaceId"` //命名空间Id
	AppID       string `ini:"AppID"`       //应用id
	Address     string `ini:"Address"`     //Nacos地址
	Port        uint64 `ini:"Port"`        //端口号
	Scheme      string `ini:"Scheme"`      //协议
	Username    string `ini:"Username"`    //服务端的API鉴权Username
	Password    string `ini:"Password"`    //服务端的API鉴权Password
	CacheDir    string `ini:"CacheDir"`    //缓存service信息的目录，默认是当前运行目录
	LogDir      string `ini:"LogDir"`      //日志存储路径
}

// Application 从apollo读取基础配置
type Application struct {
	AdminServer  string // 调度中心地址 http://domain:port/path
	AccessToken  string // xxlAdmin token
	ExecutorName string // 执行器名称
	ExecutorPort string // 相同于于HttpPort
	ExecutorIp   string // 非必须

	VideoLiveSourceUrl     string // 播控地址
	VideoLiveSiteAlias     string // 播控代号
	VideoChannelList       string // 视频场馆列表
	VideoSourceApi         string // 直播源视频场馆
	LSTYMatchClass         string // 雷速体育需要拉取的球类类型
	LSTYMatchPullVideoHost string // 雷速体育拉流域名
	LSTYMatchPushVideoHost string // 雷速体育推流域名
	LiveStreamCdnKey       string // 直播云鉴权KEY
	AnchorSourceApi        string // 播控-主播接口
	AutoloadConfig         bool   // 项目配置热加载开关，true为热加载。
	DbPools                string // 应用私有连接池配置
}

type Mysql struct {
	Address     string
	LogEnable   bool
	IdleConnect int
	MaxConnect  int
	MaxLifeTime int
}

type Elastic struct {
	Address           string
	AuthPass          string
	AuthUser          string
	EnableRequestBody bool
}

type Redis struct {
	Host     string //如果是集群模式，配置的是节点的host:port,如果是哨兵模式，配置的是哨兵节点的host:port
	Auth     string
	Master   string //哨兵模式必配，其他模式不用配
	PoolSize int
}

type Logger struct {
	LogPath  string
	LogLevel string
	LogType  string
}

type NameMapping map[string]map[string]map[string]string

func GetGameType(enName, gameType string) string {
	if gameType == "" {
		return ""
	}
	if gameTypeMap, ok := GetConfig().NameMapping["GameType"]; ok {
		if venueMap, ok := gameTypeMap[enName]; ok {
			if gameTypeNew, ok := venueMap[gameType]; ok {
				return gameTypeNew
			}
		}
	}
	return gameType
}

func GetVideoSourceApi() []*mdata.VideoSourceApi {
	videoSourceApi := make([]*mdata.VideoSourceApi, 0)
	mdata.Cjson.UnmarshalFromString(GetConfig().VideoSourceApi, &videoSourceApi)
	return videoSourceApi
}

func GetLSTYMatchClassList() []string {
	matchClass := make([]string, 0)
	matchClass = strings.Split(GetConfig().LSTYMatchClass, ",")
	return matchClass
}

// GetEnv 获取当前环境
func GetEnv() string {
	return strings.ToLower(g.app.Env)
}

// InitConfig 初始化配置
func InitConfig(confPath string) error {
	apo := new(NacosConfig)
	if err := ini.MapTo(apo, confPath); err != nil {
		return err
	}

	g = new(Configurations)
	g.app = *apo
	g.f = make(map[namespace.NacosNamespace]interface{})

	err := InitNacos()
	if err != nil {
		return err
	}
	err = g.loadConfig()
	if err != nil {
		return err
	}
	// 读apollo配置
	if GetConfig().AutoloadConfig {
		fmt.Println("开启动态加载")
		go g.poller()
	} else {
		fmt.Println("未开启动态加载")
	}
	return nil
}

func InitNacos() error {
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(g.app.NamespaceId), //当namespace是public时，此处填空字符串。
		constant.WithTimeoutMs(90000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(g.app.LogDir),
		constant.WithUsername(g.app.Username),
		constant.WithPassword(g.app.Password),
		constant.WithCacheDir(g.app.CacheDir),
	)
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			g.app.Address,
			g.app.Port,
			constant.WithScheme(g.app.Scheme)),
	}
	var err error
	g.ago, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return err
	}

	clientGlobalConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(GroupGlobal),
		constant.WithTimeoutMs(90000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(g.app.LogDir),
		constant.WithUsername(g.app.Username),
		constant.WithPassword(g.app.Password),
		constant.WithCacheDir(g.app.CacheDir),
	)
	g.global, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientGlobalConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Configurations) loadConfig() error {
	err := g.decode(namespace.Logger, &g.Logger)
	if err != nil {
		return err
	}
	err = g.decode(namespace.Application, &g.Application)
	if err != nil {
		return err
	}
	err = g.decode(namespace.ControlSlave, &g.ControlSlave)
	if err != nil {
		return err
	}
	err = g.decode(namespace.Site, &g.Site)
	if err != nil {
		return err
	}
	err = g.decode(namespace.SiteSlave, &g.SiteSlave)
	if err != nil {
		return err
	}
	err = g.decode(namespace.Video, &g.Video)
	if err != nil {
		return err
	}
	err = g.decode(namespace.VideoSlave, &g.VideoSlave)
	if err != nil {
		return err
	}
	err = g.decode(namespace.XxlJobDbNamespace, &g.XxlJob)
	if err != nil {
		return err
	}
	err = g.decode(namespace.RedisCore, &g.RedisCore)
	if err != nil {
		return err
	}
	err = g.decode(namespace.RedisGame, &g.RedisGame)
	if err != nil {
		return err
	}
	err = g.decode(namespace.Kafka, &g.Kafka)
	if err != nil {
		return err
	}
	err = g.decode(namespace.Common, &g.Common)
	if err != nil {
		return err
	}
	err = g.decode(namespace.NameMapping, &g.NameMapping)
	if err != nil {
		return err
	}

	return nil
}

func (c *Configurations) remoteCfg(ns namespace.NacosNamespace) (string, error) {
	if strings.HasPrefix(ns, GroupGlobal) { // 公共配置 Global.Database.xxx
		split := strings.Split(ns, ".")
		return c.global.GetConfig(vo.ConfigParam{DataId: ns, Group: split[1]})
	}

	// 服务配置 命名空间名称.服务名称.配置(Api.gameSite.Application)
	return c.ago.GetConfig(vo.ConfigParam{DataId: ns, Group: c.app.AppID})
}

func (c *Configurations) decode(ns namespace.NacosNamespace, v interface{}) error {
	// 注册到函数里
	if !strings.HasPrefix(ns, "Global") && !strings.ContainsAny(ns, ".") {
		ns = c.app.AppID + "." + ns
	}

	content, err := c.remoteCfg(ns)
	if err != nil {
		fmt.Printf("加载配置失败！配置名称：%v \n", ns)
		return err
	}

	fmt.Println("------------------------配置读取成功------------------", ns, content)
	c.f[ns] = v
	err = mdata.Cjson.UnmarshalFromString(content, v)
	if err != nil {
		return err
	}
	c.f[ns] = v
	return err
}

func (c *Configurations) onChange(ns namespace.NacosNamespace, v interface{}) error {
	err := c.ago.ListenConfig(vo.ConfigParam{DataId: ns, Group: "DEFAULT_GROUP", OnChange: func(namespace, group, dataId, data string) {
		oldValue, _ := mdata.Cjson.MarshalToString(v)
		fmt.Println("-----------------监听配置变动--原有配置------------ group:" + group + ", dataId:" + dataId + ", content:" + oldValue)
		err := mdata.Cjson.UnmarshalFromString(data, v)
		if err != nil {
			return
		}
		fmt.Println("-----------------监听配置变动--最新配置------------ group:" + group + ", dataId:" + dataId + ", content:" + data)
	}})
	if err != nil {
		return err
	}
	return nil
}

func (c *Configurations) poller() {
	go func() {
		for ns, v := range c.f {
			oldValue := v
			err := c.onChange(ns, v)
			if err != nil {
				fmt.Printf("Dynamic update |namespace=%s\n |old=%+v\n |new=%+v\n |err=%v\n", ns, oldValue, v, err)
			}
		}
	}()
}
