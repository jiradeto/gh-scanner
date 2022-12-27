package repositoriyepo

import (
	"fmt"

	"github.com/jiradeto/gh-scanner/app/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateOneRepositoryInput is a DTO for creating one Repository
type CreateOneRepositoryInput struct {
	RepositoryEntity *entities.Repository
}

// CreateOneRepository is a function for creating one Repository from data model
func (repo *repo) CreateOneRepository(tx *gorm.DB, input CreateOneRepositoryInput) (*entities.Repository, error) {
	const errLocation = "[Repository repository/create one Repository] %s"

	repositoryModel, err := new(models.Repository).FromEntity(input.RepositoryEntity)
	if err != nil {
		return nil, errors.Wrapf(err, errLocation, "unable to parse an entity to model")
	}

	query := repo.selectDB(tx)
	result := query.Create(repositoryModel)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf(errLocation, "unable to create Repository due to database error"))
	}

	resultEntity, err := repositoryModel.ToEntity()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(errLocation, "unable to covert from model to entity"))
	}

	return resultEntity, nil
}
