package httpDataCollection

import (
	"github.com/gin-gonic/gin"
	errorsCFG "priceComp/pkg/errors"
	logger2 "priceComp/pkg/logger"
	"priceComp/services/DataCollection/internal/useCase"
)

type App struct {
	dataCollection useCase.DataCollectionService
	router         *gin.Engine
	logger         *logger2.Logger
	errorHandler   *errorsCFG.ErrorHandler
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

func NewApp(dataCollection useCase.DataCollectionService, logger *logger2.Logger) *App {
	return &App{
		dataCollection: dataCollection,
		router:         setUpRouter(),
		logger:         logger,
		errorHandler:   errorsCFG.New(logger),
	}
}

func (a *App) Route() *gin.Engine {
	a.router.POST("/data-collection", a.addData)
	a.router.POST("/data-collection/many", a.addManyData)
	a.router.POST("/data-collection/shop", a.addShop)
	a.router.POST("/data-collection/price", a.addPrice)
	a.router.GET("/data-collection/:id", a.getData)
	a.router.GET("/data-collection/image/:id", a.getImage)
	a.router.Use(a.recoverPanic())
	return a.router
}
