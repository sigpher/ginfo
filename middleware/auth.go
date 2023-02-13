package middleware

import (
	"ginfo/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "empty request header",
			})
			ctx.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2) && parts[0] == "Bearer" {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "wrong request header format",
			})
			ctx.Abort()
			return
		}
		mc, err := util.ParseToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "invalid token",
			})
			ctx.Abort()
			return
		}
		ctx.Set("id", mc.ID)
		ctx.Set("username", mc.Username)
		ctx.Next()
	}
}
