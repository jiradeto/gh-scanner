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
func (c *FindAllScanResultsInput) Validate() error {
	const errLocation = "repositoryUsecase/FindAllScanResults/Validate: %s"
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

// FindAllScanResultsInput is an input for FindAllScanResults
type FindAllScanResultsInput struct {
	RepositoryID    *string `validate:"omitempty,uuid"`
	Limit           *int    `validate:"omitempty,min=0"`
	FromCreatedDate *time.Time
	ToCreatedDate   *time.Time
}

func (uc *useCase) FindAllScanResults(ctx context.Context, input FindAllScanResultsInput) ([]*entities.ScanResult, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	scanResults, err := uc.RepositoryRepo.FindAllScanResults(nil, repositoryrepo.FindAllScanResultsInput{
		FromCreatedDate: input.FromCreatedDate,
		ToCreatedDate:   input.ToCreatedDate,
		Limit:           input.Limit,
		RepositoryID:    input.RepositoryID,
	})
	if err != nil {
		return nil, err
	}

	return scanResults, nil
}
