package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusInternalServerError(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
		Error:   err.Error(),
	})
}

func StatusBadGateway(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusBadGateway, Response{
		Code:    http.StatusBadGateway,
		Message: message,
		Error:   err.Error(),
	})
}

func StatusServiceUnavailable(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusServiceUnavailable, Response{
		Code:    http.StatusServiceUnavailable,
		Message: message,
		Error:   err.Error(),
	})
}

func StatusGatewayTimeout(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusGatewayTimeout, Response{
		Code:    http.StatusGatewayTimeout,
		Message: message,
		Error:   err.Error(),
	})
}
