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
func (c *FindOneScanResultInput) Validate() error {
	const errLocation = "repositoryUsecase/FindOneScanResult/Validate: %s"
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

// FindOneScanResult is an input for FindOneScanResult
type FindOneScanResultInput struct {
	ID *string `validate:"required,uuid"`
}

func (uc *useCase) FindOneScanResult(ctx context.Context, input FindOneScanResultInput) (*entities.ScanResult, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	scanResult, err := uc.RepositoryRepo.FindOneScanResult(nil, repositoryrepo.FindOneScanResultInput{
		ID: input.ID,
	})
	if err != nil {
		return nil, err
	}

	return scanResult, nil
}
