package messagequeue

import (
	"encoding/json"

	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
)

type MessageQueueClient interface {
	PublishMessage(message *StartScannerMessage) error
}

type messageQueueClient struct {
	Topic string
}

type StartScannerMessage struct {
	ResultId string
	Message  string
	URL      string
}

func (m *messageQueueClient) PublishMessage(msg *StartScannerMessage) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return gerrors.InternalError{}
	}
	messages := string(messageBytes)

	producer, vErr := m.NewProducer()
	if vErr != nil {
		return vErr
	}
	defer producer.Close()

	vErr = producer.SendMessage(m.Topic, messages)
	if vErr != nil {
		return vErr
	}
	return nil

}

func New(topic string) MessageQueueClient {
	return &messageQueueClient{
		Topic: topic,
	}
}
