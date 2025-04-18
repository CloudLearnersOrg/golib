package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusMultipleChoices(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusMultipleChoices, Response{
		Code:    http.StatusMultipleChoices,
		Message: message,
		Data:    data,
	})
}

func StatusMovedPermanently(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusMovedPermanently, Response{
		Code:    http.StatusMovedPermanently,
		Message: message,
		Data:    data,
	})
}

func StatusFound(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusFound, Response{
		Code:    http.StatusFound,
		Message: message,
		Data:    data,
	})
}

func StatusSeeOther(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusSeeOther, Response{
		Code:    http.StatusSeeOther,
		Message: message,
		Data:    data,
	})
}

func StatusNotModified(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusNotModified, Response{
		Code:    http.StatusNotModified,
		Message: message,
		Data:    data,
	})
}

func StatusTemporaryRedirect(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusTemporaryRedirect, Response{
		Code:    http.StatusTemporaryRedirect,
		Message: message,
		Data:    data,
	})
}

func StatusPermanentRedirect(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusPermanentRedirect, Response{
		Code:    http.StatusPermanentRedirect,
		Message: message,
		Data:    data,
	})
}
