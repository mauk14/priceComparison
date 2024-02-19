package httpReview

import (
	"github.com/gin-gonic/gin"
	errorsCFG "priceComp/pkg/errors"
	logger2 "priceComp/pkg/logger"
	"priceComp/services/Review/internal/repository"
)

type App struct {
	reviewRep    repository.ReviewRepository
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

func NewApp(reviewRep repository.ReviewRepository, logger *logger2.Logger) *App {
	return &App{
		reviewRep:    reviewRep,
		router:       setUpRouter(),
		logger:       logger,
		errorHandler: errorsCFG.New(logger),
	}
}

func (a *App) Route() *gin.Engine {
	a.router.POST("/review", a.sendReview)
	a.router.GET("/review/:id", a.showReview)
	a.router.Use(a.recoverPanic())
	return a.router
}
