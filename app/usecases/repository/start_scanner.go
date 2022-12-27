package repositoryusecase

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/domain/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/pkg/errors"
)

// Validate is a function to validate function input
func (c *StartScannerInput) Validate() error {
	const errLocation = "repositoryUsecase/StartScanner/Validate: %s"
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

// StartScannerInput is an input for StartScanner
type StartScannerInput struct {
	ID *string `validate:"required"`
}

func (uc *useCase) StartScanner(ctx context.Context, input StartScannerInput) (*entities.ScanResult, error) {
	const errLocation = "repositoryUsecase/StartScanner: %s"
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

	now := time.Now()
	scanResult, err := uc.RepositoryRepo.CreateOneScanResult(nil, repositoryrepo.CreateOneScanResultInput{
		ScanResultEntity: &entities.ScanResult{
			RepositoryID: repository.ID,
			Status:       entities.ScanResultStatusQueued,
			QueuedAt:     &now,
		},
	})
	if err != nil {
		return nil, gerrors.InternalError{
			Code:    constants.StatusCodeDatabaseError,
			Message: constants.ErrorMessageDatabaseError,
		}.Wrap(errors.Wrapf(err, errLocation, "unable to create scan result"))
	}

	msg := &messagequeue.StartScannerMessage{
		ResultId: *scanResult.ID,
		URL:      *repository.URL,
	}
	fmt.Println("msg", msg)
	err = uc.MessageQueue.PublishMessage(msg)
	if err != nil {
		fmt.Println("Mesasgeq queue error", err.Error())
	}

	return scanResult, nil
}
