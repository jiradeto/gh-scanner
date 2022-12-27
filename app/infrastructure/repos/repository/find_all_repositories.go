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

// FindAllRepositoriesInput is an input for find one scan result
type FindAllRepositoriesInput struct {
	Name            *string
	URL             *string
	FromCreatedDate *time.Time
	ToCreatedDate   *time.Time
	Limit           *int
}

func (repo *repo) FindAllRepositories(tx *gorm.DB, input FindAllRepositoriesInput) ([]*entities.Repository, error) {
	const errLocation = "repositoryRepo/FindAllRepositories: %s"
	var resultModel models.Repositories
	query := repo.selectDB(tx)

	if input.FromCreatedDate != nil {
		query = query.Where("created_at >= ?", *input.FromCreatedDate)
	}

	if input.ToCreatedDate != nil {
		query = query.Where("created_at <= ?", *input.ToCreatedDate)
	}

	if input.Name != nil {
		query = query.Where("name LIKE ?", "%"+*input.Name+"%")
	}

	if input.URL != nil {
		query = query.Where("url LIKE ?", "%"+*input.URL+"%")
	}

	if input.Limit != nil {
		query = query.Limit(*input.Limit)
	}

	result := query.Order("created_at asc").Find(&resultModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gerrors.RecordNotFoundError{
				Code:    constants.StatusCodeEntryNotFound,
				Message: constants.ErrorMessageNotFound,
			}
		}
		return nil, errors.Wrap(result.Error, fmt.Sprintf(errLocation, "unable to find repositories due to database error"))
	}

	resultEntities, err := resultModel.ToEntities()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(errLocation, "unable to covert from model to entities"))
	}

	return resultEntities, nil
}
