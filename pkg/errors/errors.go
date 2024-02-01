package errorsCFG

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"priceComp/pkg/logger"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type envelope map[string]any

type ErrorHandler struct {
	logger *logger.Logger
}

func New(logger *logger.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (e *ErrorHandler) ErrorResponse(c *gin.Context, status int, message any) {
	env := envelope{"error": message}
	c.IndentedJSON(status, env)
}

func (e *ErrorHandler) ServerErrorResponse(c *gin.Context, err error) {
	e.logger.LogError(c.Request, err)
	message := "the server encountered a problem and could not process your request"
	e.ErrorResponse(c, http.StatusInternalServerError, message)
}

func (e *ErrorHandler) BadRequestResponse(c *gin.Context, err error) {
	e.ErrorResponse(c, http.StatusBadRequest, err.Error())
}

func (e *ErrorHandler) FailedValidationResponse(c *gin.Context, errors map[string]string) {
	e.ErrorResponse(c, http.StatusUnprocessableEntity, errors)
}

func (e *ErrorHandler) InvalidCredentialsResponse(c *gin.Context) {
	message := "invalid authentication credentials"
	e.ErrorResponse(c, http.StatusUnauthorized, message)
}

func (e *ErrorHandler) InvalidAuthenticationTokenResponse(c *gin.Context) {
	message := "invalid or missing authentication token"
	e.ErrorResponse(c, http.StatusUnauthorized, message)
}
