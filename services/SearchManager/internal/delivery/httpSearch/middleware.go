package httpSearch

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (app *App) recoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Connection", "close")
				app.errorHandler.ServerErrorResponse(c, fmt.Errorf("%s", err))
			}
		}()
		c.Next()
	}
}
