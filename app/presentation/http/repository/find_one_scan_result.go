package repositoryhttp

import (
	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/response"
)

func (handler *httpHandler) FindOneScanResult(c *gin.Context) {
	resultID := c.Param("resultID")
	scanResult, err := handler.RepositoryUsecase.FindOneScanResult(c.Request.Context(), repositoryusecase.FindOneScanResultInput{
		ID: pointer.ToString(resultID),
	})

	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseSuccess(c, scanResult.MapResponse())
}
