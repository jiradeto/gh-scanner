package scannerworker

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/jiradeto/gh-scanner/app/entities"
	"github.com/jiradeto/gh-scanner/app/environments"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
)

type WorkerHandler interface {
	Start(params messagequeue.StartScannerMessage)
}

type ScannerWorker struct {
	RepositoryUsecase repositoryusecase.UseCase
	KafkaTopic        string
	configs           *messagequeue.StartScannerMessage
	findings          []entities.ScanFinding
}

func RegisterScannerWorker(scanner *ScannerWorker) error {
	consumer, err := sarama.NewConsumer([]string{environments.KafkaBrokerAddress}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer consumer.Close()
	partitionConsumer, err := consumer.ConsumePartition(scanner.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()
	for {
		msg := <-partitionConsumer.Messages()
		go scanner.Start(msg)
	}
}
