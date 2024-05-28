package httpFavorites

import (
	"github.com/gin-gonic/gin"
	errorsCFG "priceComp/pkg/errors"
	logger2 "priceComp/pkg/logger"
	"priceComp/services/Favorites/internal/repository"
)

type App struct {
	favoritesRep repository.FavoritesRep
	router       *gin.Engine
	logger       *logger2.Logger
	errorHandler *errorsCFG.ErrorHandler
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

func NewApp(reviewRep repository.FavoritesRep, logger *logger2.Logger) *App {
	return &App{
		favoritesRep: reviewRep,
		router:       setUpRouter(),
		logger:       logger,
		errorHandler: errorsCFG.New(logger),
	}
}

func (a *App) Route() *gin.Engine {
	a.router.POST("/favorites", a.addFavorites)
	a.router.GET("/favorites/:id", a.showFavorites)
	a.router.DELETE("/favorites/:id", a.deleteFavorites)
	a.router.Use(a.recoverPanic())
	return a.router
}
