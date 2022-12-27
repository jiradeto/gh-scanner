package healthcheckhttp

import (
	"github.com/gin-gonic/gin"
	healthcheckusecase "github.com/jiradeto/gh-scanner/app/usecases/health_check"
)

// A HTTPHandler for HealthCheck endpoints
type HTTPHandler interface {
	CheckLiveness(c *gin.Context)
	CheckServiceVersion(c *gin.Context)
}

type httpHandler struct {
	Usecase healthcheckusecase.UseCase
}

// New is a function to create a new HealthCheck endpoint
func New(usecase healthcheckusecase.UseCase) HTTPHandler {
	return &httpHandler{
		Usecase: usecase,
	}
}
