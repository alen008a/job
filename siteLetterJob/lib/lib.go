package lib

import (
	"siteLetterJob/config"
	"siteLetterJob/lib/kfk"
)

func InitLib() error {

	kfk.InitProducer(config.GetConfig().Kafka.KafkaAddr)
	return nil
}
