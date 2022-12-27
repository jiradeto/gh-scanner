package workerhandler

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/jiradeto/gh-scanner/app/environments"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/github"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
)

type WorkerHandler interface {
	ParseScannerMessage() messagequeue.StartScannerMessage
	PerformScan(params messagequeue.StartScannerMessage)
}

type Worker struct {
	githubClient      github.GithubClient
	RepositoryUsecase repositoryusecase.UseCase
	KafkaTopic        string
}

func (w *Worker) ParseScannerMessage(msg *sarama.ConsumerMessage) *messagequeue.StartScannerMessage {
	message := string(msg.Value)
	fmt.Println("incoming message", message)
	var params messagequeue.StartScannerMessage
	err := json.Unmarshal([]byte(message), &params)
	if err != nil {
		return nil
	}
	return &params
}

func RegisterWorker(w *Worker) {
	brokerAddress := strings.Join([]string{environments.KafkaHost, environments.KafkaPort}, ":")
	consumer, err := sarama.NewConsumer([]string{brokerAddress}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer consumer.Close()
	// Subscribe to the topic
	partitionConsumer, err := consumer.ConsumePartition(w.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln(err)
	}
	defer partitionConsumer.Close()
	// Consume messages
	for {
		msg := <-partitionConsumer.Messages()
		params := w.ParseScannerMessage(msg)
		go w.PerformScan(params)
		fmt.Println("! > done")
	}
}
