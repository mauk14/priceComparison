package httpSearch

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (a *App) listProducts(c *gin.Context) {
	searchKeyword := c.Query("search")
	category := c.Query("category")
	brand := c.QueryArray("brand[]")
	sortBy := c.Query("sortBy")
	fmt.Println(c.Request.URL.String())
	fmt.Println(brand)

	limit := 9
	offset := 0
	priceFrom := 0
	priceTo := 1500000

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

	if c.Query("priceFrom") != "" {
		if newPriceFrom, err := strconv.Atoi(c.Query("priceFrom")); err == nil {
			priceFrom = newPriceFrom
		}
	}

	if c.Query("priceTo") != "" {
		if newPriceTo, err := strconv.Atoi(c.Query("priceTo")); err == nil {
			priceTo = newPriceTo
		}
	}

	products, err := a.searchManager.ListProducts(context.Background(), searchKeyword, category, sortBy, brand, limit, offset, priceFrom, priceTo)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, products)
}
