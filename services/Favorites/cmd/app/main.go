package main

import (
	"fmt"
	"priceComp/pkg/config"
	logger2 "priceComp/pkg/logger"
	"priceComp/pkg/postgres"
	httpFavorites "priceComp/services/Favorites/internal/delivery/http"
	"priceComp/services/Favorites/internal/repository"
)

func main() {
	cfg := config.MustLoad()
	logger := logger2.SetUpLogger()
	db, err := postgres.OpenDb(cfg.PostgresDsn)
	if err != nil {
		logger.LogError(nil, err)
		return
	}

	app := httpFavorites.NewApp(repository.New(db), logger)

	err = app.Route().Run(fmt.Sprintf(":%d", cfg.Favorites))
	if err != nil {
		logger.LogError(nil, err)
		return
	}
}
