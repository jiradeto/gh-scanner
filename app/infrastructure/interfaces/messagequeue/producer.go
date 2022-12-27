package messagequeue

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/jiradeto/gh-scanner/app/environments"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
)

type Producer struct {
	conn sarama.SyncProducer
}

// NewProducer init kafka producer
func (m *messageQueueClient) NewProducer() (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	brokerAddress := strings.Join([]string{environments.KafkaHost, environments.KafkaPort}, ":")
	conn, err := sarama.NewSyncProducer([]string{brokerAddress}, config)
	if err != nil {
		return nil, gerrors.NewInternalError(err)
	}

	return &Producer{
		conn: conn,
	}, nil
}

// Close close kafka producer
func (p *Producer) Close() error {
	err := p.conn.Close()
	if err != nil {
		return gerrors.NewInternalError(err)
	}
	return nil
}

// SendMessage start sending message to kafka with topic
func (p *Producer) SendMessage(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := p.conn.SendMessage(msg)
	if err != nil {
		return gerrors.NewInternalError(err)
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
