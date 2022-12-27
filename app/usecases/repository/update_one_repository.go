package repositoryusecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/entities"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/pkg/errors"
)

// Validate is a function to validate function input
func (c *UpdateOneRepositoryInput) Validate() error {
	const errLocation = "repositoryUsecase/UpdateOneRepository/Validate: %s"
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		ve, ok := err.(validator.ValidationErrors)
		if !ok {
			return gerrors.InternalError{
				Code:    constants.StatusCodeInvalidParameters,
				Message: constants.ErrorMessageUnableProcessParameter,
			}.Wrap(errors.Wrapf(err, errLocation, "failed to convert validation error"))
		}
		return gerrors.ParameterError{
			Code:            constants.StatusCodeInvalidParameters,
			ValidatorErrors: &ve,
		}.Wrap(errors.Wrapf(err, errLocation, "unable to process the request due to some parameter(s) are invalid"))
	}
	return nil
}

// UpdateOneRepositoryInput is an input for UpdateOneRepository
type UpdateOneRepositoryInput struct {
	ID   *string `validate:"required,uuid"`
	Name *string `validate:"omitempty"`
	URL  *string `validate:"omitempty,url,contains=//github.com"`
}

func (uc *useCase) UpdateOneRepository(ctx context.Context, input UpdateOneRepositoryInput) (*entities.Repository, error) {
	const errLocation = "repositoryUsecase/UpdateOneRepository: %s"
	if err := input.Validate(); err != nil {
		return nil, err
	}
	repository, err := uc.RepositoryRepo.FindOneRepository(nil, repositoryrepo.FindOneRepositoryInput{
		ID: input.ID,
	})

	if err != nil {
		return nil, gerrors.InternalError{
			Code:    constants.StatusCodeDatabaseError,
			Message: constants.ErrorMessageDatabaseError,
		}.Wrap(errors.Wrapf(err, errLocation, "unable to find repository"))
	}

	if repository == nil {
		return nil, gerrors.RecordNotFoundError{
			Code:    constants.StatusCodeEntryNotFound,
			Message: constants.ErrorMessageNotFound,
		}.Wrap(errors.Errorf(errLocation, "not found repository"))
	}

	if input.Name != nil {
		repository.Name = input.Name
	}
	if input.URL != nil {
		repository.URL = input.URL
	}

	err = uc.RepositoryRepo.UpdateOneRepository(nil, repositoryrepo.UpdateOneRepositoryInput{
		RepositoryEntity: repository,
	})
	if err != nil {
		return nil, err
	}

	return repository, nil
}
