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
