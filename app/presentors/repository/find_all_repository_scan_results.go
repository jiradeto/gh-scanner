package repositoryhttp

import (
	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/response"
)

func (handler *httpHandler) FindAllRepositoryScanResults(c *gin.Context) {
	repositoryID := c.Param("repositoryID")
	scanResults, err := handler.RepositoryUsecase.FindAllScanResults(c.Request.Context(),
		repositoryusecase.FindAllScanResultsInput{
			RepositoryID: pointer.ToString(repositoryID),
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
