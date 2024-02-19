package main

import (
	"fmt"
	"priceComp/pkg/config"
	logger2 "priceComp/pkg/logger"
	"priceComp/pkg/postgres"
	httpReview "priceComp/services/Review/internal/delivery/http"
	"priceComp/services/Review/internal/repository"
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

	app := httpReview.NewApp(repository.New(db), logger)

	err = app.Route().Run(fmt.Sprintf(":%d", cfg.DataCollection))
	if err != nil {
		logger.LogError(nil, err)
		return
	}
}
