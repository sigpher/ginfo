package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)
	id := ctx.MustGet("id").(uint)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "sucess",
		"data": gin.H{
			"id":       id,
			"username": username,
		},
	})
}
