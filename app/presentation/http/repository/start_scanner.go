package repositoryhttp

import (
	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/response"
)

func (handler *httpHandler) StartScanner(c *gin.Context) {
	repositoryID := c.Param("repositoryID")
	scanResult, err := handler.RepositoryUsecase.StartScanner(c.Request.Context(), repositoryusecase.StartScannerInput{
		ID: pointer.ToString(repositoryID),
	})

	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseSuccess(c, scanResult)
}
