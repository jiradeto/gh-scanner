package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jiradeto/gh-scanner/app/environments"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/connectors"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/loggers"
	scannerworker "github.com/jiradeto/gh-scanner/app/worker"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "Start app",
		Run:   startApp,
	}
	rootCmd.Flags().Bool("production", false, "A flag with a string value")
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
}

func startApp(cmd *cobra.Command, _ []string) {

	useProductionEnv, err := cmd.Flags().GetBool("production")
	if err != nil {
		log.Fatal("Error parsing production flag")
	}

	environments.Init(useProductionEnv)
	if useProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	loggers.JSON.Info("Starting gh-scanner worker...")
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
	db := connectors.ConnectPostgresDB()
	mq := messagequeue.New(environments.KafkaTopic)

	repositoryRepo := repositoryrepo.New(db)
	repositoryUsecase := repositoryusecase.New(repositoryRepo, mq)
	scannerWorker := &scannerworker.ScannerWorker{
		KafkaTopic:        environments.KafkaTopic,
		RepositoryUsecase: repositoryUsecase,
	}
	if err := scannerworker.RegisterScannerWorker(scannerWorker); err != nil {
		panic(err)
	}

	loggers.JSON.Info("Server exiting")
}
