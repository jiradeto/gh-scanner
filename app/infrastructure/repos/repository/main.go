package repositoriyepo

import (
	"github.com/jiradeto/gh-scanner/app/entities"
	"gorm.io/gorm"
)

// Repo ...
type Repo interface {
	// repository
	CreateOneRepository(tx *gorm.DB, input CreateOneRepositoryInput) (*entities.Repository, error)
	FindOneRepository(tx *gorm.DB, input FindOneRepositoryInput) (*entities.Repository, error)
	FindAllRepositories(tx *gorm.DB, input FindAllRepositoriesInput) ([]*entities.Repository, error)
	UpdateOneRepository(tx *gorm.DB, input UpdateOneRepositoryInput) error
	DeleteOneRepository(tx *gorm.DB, input DeleteOneRepositoryInput) error

	// scan result
	FindAllScanResults(tx *gorm.DB, input FindAllScanResultsInput) ([]*entities.ScanResult, error)
	FindOneScanResult(tx *gorm.DB, input FindOneScanResultInput) (*entities.ScanResult, error)
	CreateOneScanResult(tx *gorm.DB, input CreateOneScanResultInput) (*entities.ScanResult, error)
	UpdateOneScanResult(tx *gorm.DB, input UpdateOneScanResultInput) error
}

type repo struct {
	DB *gorm.DB
}

// New is a constructor method of Repo
func New(db *gorm.DB) Repo {
	return &repo{
		DB: db,
	}
}

func (repo *repo) selectDB(injectedDB *gorm.DB) *gorm.DB {
	if injectedDB == nil {
		return repo.DB
	}
	return injectedDB
}
