package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"priceComp/pkg/config"
	logger2 "priceComp/pkg/logger"
	"priceComp/pkg/postgres"
	"priceComp/services/UserManagment/internal/delivery/httpUser"
	"priceComp/services/UserManagment/internal/repository"
	"priceComp/services/UserManagment/internal/useCase"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	logger := logger2.SetUpLogger()
	db, err := postgres.OpenDb(cfg.PostgresDsn)
	if err != nil {
		logger.LogError(nil, err)
		return
	}
	app := httpUser.NewApp(useCase.New(repository.New(db)), logger)

	err = app.Route().Run(fmt.Sprintf(":%d", cfg.UserManager))
	if err != nil {
		logger.LogError(nil, err)
		return
	}
}
