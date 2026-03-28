package kfk

import (
	"context"
	"github.com/Shopify/sarama"
	"siteLetterJob/internal/glog"
	"strings"
)

type TopicConsumerHandler func(string, []byte, []byte) //消费数据处理器,需开发者实现对数据的处理

//kfkAddrs kfka地址列 多个以逗号分割
//topics 多个以逗号分割
//hander 消费数据的业务处理
func RegisterKfkConsumerTopicsListener(kfkAddress, topics, groupId, version string, handler TopicConsumerHandler) {
	glog.Infof("kfk consumer begin listen, address=%s topic=%s groupId=%s handler=%v version=%v", kfkAddress, topics, groupId, handler, version)
	if topics == "" {
		glog.Emergency("listener topic is require!! address=%s handler=%v version=%v", kfkAddress, handler, version)
		return
	}
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Errors = true // 是否等待成功和失败后的响应,只有上面的 RequiredAcks 设置不是 NoReponse 这里才有用
	saramaConfig.Version = sarama.V2_1_0_0     //设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	gc := &groupConsumer{
		topics:  strings.Split(topics, ","),
		ready:   make(chan bool),
		handler: handler,
		version: version,
	}
	if groupId == "" {
		groupId = topics + "_group1"
	}
	client, err := sarama.NewConsumerGroup(strings.Split(kfkAddress, ","), groupId, saramaConfig)
	if err != nil {
		glog.Emergency("kfk NewConsumerGroup, topics=%s, groupId=%v version=%v err=%v", topics, groupId, version, err)
		return
	}
	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(context.Background(), strings.Split(topics, ","), gc); err != nil {
				glog.Emergency("Error from consumer topics:%v version=%v error= %v", topics, version, err)
			}
			gc.ready = make(chan bool)
		}
	}()
	<-gc.ready
	glog.Infof("kfk consumer started, topics=%s handler=%v groupId=%v version=%v", topics, handler, groupId, version)
	return
}
