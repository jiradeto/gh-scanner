package workerhandler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/jiradeto/gh-scanner/app/environments"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
)

type WorkerHandler interface {
	ParseScannerMessage() messagequeue.StartScannerMessage
	PerformScan(params messagequeue.StartScannerMessage)
}

type ScannerWorker struct {
	RepositoryUsecase repositoryusecase.UseCase
	KafkaTopic        string
}

func (w *ScannerWorker) ParseScannerMessage(msg *sarama.ConsumerMessage) *messagequeue.StartScannerMessage {
	message := string(msg.Value)
	fmt.Println("incoming message", message)
	var params messagequeue.StartScannerMessage
	err := json.Unmarshal([]byte(message), &params)
	if err != nil {
		return nil
	}
	return &params
}

func RegisterScannerWorker(w *ScannerWorker) error {
	consumer, err := sarama.NewConsumer([]string{environments.KafkaBrokerAddress}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer consumer.Close()
	partitionConsumer, err := consumer.ConsumePartition(w.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()
	for {
		msg := <-partitionConsumer.Messages()
		params := w.ParseScannerMessage(msg)
		go w.PerformScan(params)
	}
}
