package repositoryhttp

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jiradeto/gh-scanner/app/constants"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/jiradeto/gh-scanner/app/utils/response"
	"github.com/pkg/errors"
)

type FindAllScanResultsRequest struct {
	Limit           *int       `form:"limit,default=20"`
	Offset          *int       `form:"start,default=0"`
	FromCreatedDate *time.Time `form:"from"`
	ToCreatedDate   *time.Time `form:"to"`
	RepositoryID    *string    `form:"repositoryID"`
}

func (handler *httpHandler) FindAllScanResults(c *gin.Context) {
	errLocation := "repositoryHTTP/FindAllScanResults: %s"
	var req FindAllScanResultsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		err := gerrors.ParameterError{
			Code:    constants.StatusCodeInvalidParameters,
			Message: constants.ErrorMessageParameterInvalid,
		}.Wrap(errors.Wrap(err, fmt.Sprintf(errLocation, "unable to process parameter(s)")))
		response.ResponseError(c, err)
		return
	}

	scanResults, err := handler.RepositoryUsecase.FindAllScanResults(c.Request.Context(), repositoryusecase.FindAllScanResultsInput{
		FromCreatedDate: req.FromCreatedDate,
		ToCreatedDate:   req.ToCreatedDate,
		Limit:           req.Limit,
		Offset:          req.Offset,
		RepositoryID:    req.RepositoryID,
	})

	if err != nil {
		response.ResponseError(c, err)
		return
	}

	var scanResultsResponse []map[string]interface{}
	for _, scanResult := range scanResults {
		scanResultsResponse = append(scanResultsResponse, scanResult.MapResponse())
	}
	response.ResponseSuccess(c, scanResultsResponse)
}
