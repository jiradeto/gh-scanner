package repositoryhttp

import (
	"github.com/gin-gonic/gin"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
)

// HTTPHandler ...
type HTTPHandler interface {
	// repository
	CreateOneRepository(c *gin.Context)
	FindOneRepository(c *gin.Context)
	UpdateOneRepository(c *gin.Context)
	DeleteOneRepository(c *gin.Context)
	FindAllRepositories(c *gin.Context)

	// scanner
	StartScanner(c *gin.Context)
	FindAllRepositoryScanResults(c *gin.Context)

	// scan result
	FindOneScanResult(c *gin.Context)
	FindAllScanResults(c *gin.Context)
}

type httpHandler struct {
	RepositoryUsecase repositoryusecase.UseCase
}

// New is a constructor method of HTTPHandler
func New(
	repositoryUsecase repositoryusecase.UseCase,
) HTTPHandler {
	return &httpHandler{
		RepositoryUsecase: repositoryUsecase,
	}
}
