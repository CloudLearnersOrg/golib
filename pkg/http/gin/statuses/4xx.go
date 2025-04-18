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
}

func StatusUnauthorized(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: message,
		Error:   err.Error(),
	})
}

func StatusForbidden(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: message,
		Error:   err.Error(),
	})
}

func StatusNotFound(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: message,
		Error:   err.Error(),
	})
}

func StatusRequestTimeout(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusRequestTimeout, Response{
		Code:    http.StatusRequestTimeout,
		Message: message,
		Error:   err.Error(),
	})
}

func StatusConflict(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusConflict, Response{
		Code:    http.StatusConflict,
		Message: message,
		Error:   err.Error(),
	})
}

func StatusUnprocessableEntity(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusUnprocessableEntity, Response{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Error:   err.Error(),
	})
}

func StatusTooManyRequests(ctx *gin.Context, message string, err error) {
	ctx.JSON(http.StatusTooManyRequests, Response{
		Code:    http.StatusTooManyRequests,
		Message: message,
		Error:   err.Error(),
	})
}
