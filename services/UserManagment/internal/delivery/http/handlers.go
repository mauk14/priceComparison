package http

import (
	"github.com/gin-gonic/gin"
	"priceComp/services/UserManagment/internal/domain"
)

func (a *App) registerUser(c *gin.Context) {
	var user domain.User

	if err := c.BindJSON(&user); err != nil {

	}
}
