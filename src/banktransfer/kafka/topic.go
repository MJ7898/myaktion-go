package kafka

import (
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

const (
	Topic ="banktransfer"
)

func EnsureTransactionTopic() {
	if err := ensureTopic(Topic, 10); err != nil {
		panic(err.Error())
	}
}

func ensureTopic(topic string, numPartitions int) error {
	conn, _ := kafka.Dial("tcp", connect)
	defer conn.Close()
	controller, _ := conn.Controller()
	var leaderConn *kafka.Conn
	leaderConn, _ = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	defer leaderConn.Close()

	topicConfigs := []kafka.TopicConfig{
	{
		Topic: topic,
		NumPartitions: numPartitions,
		ReplicationFactor: 1,
	},
}
_ = leaderConn.CreateTopics(topicConfigs...)
return nil
}