package repositoryusecase

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/domain/entities"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/pkg/errors"
)

// Validate is a function to validate function input
func (c *FindAllRepositoriesInput) Validate() error {
	const errLocation = "repositoryUsecase/FindAllRepositories/Validate: %s"
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

// FindAllRepositoriesInput is an input for FindAllRepositories
type FindAllRepositoriesInput struct {
	Name            *string `validate:"omitempty,max=128"`
	URL             *string `validate:"omitempty,max=128"`
	Offset          *int    `validate:"omitempty,min=1,max=100"`
	Limit           *int    `validate:"omitempty,min=0"`
	FromCreatedDate *time.Time
	ToCreatedDate   *time.Time
}

func (uc *useCase) FindAllRepositories(ctx context.Context, input FindAllRepositoriesInput) ([]*entities.Repository, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	repositories, err := uc.RepositoryRepo.FindAllRepositories(nil, repositoryrepo.FindAllRepositoriesInput{
		Name:            input.Name,
		URL:             input.URL,
		FromCreatedDate: input.FromCreatedDate,
		ToCreatedDate:   input.ToCreatedDate,
		Offset:          input.Offset,
		Limit:           input.Limit,
	})
	if err != nil {
		return nil, err
	}

	return repositories, nil
}
