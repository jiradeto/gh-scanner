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
func (c *UpdateOneScanResultInput) Validate() error {
	const errLocation = "repositoryUsecase/UpdateOneScanResult/Validate: %s"
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

// UpdateOneScanResultInput is an input for UpdateOneScanResult
type UpdateOneScanResultInput struct {
	ID         *string `validate:"required,uuid"`
	Status     *string `validate:"required,min=1,oneof='queued' 'in_progress' 'success' 'failure'"`
	Findings   []entities.ScanFinding
	QueuedAt   *time.Time
	ScanningAt *time.Time
	FinishedAt *time.Time
}

func (uc *useCase) UpdateOneScanResult(ctx context.Context, input UpdateOneScanResultInput) (*entities.ScanResult, error) {
	const errLocation = "repositoryUsecase/UpdateOneScanResult: %s"
	if err := input.Validate(); err != nil {
		return nil, err
	}
	scanResult, err := uc.RepositoryRepo.FindOneScanResult(nil, repositoryrepo.FindOneScanResultInput{
		ID: input.ID,
	})

	if err != nil {
		return nil, gerrors.InternalError{
			Code:    constants.StatusCodeDatabaseError,
			Message: constants.ErrorMessageDatabaseError,
		}.Wrap(errors.Wrapf(err, errLocation, "unable to find scan result"))
	}

	if scanResult == nil {
		return nil, gerrors.RecordNotFoundError{
			Code:    constants.StatusCodeEntryNotFound,
			Message: constants.ErrorMessageNotFound,
		}.Wrap(errors.Errorf(errLocation, "not found scan result"))
	}

	if input.Status != nil {
		scanResult.Status = new(entities.ScanResultStatus).Parse(*input.Status)
	}

	if input.QueuedAt != nil {
		scanResult.QueuedAt = input.QueuedAt
	}
	if input.ScanningAt != nil {
		scanResult.ScanningAt = input.ScanningAt
	}
	if input.FinishedAt != nil {
		scanResult.FinishedAt = input.FinishedAt
	}
	if input.Findings != nil {
		scanResult.Findings = input.Findings
	}

	err = uc.RepositoryRepo.UpdateOneScanResult(nil, repositoryrepo.UpdateOneScanResultInput{
		ScanResultEntity: scanResult,
	})
	if err != nil {
		return nil, err
	}

	return scanResult, nil
}
