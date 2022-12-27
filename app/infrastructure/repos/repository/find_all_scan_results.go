package repositoriyepo

import (
	"fmt"
	"time"

	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/domain/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/models"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// FindAllScanResultsInput is an input for find one scan result
type FindAllScanResultsInput struct {
	RepositoryID    *string
	FromCreatedDate *time.Time
	ToCreatedDate   *time.Time
	Offset          *int
	Limit           *int
}

func (repo *repo) FindAllScanResults(tx *gorm.DB, input FindAllScanResultsInput) ([]*entities.ScanResult, error) {
	const errLocation = "repositoryRepo/FindAllScanResults: %s"
	var resultModel models.ScanResults
	query := repo.selectDB(tx)

	if input.FromCreatedDate != nil {
		query = query.Where("created_at >= ?", *input.FromCreatedDate)
	}

	if input.ToCreatedDate != nil {
		query = query.Where("created_at <= ?", *input.ToCreatedDate)
	}

	if input.RepositoryID != nil {
		query = query.Where("repository_id = ?", *input.RepositoryID)
	}

	if input.Limit != nil {
		query = query.Limit(*input.Limit)
	}

	if input.Offset != nil {
		query = query.Offset(*input.Offset)
	}

	result := query.Order("created_at asc").Find(&resultModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gerrors.RecordNotFoundError{
				Code:    constants.StatusCodeEntryNotFound,
				Message: constants.ErrorMessageNotFound,
			}
		}
		return nil, errors.Wrap(result.Error, fmt.Sprintf(errLocation, "unable to find scan result due to database error"))
	}

	resultEntities, err := resultModel.ToEntities()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(errLocation, "unable to covert from model to entities"))
	}

	return resultEntities, nil
}
