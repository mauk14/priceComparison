package httpUser

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
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

func (app *App) requireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		app.errorHandler.InvalidAuthenticationTokenResponse(c)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header)
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		app.errorHandler.InvalidAuthenticationTokenResponse(c)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			app.errorHandler.InvalidAuthenticationTokenResponse(c)
			return
		}

		user, err := app.userManager.UserInfo(context.Background(), claims["sub"].(string))
		if err != nil {
			app.errorHandler.InvalidAuthenticationTokenResponse(c)
			return
		}
		c.Set("user", user)

		c.Next()
	}
}
