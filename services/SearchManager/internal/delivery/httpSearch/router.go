package httpSearch

import (
	"github.com/gin-gonic/gin"
	errorsCFG "priceComp/pkg/errors"
	logger2 "priceComp/pkg/logger"
	"priceComp/services/SearchManager/internal/usecase"
)

type App struct {
	searchManager usecase.SearchManager
	router        *gin.Engine
	logger        *logger2.Logger
	errorHandler  *errorsCFG.ErrorHandler
}

func setUpRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}

func NewApp(searchManager usecase.SearchManager, logger *logger2.Logger) *App {
	return &App{
		searchManager: searchManager,
		router:        setUpRouter(),
		logger:        logger,
		errorHandler:  errorsCFG.New(logger),
	}
}

func (a *App) Route() *gin.Engine {
	a.router.GET("/search", a.listProducts)
	a.router.Use(a.recoverPanic())
	return a.router
}
