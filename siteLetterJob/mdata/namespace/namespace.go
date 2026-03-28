package namespace

type NacosNamespace = string

const (
	Application NacosNamespace = "Application"
	Logger                     = "Logger"
	Kafka                      = "Global.MQ.Kafka"
	RedisCore                  = "Global.Redis.RedisCore"
	KafkaTopic                 = "Global.Config.KafkaTopic"
	Common                     = "Global.Config.Common"
	Site                       = "Global.Database.Site"
	SiteSlave                  = "Global.Database.SiteSlave"
	EdgeDB                     = "Global.Database.EdgeDB"
)
