package httpFavorites

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"priceComp/pkg/validator"
	"priceComp/services/Favorites/internal/domain"
	"strconv"
)

func (a *App) addFavorites(c *gin.Context) {
	var input struct {
		User_id    int64 `json:"user_id"`
		Product_id int64 `json:"product_id"`
	}

	if err := c.BindJSON(&input); err != nil {
		a.errorHandler.BadRequestResponse(c, err)
		return
	}

	var favorite domain.Favorites

	favorite.User_id = input.User_id
	favorite.Product_id = input.Product_id

	err := a.favoritesRep.Insert(context.Background(), &favorite)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, favorite)
}

func (a *App) showFavorites(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.ParseInt(params, 10, 64)
	v := validator.New()
	if err != nil || id < 1 {
		v.AddError("id", "invalid id parameter")
		a.errorHandler.FailedValidationResponse(c, v.Errors)
		return
	}
	favorites, err := a.favoritesRep.Get(context.Background(), id)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, favorites)
}

func (a *App) deleteFavorites(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.ParseInt(params, 10, 64)
	v := validator.New()
	if err != nil || id < 1 {
		v.AddError("id", "invalid id parameter")
		a.errorHandler.FailedValidationResponse(c, v.Errors)
		return
	}
	err = a.favoritesRep.Delete(context.Background(), id)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}

	var output struct {
		Message string `json:"message"`
	}

	output.Message = fmt.Sprintf("Succesfully delete favorite with %d id", id)
	c.IndentedJSON(http.StatusOK, output)
}
