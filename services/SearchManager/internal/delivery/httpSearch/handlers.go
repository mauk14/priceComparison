package httpSearch

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (a *App) listProducts(c *gin.Context) {
	searchKeyword := c.Query("search")
	category := c.Query("category")
	brand := c.Query("brand")
	sortBy := c.Query("sortBy")

	limit := 9
	offset := 0

	if c.Query("limit") != "" {
		if newLimit, err := strconv.Atoi(c.Query("limit")); err == nil {
			limit = newLimit
		}
	}

	if c.Query("offset") != "" {
		if newOffset, err := strconv.Atoi(c.Query("offset")); err == nil {
			offset = newOffset
		}
	}

	products, err := a.searchManager.ListProducts(context.Background(), searchKeyword, category, brand, sortBy, limit, offset)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, products)
}
