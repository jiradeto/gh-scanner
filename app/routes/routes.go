package routes

import (
	"github.com/gin-gonic/gin"
	healthcheckhttp "github.com/jiradeto/gh-scanner/app/presentation/http/health_check"
	repositoryhttp "github.com/jiradeto/gh-scanner/app/presentation/http/repository"
)

// HTTPRoutes ...
type HTTPRoutes struct {
	HealthCheck healthcheckhttp.HTTPHandler
	Repository  repositoryhttp.HTTPHandler
}

// RegisterHealthCheckRoutes represents routing for healthcheck group
func RegisterHealthCheckRoutes(r *gin.Engine, httpRoutes *HTTPRoutes) {
	apiRoute := r.Group("/health")
	{
		apiRoute.GET("/check", httpRoutes.HealthCheck.CheckLiveness)
		apiRoute.GET("/version", httpRoutes.HealthCheck.CheckServiceVersion)
	}
}

// RegisterAPIRoutes represents routing for api group
func RegisterAPIRoutes(r *gin.Engine, httpRoutes *HTTPRoutes) {
	apiRoute := r.Group("/api")
	v1GroupRoute := apiRoute.Group("v1") // , middlewares.JWTAuth.MiddlewareFunc())
	{
		// repository
		{
			v1GroupRoute.POST("/repository", httpRoutes.Repository.CreateOneRepository)
			v1GroupRoute.GET("/repository/:repositoryID", httpRoutes.Repository.FindOneRepository)
			v1GroupRoute.GET("/repository/list", httpRoutes.Repository.FindAllRepositories)
			v1GroupRoute.PUT("/repository/:repositoryID", httpRoutes.Repository.UpdateOneRepository)
			v1GroupRoute.DELETE("/repository/:repositoryID", httpRoutes.Repository.DeleteOneRepository)
		}
		// scanner
		{
			v1GroupRoute.POST("/repository/:repositoryID/scan", httpRoutes.Repository.StartScanner)
			v1GroupRoute.GET("/repository/:repositoryID/scan_result/list", httpRoutes.Repository.FindAllRepositoryScanResults)
		}
		// scan result
		{
			v1GroupRoute.GET("/scan_result/list", httpRoutes.Repository.FindAllScanResults)
			v1GroupRoute.GET("/scan_result/:resultID", httpRoutes.Repository.FindOneScanResult)
		}
	}
}
