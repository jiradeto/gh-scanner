package repositoriyepo

import (
	"fmt"

	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/models"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// FindOneRepositoryInput is an input for find one wallet
type FindOneRepositoryInput struct {
	ID *string
}

func (repo *repo) FindOneRepository(tx *gorm.DB, input FindOneRepositoryInput) (*entities.Repository, error) {
	const errLocation = "repositoryRepo/FindOneRepository: %s"
	var resultModel models.Repository

	query := repo.selectDB(tx)

	if input.ID != nil {
		query = query.Where(`id = ?`, *input.ID)
	}

	result := query.First(&resultModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gerrors.RecordNotFoundError{
				Code:    constants.StatusCodeEntryNotFound,
				Message: constants.ErrorMessageNotFound,
			}
		}
		return nil, errors.Wrap(result.Error, fmt.Sprintf(errLocation, "unable to find Wallet due to database error"))
	}

	resultEntity, err := resultModel.ToEntity()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(errLocation, "unable to covert from model to entity"))
	}

	return resultEntity, nil
}
