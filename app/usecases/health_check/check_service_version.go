package healthcheckusecase

import (
	"context"
	"syscall"

	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/pkg/errors"
)

// CheckServiceVersionResponse ..
type CheckServiceVersionResponse struct {
	AppServiceVersion string `json:"app_version"`
}

// CheckServiceVersion ..
func (uc *useCase) CheckServiceVersion(context.Context) (*CheckServiceVersionResponse, error) {
	const errLocation = "healthcheckUseCase/CheckServiceVersion: %s"
	appVersionNo, found := syscall.Getenv("APP_VERSION_NO")
	if !found {
		return nil, gerrors.RecordNotFoundError{
			Code:    constants.StatusCodeEntryNotFound,
			Message: constants.ErrorMessageNotFound,
		}.Wrap(errors.Errorf(errLocation, "not found app version number"))
	}

	checkServiceVersion := &CheckServiceVersionResponse{
		AppServiceVersion: appVersionNo,
	}
	return checkServiceVersion, nil
}
