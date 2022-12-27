package repositoriyepo

import (
	"fmt"

	"github.com/jiradeto/gh-scanner/app/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// UpdateOneRepositoryInput is a DTO for updating one Repository
type UpdateOneRepositoryInput struct {
	RepositoryEntity *entities.Repository
}

// UpdateOneRepository is a function for updating repository
func (repo *repo) UpdateOneRepository(tx *gorm.DB, input UpdateOneRepositoryInput) error {
	const errLocation = "repositoryRepo/UpdateOneRepository %s"

	repositoryModel, err := new(models.Repository).FromEntity(input.RepositoryEntity)
	if err != nil {
		return errors.Wrapf(err, errLocation, "unable to parse an entity to model")
	}

	query := repo.selectDB(tx)
	result := query.Updates(&repositoryModel)
	if result.Error != nil {
		return errors.Wrap(result.Error, fmt.Sprintf(errLocation, "unable to update Repository aggregates due to database error"))
	}
	return nil
}
