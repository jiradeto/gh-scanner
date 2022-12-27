package repositoryusecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jiradeto/gh-scanner/app/constants"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/pkg/errors"
)

// Validate is a function to validate function input
func (c *DeleteOneRepositoryInput) Validate() error {
	const errLocation = "repositoryUsecase/DeleteOneRepository/Validate: %s"
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

// DeleteOneRepositoryInput is an input for DeleteOneRepository
type DeleteOneRepositoryInput struct {
	ID *string `validate:"required"`
}

func (uc *useCase) DeleteOneRepository(ctx context.Context, input DeleteOneRepositoryInput) error {
	const errLocation = "repositoryUsecase/DeleteOneRepository: %s"
	if err := input.Validate(); err != nil {
		return err
	}
	repository, err := uc.RepositoryRepo.FindOneRepository(nil, repositoryrepo.FindOneRepositoryInput{
		ID: input.ID,
	})
	if err != nil {
		return gerrors.InternalError{
			Code:    constants.StatusCodeDatabaseError,
			Message: constants.ErrorMessageDatabaseError,
		}.Wrap(errors.Wrapf(err, errLocation, "unable to find repository"))
	}

	if repository == nil {
		return gerrors.RecordNotFoundError{
			Code:    constants.StatusCodeEntryNotFound,
			Message: constants.ErrorMessageNotFound,
		}.Wrap(errors.Errorf(errLocation, "not found repository"))
	}

	err = uc.RepositoryRepo.DeleteOneRepository(nil, repositoryrepo.DeleteOneRepositoryInput{
		RepositoryEntity: repository,
	})

	return err
}
