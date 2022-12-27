package repositoryhttp

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jiradeto/gh-scanner/app/constants"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/jiradeto/gh-scanner/app/utils/response"
	"github.com/pkg/errors"
)

// CreateOneRepositoryRequest request body for CreateOneRepository
type CreateOneRepositoryRequest struct {
	Name *string `json:"name"`
	URL  *string `json:"url"`
}

func (handler *httpHandler) CreateOneRepository(c *gin.Context) {
	errLocation := "repositoryHTTP/CreateOneRepository: %s"

	var req CreateOneRepositoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, gerrors.ParameterError{
			Code:    constants.StatusCodeInvalidParameters,
			Message: fmt.Sprintf(constants.ErrorMessageFmtInvalidFormat, "parameter(s)"),
		}.Wrap(errors.Wrap(err, fmt.Sprintf(errLocation, "unable to process parameter(s)"))))
		return
	}

	repository, err := handler.RepositoryUsecase.CreateOneRepository(c.Request.Context(), repositoryusecase.CreateOneRepositoryInput{
		Name: req.Name,
		URL:  req.URL,
	})
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseSuccess(c, repository)
}
