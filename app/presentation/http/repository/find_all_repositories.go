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

type FindAllRepositoriesRequest struct {
	Name            *string    `form:"name"`
	URL             *string    `form:"url"`
	Limit           *int       `form:"limit,default=20"`
	FromCreatedDate *time.Time `form:"from"`
	ToCreatedDate   *time.Time `form:"to"`
}

func (handler *httpHandler) FindAllRepositories(c *gin.Context) {
	errLocation := "repositoryHTTP/FindAllRepositories: %s"
	var req FindAllRepositoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResponseError(c, gerrors.ParameterError{
			Code:    constants.StatusCodeInvalidParameters,
			Message: constants.ErrorMessageParameterInvalid,
		}.Wrap(errors.Wrap(err, fmt.Sprintf(errLocation, "unable to process parameter(s)"))))
		return
	}

	repositories, err := handler.RepositoryUsecase.FindAllRepositories(c.Request.Context(), repositoryusecase.FindAllRepositoriesInput{
		Name:            req.Name,
		URL:             req.URL,
		FromCreatedDate: req.FromCreatedDate,
		ToCreatedDate:   req.ToCreatedDate,
		Limit:           req.Limit,
	})
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseSuccess(c, repositories)
}
