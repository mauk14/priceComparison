package main

import (
	"fmt"
	"priceComp/pkg/config"
	logger2 "priceComp/pkg/logger"
	"priceComp/pkg/postgres"
	"priceComp/services/SearchManager/internal/delivery/httpSearch"
	"priceComp/services/SearchManager/internal/repository"
	"priceComp/services/SearchManager/internal/usecase"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg.MongoDb_uri)
	logger := logger2.SetUpLogger()
	db, err := postgres.OpenDb(cfg.PostgresDsn)
	if err != nil {
		logger.LogError(nil, err)
		return
	}

	app := httpSearch.NewApp(usecase.New(repository.New(db)), logger)

	err = app.Route().Run(fmt.Sprintf(":%d", cfg.SearchManager))
	if err != nil {
		logger.LogError(nil, err)
		return
	}
}
