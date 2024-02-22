package httpUser

import (
	"github.com/gin-gonic/gin"
	errorsCFG "priceComp/pkg/errors"
	logger2 "priceComp/pkg/logger"
	"priceComp/services/UserManagment/internal/useCase"
)

type App struct {
	userManager  useCase.UserManagment
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

func NewApp(userManager useCase.UserManagment, logger *logger2.Logger) *App {
	return &App{
		userManager:  userManager,
		router:       setUpRouter(),
		logger:       logger,
		errorHandler: errorsCFG.New(logger),
	}
}

func (a *App) Route() *gin.Engine {
	a.router.POST("/users", a.registerUser)
	a.router.POST("/users/login", a.loginUser)
	a.router.GET("/users", a.requireAuth, a.fetchUser)
	a.router.GET("/users/:id", a.fetchUserId)
	a.router.Use(a.recoverPanic())
	return a.router
}
