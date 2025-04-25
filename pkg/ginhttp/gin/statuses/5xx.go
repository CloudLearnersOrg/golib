package ginhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusInternalServerError(ctx *gin.Context, message string, err any) {
	e := toError(err)
	ctx.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
		Error:   e.Error(),
	})
	ctx.Abort()
}

func StatusBadGateway(ctx *gin.Context, message string, err any) {
	e := toError(err)
	ctx.JSON(http.StatusBadGateway, Response{
		Code:    http.StatusBadGateway,
		Message: message,
		Error:   e.Error(),
	})
	ctx.Abort()
}

func StatusServiceUnavailable(ctx *gin.Context, message string, err any) {
	e := toError(err)
	ctx.JSON(http.StatusServiceUnavailable, Response{
		Code:    http.StatusServiceUnavailable,
		Message: message,
		Error:   e.Error(),
	})
	ctx.Abort()
}

func StatusGatewayTimeout(ctx *gin.Context, message string, err any) {
	e := toError(err)
	ctx.JSON(http.StatusGatewayTimeout, Response{
		Code:    http.StatusGatewayTimeout,
		Message: message,
		Error:   e.Error(),
	})
	ctx.Abort()
}
