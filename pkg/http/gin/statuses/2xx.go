package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusOK(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
	ctx.Next()
}

func StatusCreated(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	})
	ctx.Next()
}
