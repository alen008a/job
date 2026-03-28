package kfk

import (
	"siteLetterJob/config"
	"siteLetterJob/lib/httpclient"
	"strings"
	"time"

	"siteLetterJob/internal/glog"
	"siteLetterJob/mdata"

	"github.com/Shopify/sarama"
)

var asyncProducer sarama.AsyncProducer

// 初始化 kafka 的生产者
func InitProducer(kafkaAddr string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 等待服务器所有副本都保存成功后才响应
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机的分区类型
	config.Producer.Return.Successes = true                   // 是否等待成功和失败后的响应,只有上面的 RequiredAcks 设置不是 NoReponse 这里才有用
	config.Producer.Return.Errors = true
	config.Version = sarama.V2_1_0_0 //设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置

	// 异步生产者
	go AsyncProducer(config, kafkaAddr)
}

// AsyncProducer 异步 kafka 的生产者
func AsyncProducer(config *sarama.Config, kafkaAddr string) {
	addrArr := strings.Split(kafkaAddr, ",")
	var err error

	// 使用配置，新建一个异步生产者
	asyncProducer, err = sarama.NewAsyncProducer(addrArr, config)
	if err != nil {
		glog.Emergency("kafka异常 -- err=%v -- addr=%v -- config=%v", err, addrArr, mdata.MustMarshal2String(config))
		return
	}

	for {
		// 判断发送是否成功
		select {
		case <-asyncProducer.Successes():
		case fail := <-asyncProducer.Errors():
			glog.Emergency("msgPush producer err: %v", fail.Err)
		}
	}
}

// MsgPushKafka 消息推送kafka中
func MsgPushKafka(data []byte, topic string, key ...string) {
	var k string
	if len(key) != 1 {
		k = topic
	} else {
		k = key[0]
	}

	if asyncProducer == nil {
		glog.Error(">>> kafka警告 -- MsgPushKafka -- asyncProducer 为空 走不下去了")
		return
	}

	//  发送的消息， 主题， key
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(k),
		Value: sarama.ByteEncoder(data),
	}

	glog.Infof(">>> MsgPushKafka -- kafka推送信息 -- Topic=%v -- data=%v", topic, string(data))
	asyncProducer.Input() <- msg
}

func SendSlack(siteId, code, msg string) {
	resp, err := SendBotNotify(siteId, code, msg, 3)
	glog.Emergency(">>> SendSlack --发送告警信息  --serviceCode=%v --msg=%v --result=%v --error=%v", code, msg, resp, err)
}

// 发送告警
func SendBotNotify(siteId, serviceCode, msg string, botType int) (string, error) {
	notiryUrl := config.GetConfig().VerifyCodeDomain + "/verifycode/bot/v1/send"
	paramsStr := mdata.MustMarshal2Byte(map[string]interface{}{"serviceCode": serviceCode, "botType": botType, "msg": msg})
	byteNotifys, err := httpclient.POSTJson(notiryUrl, paramsStr,
		map[string]string{"Content-Type": "application/json", mdata.HeaderSite: siteId}, httpclient.GetShortProxyNotifyClient(time.Second*10))
	return string(byteNotifys), err
}
