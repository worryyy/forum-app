package response

import (
	"net/http"

	apperrors "campus-forum/internal/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Body = Response

type successResponder struct{}

var Success successResponder

func (successResponder) RespData(c *gin.Context, data any) {
	write(c, responseStatus(c, http.StatusOK), "", data)
}

func (successResponder) RespMessage(c *gin.Context, message string) {
	write(c, responseStatus(c, http.StatusOK), message, nil)
}

func Fail(c *gin.Context, err error) {
	fail(c, err, "")
}

func FailMessage(c *gin.Context, err error, message string) {
	fail(c, err, message)
}

func fail(c *gin.Context, err error, message string) {
	status := http.StatusInternalServerError
	responseMessage := "internal server error"

	if appErr, ok := apperrors.AsAppError(err); ok {
		status = appErr.Status
		if status == 0 {
			status = http.StatusInternalServerError
		}
		responseMessage = appErr.Message
		if responseMessage == "" {
			responseMessage = http.StatusText(status)
		}
	}

	if message != "" {
		responseMessage = message
	}

	write(c, status, responseMessage, nil)
}

func write(c *gin.Context, status int, message string, data any) {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	c.JSON(status, Response{
		Code:    status,
		Message: message,
		Data:    data,
	})
}

func responseStatus(c *gin.Context, defaultStatus int) int {
	status := c.Writer.Status()
	if status == 0 {
		return defaultStatus
	}
	return status
}
