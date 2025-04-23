package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusBadRequest(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
		Error:   err.Error(),
	})
	ctx.Abort()
}

func StatusUnauthorized(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: message,
		Error:   err.Error(),
	})
	ctx.Abort()
}

func StatusForbidden(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: message,
		Error:   err.Error(),
	})
	ctx.Abort()
}

func StatusNotFound(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: message,
		Error:   err.Error(),
	})
	ctx.Abort()
}

func StatusRequestTimeout(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusRequestTimeout, Response{
		Code:    http.StatusRequestTimeout,
		Message: message,
		Error:   err.Error(),
	})
	ctx.Abort()
}

func StatusConflict(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusConflict, Response{
		Code:    http.StatusConflict,
		Message: message,
		Error:   err.Error(),
	})
	ctx.Abort()
}

func StatusUnprocessableEntity(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusUnprocessableEntity, Response{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Error:   err.Error(),
	})
	ctx.Abort()
}

func StatusTooManyRequests(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusTooManyRequests, Response{
		Code:    http.StatusTooManyRequests,
		Message: message,
		Error:   err.Error(),
	})
	ctx.Abort()
}
