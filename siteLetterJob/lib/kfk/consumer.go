package kfk

// kafka consumer
import (
	"github.com/Shopify/sarama"
	"github.com/panjf2000/ants/v2"
	"siteLetterJob/config"
	"siteLetterJob/internal/glog"
	"time"
)

type groupConsumer struct {
	topics  []string
	ready   chan bool
	handler TopicConsumerHandler
	version string //使用kafak属于第几套 第一套拉单专用 第二套主战，代理， 活动使用,第三套财务使用
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *groupConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	glog.Infof("kfk consumer setup, topics=%s handler=%v", consumer.topics, consumer.handler)
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *groupConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	glog.Infof("kfk consumer cleanup,topics=%s handler=%v", consumer.topics, consumer.handler)
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *groupConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	poolNum := config.GetSiteMsgPoolNumFormKafka()
	pool, err := ants.NewPool(poolNum)
	if err != nil {
		glog.Errorf("ConsumeClaim initialization ants.coroutine sendFeePool err:%v", err)
		return err
	}
	glog.Infof("ConsumeClaim initialization ants.coroutine is successful，and topic:%v Partition:%v MemberID:%v poolNum:%v", claim.Topic(), claim.Partition(), session.MemberID(), poolNum)

	defer pool.Release()

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		strMsg := string(message.Value)
		if len(strMsg) > 500 {
			strMsg = strMsg[0:500] + "..."
		}

		now := time.Now()
		if now.Sub(message.Timestamp).Minutes() > 1 {
			glog.Infof("kafka消费延迟：message = %s, timestamp = %v, topic = %s, partition = %d, offset = %d, msg chan size = %d, free pool size = %d", strMsg,
				message.Timestamp, message.Topic, message.Partition, message.Offset, len(claim.Messages()), pool.Free())
		} else {
			glog.Infof("开始消费kafka消息：message = %s, timestamp = %v, topic = %s, partition = %d, offset = %d", strMsg, message.Timestamp,
				message.Topic, message.Partition, message.Offset)
		}

		var serverInfo []byte
		if len(message.Headers) > 0 {
			serverInfo = message.Headers[0].Value
		} else {
			glog.Infof("开始消费kafka消息 不存在服务信息 ：message = %s, timestamp = %v, topic = %s, partition = %d, offset = %d", strMsg, message.Timestamp,
				message.Topic, message.Partition, message.Offset)
		}
		doData := message.Value
		_ = pool.Submit(func() {
			consumer.handler(consumer.version, doData, serverInfo)
			glog.Infof("kafka消息消费完成：topic = %s, partition = %d, offset = %d, headers = %+v", message.Topic, message.Partition, message.Offset,
				string(serverInfo))
			session.MarkMessage(message, "")
		})
	}
	return nil
}
