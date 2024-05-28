package httpUser

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	errorsCFG "priceComp/pkg/errors"
	"priceComp/pkg/validator"
	"priceComp/services/UserManagment/internal/domain"
	"strconv"
)

func (a *App) registerUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		a.errorHandler.BadRequestResponse(c, err)
		return
	}

	var user domain.User
	user.Name = input.Name
	user.Email = input.Email
	err := user.Password.Set(input.Password)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}
	v := validator.New()
	if domain.ValidateUser(v, &user); !v.Valid() {
		a.errorHandler.FailedValidationResponse(c, v.Errors)
		return
	}

	err = a.userManager.RegisterUser(context.Background(), &user)
	if err != nil {
		switch {
		case errors.Is(err, errorsCFG.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			a.errorHandler.FailedValidationResponse(c, v.Errors)
		default:
			a.errorHandler.ServerErrorResponse(c, err)
		}
		return
	}

	c.IndentedJSON(http.StatusCreated, user)

}

func (a *App) loginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		a.errorHandler.BadRequestResponse(c, err)
		return
	}

	v := validator.New()

	domain.ValidateEmail(v, input.Email)
	domain.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		a.errorHandler.FailedValidationResponse(c, v.Errors)
		return
	}

	token, err := a.userManager.LoginUser(context.Background(), input.Email, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, errorsCFG.ErrRecordNotFound):
			a.errorHandler.InvalidCredentialsResponse(c)
			return
		default:
			a.errorHandler.ServerErrorResponse(c, err)
			return
		}
	}

	if token == "" {
		a.errorHandler.InvalidCredentialsResponse(c)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24*30, "", "", false, true)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

func (a *App) fetchUser(c *gin.Context) {
	userJson, ok := c.Get("user")
	if !ok {
		a.errorHandler.InvalidAuthenticationTokenResponse(c)
		return
	}
	user := userJson.(*domain.User)
	c.IndentedJSON(http.StatusCreated, user)
}

func (a *App) fetchUserId(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}
	user, err := a.userManager.UserInfoById(context.Background(), id)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, user)
}

func (a *App) ChangePass(c *gin.Context) {
	var input struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BindJSON(&input); err != nil {
		a.errorHandler.BadRequestResponse(c, err)
		return
	}
	v := validator.New()

	err := a.userManager.ChangePassword(context.Background(), input.Email, input.OldPassword, input.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, errorsCFG.ErrNotCorrectPassword):
			v.AddError("old_password", "Old password is not correct")
			a.errorHandler.FailedValidationResponse(c, v.Errors)
		default:
			a.errorHandler.ServerErrorResponse(c, err)
		}
		return
	}

	var output struct {
		Message string `json:"message"`
	}
	output.Message = "Password success changed"
	c.IndentedJSON(http.StatusOK, output)
}

func (a *App) ChangeInfo(c *gin.Context) {
	var input struct {
		EmailOld  string `json:"emailOld"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		Activated bool   `json:"activated"`
	}

	if err := c.BindJSON(&input); err != nil {
		a.errorHandler.BadRequestResponse(c, err)
		return
	}

	token, err := a.userManager.ChangePersonalInfo(context.Background(), input.EmailOld, input.Email, input.Name, input.Activated)
	if err != nil {
		a.errorHandler.ServerErrorResponse(c, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24*30, "", "", false, true)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"token": token,
	})
}
