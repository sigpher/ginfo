package middleware

import (
	"fmt"
	"ginfo/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func DoubleTokenMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		fmt.Println(ctx.Request.Header)
		if authHeader == "" {
			ctx.JSON(200, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			ctx.Abort()
			return
		}
		fmt.Println("authHeader = ", authHeader)
		parts := strings.Split(authHeader, " ")
		fmt.Println("len = ", len(parts))
		fmt.Println("parts[0] = ", parts[0])
		if !(len(parts) == 3 && parts[0] == "Bearer") {
			ctx.JSON(200, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			ctx.Abort()
			return
		}

		parseToken, shouldUpdate, err := util.ParseDoubleToken(parts[1], parts[2])
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			ctx.Abort()
			return
		}
		if shouldUpdate {
			parts[1], parts[2] = util.GenDoubleToken(parseToken.ID, parseToken.Username)
			// 如果需要刷新双Token时，返回双Token
			ctx.JSON(200, gin.H{
				"code": 200,
				"msg":  "鉴权成功",
				"data": gin.H{
					"accessToken":  parts[1],
					"refreshToken": parts[2],
				},
			})
		}
		ctx.Set("id", parseToken.ID)
		ctx.Set("username", parseToken.Username)
		ctx.Next()
	}

}
