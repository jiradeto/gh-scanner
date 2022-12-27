package repositoriyepo

import (
	"fmt"

	"github.com/jiradeto/gh-scanner/app/domain/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// UpdateOneScanResultInput is a DTO for updating one Repository
type UpdateOneScanResultInput struct {
	ScanResultEntity *entities.ScanResult
}

// UpdateOneScanResult is a function for updating scan result
func (repo *repo) UpdateOneScanResult(tx *gorm.DB, input UpdateOneScanResultInput) error {
	const errLocation = "repositoryRepo/UpdateOneScanResult %s"

	scanResultModel, err := new(models.ScanResult).FromEntity(input.ScanResultEntity)
	if err != nil {
		return errors.Wrapf(err, errLocation, "unable to parse an entity to model")
	}

	query := repo.selectDB(tx)
	result := query.Updates(&scanResultModel)
	if result.Error != nil {
		return errors.Wrap(result.Error, fmt.Sprintf(errLocation, "unable to update Repository aggregates due to database error"))
	}
	return nil
}
