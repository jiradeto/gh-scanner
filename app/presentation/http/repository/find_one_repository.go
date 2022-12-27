package repositoryhttp

import (
	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/response"
)

func (handler *httpHandler) FindOneRepository(c *gin.Context) {
	repositoryID := c.Param("repositoryID")
	repository, err := handler.RepositoryUsecase.FindOneRepository(c.Request.Context(), repositoryusecase.FindOneRepositoryInput{
		ID: pointer.ToString(repositoryID),
	})

	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseSuccess(c, repository)
}
