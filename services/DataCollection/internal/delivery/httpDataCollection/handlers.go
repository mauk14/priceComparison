package httpDataCollection

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	errorsCFG "priceComp/pkg/errors"
	"priceComp/pkg/validator"
	"strconv"
	"strings"
)

func (a *App) addData(c *gin.Context) {
	var input struct {
		Url  string `json:"url"`
		Shop string `json:"shop"`
	}
	if err := c.BindJSON(&input); err != nil {
		a.errorHandler.BadRequestResponse(c, err)
		return
	}
	input.Shop = strings.ToLower(input.Shop)
	ok := false
	v := validator.New()
	shops := []string{"dns"}
	for _, s := range shops {
		if s == input.Shop {
			ok = true
			break
		}
	}
	if !ok {
		v.AddError("shop", "this shop we cannot parse")
		a.errorHandler.FailedValidationResponse(c, v.Errors)
		return
	}
	product, err := a.dataCollection.AddNewProduct(context.Background(), input.Url, input.Shop)
	if err != nil {
		switch {
		case errors.Is(err, errorsCFG.ErrBadUrl):
			v.AddError("url", "a bad url")
			a.errorHandler.FailedValidationResponse(c, v.Errors)
		default:
			a.errorHandler.ServerErrorResponse(c, err)
		}
		return
	}

	c.IndentedJSON(http.StatusCreated, product)
}

func (a *App) getData(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.ParseInt(params, 10, 64)
	v := validator.New()
	if err != nil || id < 1 {
		v.AddError("id", "invalid id parameter")
		a.errorHandler.FailedValidationResponse(c, v.Errors)
		return
	}

	product, err := a.dataCollection.GetProduct(context.Background(), id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			v.AddError("id", "invalid id parameter")
			a.errorHandler.FailedValidationResponse(c, v.Errors)
		default:
			a.errorHandler.ServerErrorResponse(c, err)
		}
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}

func (a *App) getImage(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.ParseInt(params, 10, 64)
	v := validator.New()
	if err != nil || id < 1 {
		v.AddError("id", "invalid id parameter")
		a.errorHandler.FailedValidationResponse(c, v.Errors)
		return
	}

	image, err := a.dataCollection.GetImage(context.Background(), id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			v.AddError("id", "invalid id parameter")
			a.errorHandler.FailedValidationResponse(c, v.Errors)
		default:
			a.errorHandler.ServerErrorResponse(c, err)
		}
		return
	}
	c.Data(http.StatusOK, "image/jpeg", image.Image_data)
}

func (a *App) addManyData(c *gin.Context) {
	var input struct {
		Url  string `json:"url"`
		Shop string `json:"shop"`
	}
	if err := c.BindJSON(&input); err != nil {
		a.errorHandler.BadRequestResponse(c, err)
		return
	}
	input.Shop = strings.ToLower(input.Shop)
	ok := false
	v := validator.New()
	shops := []string{"dns"}
	for _, s := range shops {
		if s == input.Shop {
			ok = true
			break
		}
	}
	if !ok {
		v.AddError("shop", "this shop we cannot parse")
		a.errorHandler.FailedValidationResponse(c, v.Errors)
		return
	}

	products, err := a.dataCollection.AddNewManyProduct(context.Background(), input.Url, input.Shop)
	if err != nil {
		switch {
		case errors.Is(err, errorsCFG.ErrBadUrl):
			v.AddError("url", "a bad url")
			a.errorHandler.FailedValidationResponse(c, v.Errors)
		default:
			a.errorHandler.ServerErrorResponse(c, err)
		}
		return
	}
	c.IndentedJSON(http.StatusCreated, products)
}

func (a *App) addShop(c *gin.Context) {
	var input struct {
		Url  string `json:"url"`
		Shop string `json:"shop"`
	}
	if err := c.BindJSON(&input); err != nil {
		a.errorHandler.BadRequestResponse(c, err)
		return
	}
	input.Shop = strings.ToLower(input.Shop)
	shop, err := a.dataCollection.AddShop(context.Background(), input.Shop, input.Url)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, shop)
}

func (a *App) addPrice(c *gin.Context) {
	var input struct {
		Shop string `json:"shop"`
	}
	if err := c.BindJSON(&input); err != nil {
		a.errorHandler.BadRequestResponse(c, err)
		return
	}
	input.Shop = strings.ToLower(input.Shop)

	v := validator.New()
	ok := false
	shops := []string{"dns"}
	for _, s := range shops {
		if s == input.Shop {
			ok = true
			break
		}
	}
	if !ok {
		v.AddError("shop", "this shop we cannot parse")
		a.errorHandler.FailedValidationResponse(c, v.Errors)
		return
	}

	prices, err := a.dataCollection.AddPrice(context.Background(), input.Shop)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, prices)
}
