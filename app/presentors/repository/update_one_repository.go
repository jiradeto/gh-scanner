package repositoryhttp

import (
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	"github.com/jiradeto/gh-scanner/app/constants"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/jiradeto/gh-scanner/app/utils/response"
	"github.com/pkg/errors"
)

// UpdateOneRepositoryRequest Implementation is a struct for receiving request from HTTP
type UpdateOneRepositoryRequest struct {
	Name *string `json:"name"`
	URL  *string `json:"url"`
}

func (handler *httpHandler) UpdateOneRepository(c *gin.Context) {
	errLocation := "repositoryHTTP/UpdateOneRepository: %s"
	var body UpdateOneRepositoryRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ResponseError(c, gerrors.ParameterError{
			Code:    constants.StatusCodeInvalidParameters,
			Message: fmt.Sprintf(constants.ErrorMessageFmtInvalidFormat, "parameter(s)"),
		}.Wrap(errors.Wrapf(err, fmt.Sprintf(errLocation, "unable to process parameter(s)"))))
		return
	}

	repositoryID := c.Param("repositoryID")
	repository, err := handler.RepositoryUsecase.UpdateOneRepository(c.Request.Context(), repositoryusecase.UpdateOneRepositoryInput{
		ID:   pointer.ToString(repositoryID),
		Name: body.Name,
		URL:  body.URL,
	})

	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseSuccess(c, repository)
}
