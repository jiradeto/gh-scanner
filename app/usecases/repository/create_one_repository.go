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
func (c *CreateOneRepositoryInput) Validate() error {
	const errLocation = "repositoryUsecase/CreateOneRepositoryInput/Validate: %s"
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

// CreateOneRepositoryInput is an input for CreateOneRepository
type CreateOneRepositoryInput struct {
	Name *string `json:"name" validate:"required,min=1"`
	URL  *string `json:"url" validate:"url,contains=//github.com"`
}

func (uc *useCase) CreateOneRepository(ctx context.Context, input CreateOneRepositoryInput) (*entities.Repository, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	repository, err := uc.RepositoryRepo.CreateOneRepository(nil, repositoryrepo.CreateOneRepositoryInput{
		RepositoryEntity: &entities.Repository{
			Name: input.Name,
			URL:  input.URL,
		},
	})
	if err != nil {
		return nil, err
	}

	return repository, nil
}
