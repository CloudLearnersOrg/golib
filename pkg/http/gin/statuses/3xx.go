package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusTemporaryRedirect(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusTemporaryRedirect, location)
}

func StatusPermanentRedirect(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusPermanentRedirect, location)
}

func StatusFound(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusFound, location)
}

func StatusMovedPermanently(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusMovedPermanently, location)
}
