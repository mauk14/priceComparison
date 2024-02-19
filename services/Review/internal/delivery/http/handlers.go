package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"priceComp/services/Review/internal/domain"
	"strconv"
)

func (app *App) sendReview(c *gin.Context) {
	var input struct {
		Message    string `json:"message"`
		Rating     uint   `json:"rating"`
		User_id    int64  `json:"user_id"`
		Product_id int64  `json:"product_id"`
	}

	var review domain.Review
	review.Message = input.Message
	review.Rating = input.Rating
	review.User_id = input.User_id
	review.Product_id = review.Product_id

	err := app.reviewRep.InsertReview(context.Background(), &review)
	if err != nil {
		app.errorHandler.ServerErrorResponse(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, review)
}

func (app *App) showReview(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		app.errorHandler.ServerErrorResponse(c, err)
		return
	}
	review, err := app.reviewRep.GetReview(context.Background(), id)
	if err != nil {
		app.errorHandler.ServerErrorResponse(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, review)
}
