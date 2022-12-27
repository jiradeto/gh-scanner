package repositoryusecase

import (
	"context"

	"github.com/jiradeto/gh-scanner/app/domain/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
)

// UseCase is a declarative interface for repository usecase
type UseCase interface {
	// repository
	CreateOneRepository(ctx context.Context, input CreateOneRepositoryInput) (*entities.Repository, error)
	FindOneRepository(ctx context.Context, input FindOneRepositoryInput) (*entities.Repository, error)
	UpdateOneRepository(ctx context.Context, input UpdateOneRepositoryInput) (*entities.Repository, error)
	DeleteOneRepository(ctx context.Context, input DeleteOneRepositoryInput) error
	FindAllRepositories(ctx context.Context, input FindAllRepositoriesInput) ([]*entities.Repository, error)

	// scanner
	StartScanner(ctx context.Context, input StartScannerInput) (*entities.ScanResult, error)

	// scan result
	UpdateOneScanResult(ctx context.Context, input UpdateOneScanResultInput) (*entities.ScanResult, error)
	FindOneScanResult(ctx context.Context, input FindOneScanResultInput) (*entities.ScanResult, error)
	FindAllScanResults(ctx context.Context, input FindAllScanResultsInput) ([]*entities.ScanResult, error)
}

type useCase struct {
	RepositoryRepo repositoryrepo.Repo
	MessageQueue   messagequeue.MessageQueueClient
}

// New is a constructor method of UseCase
func New(
	repositoryRepo repositoryrepo.Repo,
	messageQueue messagequeue.MessageQueueClient,
) UseCase {
	return &useCase{
		RepositoryRepo: repositoryRepo,
		MessageQueue:   messageQueue,
	}
}
