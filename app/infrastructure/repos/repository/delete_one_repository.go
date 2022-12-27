package repositoriyepo

import (
	"fmt"

	"github.com/jiradeto/gh-scanner/app/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// DeleteOneRepositoryInput is an input for find one scan result
type DeleteOneRepositoryInput struct {
	RepositoryEntity *entities.Repository
}

// DeleteOneRepository is a function for deleting repository
func (repo *repo) DeleteOneRepository(tx *gorm.DB, input DeleteOneRepositoryInput) error {
	const errLocation = "repositoryRepo/DeleteOneRepository %s"

	repositoryModel, err := new(models.Repository).FromEntity(input.RepositoryEntity)
	if err != nil {
		return errors.Wrapf(err, errLocation, "unable to parse an entity to model")
	}

	query := repo.selectDB(tx)
	result := query.Delete(&repositoryModel)
	if result.Error != nil {
		return errors.Wrap(result.Error, fmt.Sprintf(errLocation, "unable to update Repository aggregates due to database error"))
	}
	return nil
}
