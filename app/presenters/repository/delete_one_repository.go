package repositoryhttp

import (
	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/response"
)

func (handler *httpHandler) DeleteOneRepository(c *gin.Context) {
	repositoryID := c.Param("repositoryID")
	err := handler.RepositoryUsecase.DeleteOneRepository(c.Request.Context(), repositoryusecase.DeleteOneRepositoryInput{
		ID: pointer.ToString(repositoryID),
	})

	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseSuccess(c, nil)
}
