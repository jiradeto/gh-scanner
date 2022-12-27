package repositoriyepo

import (
	"fmt"

	"github.com/jiradeto/gh-scanner/app/domain/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateOneScanResultInput is a DTO for creating one ScanResult
type CreateOneScanResultInput struct {
	ScanResultEntity *entities.ScanResult
}

// CreateOneScanResult is a function for creating one ScanResult from data model
func (repo *repo) CreateOneScanResult(tx *gorm.DB, input CreateOneScanResultInput) (*entities.ScanResult, error) {
	const errLocation = "[ScanResult repository/create one ScanResult] %s"

	scanResultModel, err := new(models.ScanResult).FromEntity(input.ScanResultEntity)
	if err != nil {
		return nil, errors.Wrapf(err, errLocation, "unable to parse an entity to model")
	}

	query := repo.selectDB(tx)
	result := query.Create(scanResultModel)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf(errLocation, "unable to create ScanResult due to database error"))
	}

	resultEntity, err := scanResultModel.ToEntity()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(errLocation, "unable to covert from model to entity"))
	}

	return resultEntity, nil
}
