package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/environments"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/connectors"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	"github.com/jiradeto/gh-scanner/app/infrastructure/migrations"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	healthcheckhttp "github.com/jiradeto/gh-scanner/app/presentors/health_check"
	repositoryhttp "github.com/jiradeto/gh-scanner/app/presentors/repository"
	"github.com/jiradeto/gh-scanner/app/routes"
	healthcheckusecase "github.com/jiradeto/gh-scanner/app/usecases/health_check"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/jiradeto/gh-scanner/app/utils/loggers"
	"github.com/jiradeto/gh-scanner/app/utils/response"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "Start application",
		Run:   startApp,
	}
	rootCmd.Flags().Bool("production", false, "whether an app running in production")
	rootCmd.Flags().Bool("check-migration", false, "check and run migratioh if necessary")
	rootCmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Migrate database schema",
		Run:   runMigration,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "rollbackmigration",
		Short: "Rollback Last migration",
		Run:   rollbackMigration,
	})

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
}

func startApp(cmd *cobra.Command, wow []string) {
	useProductionEnv, err := cmd.Flags().GetBool("production")
	if err != nil {
		log.Fatal("Error parsing 'production' flag")
	}
	checkMigration, err := cmd.Flags().GetBool("check-migration")
	if err != nil {
		log.Fatal("Error parsing 'check-migration' flag")
	}

	environments.Init(useProductionEnv)
	if useProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	if checkMigration {
		runMigration(nil, nil)
	}

	middlewareLog := loggers.New()
	loggers.JSON.Info("Starting gh-scanner api...")

	dbConfig := connectors.DatabaseConfig{
		Host:      environments.DBHost,
		Port:      environments.DBPort,
		User:      environments.DBUser,
		Password:  environments.DBPassword,
		DB:        environments.DBName,
		DBLogMode: environments.DevMode,
	}

	fmt.Println("dbConfig", dbConfig)
	if err := dbConfig.Validate(); err != nil {
		panic(err)
	}
	connectors.InitPostgresDB(dbConfig)
	db := connectors.ConnectPostgresDB()
	// message queue
	mq := messagequeue.New(environments.KafkaTopic)

	// repos
	repositoryRepo := repositoryrepo.New(db)

	// usecases
	healthcheckUseCase := healthcheckusecase.New()
	repositoryUsecase := repositoryusecase.New(repositoryRepo, mq)

	// http
	healthcheckHTTP := healthcheckhttp.New(healthcheckUseCase)
	repositoryHTTP := repositoryhttp.New(repositoryUsecase)

	app := gin.Default()
	app.Use(middlewareLog)
	app.NoRoute(func(c *gin.Context) {
		response.ResponseError(c, gerrors.RecordNotFoundError{
			Code:    constants.StatusCodeEntryNotFound,
			Message: "endpoint not found",
		}.Wrap(errors.New("the requested endpoint is not registered")))
	})

	routes.RegisterHealthCheckRoutes(app, &routes.HTTPRoutes{
		HealthCheck: healthcheckHTTP,
	})

	routes.RegisterAPIRoutes(app, &routes.HTTPRoutes{
		Repository: repositoryHTTP,
	})

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: app,
	}

	var exit = make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// service connections
		loggers.JSON.Info(fmt.Sprintf("Listening and serving HTTP on %s\n", port))
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			loggers.JSON.Error(fmt.Sprintf("Listening and serving HTTP Error: %s\n", err))
		}
	}()

	<-exit
	loggers.JSON.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		loggers.JSON.Error(fmt.Sprint("Server forced to shutdown:", err))
	}
	loggers.JSON.Info("Server exiting")
}

func runMigration(_ *cobra.Command, _ []string) {
	environments.Init(false)
	loggers.JSON.Info("Running migration...")

	dbConfig := connectors.DatabaseConfig{
		Host:      environments.DBHost,
		Port:      environments.DBPort,
		User:      environments.DBUser,
		Password:  environments.DBPassword,
		DB:        environments.DBName,
		DBLogMode: environments.DevMode,
	}
	if err := dbConfig.Validate(); err != nil {
		panic(err)
	}
	connectors.InitPostgresDB(dbConfig)

	err := migrations.Migrate()
	if err != nil {
		loggers.JSON.Error(err.Error())
		panic(err)
	}
}

// rollback last migration
func rollbackMigration(_ *cobra.Command, _ []string) {
	environments.Init(false)
	loggers.JSON.Info("Rollback last migration...")

	err := migrations.RollbackLast()
	if err != nil {
		loggers.JSON.Error(err.Error())
		panic(err)
	}
}
