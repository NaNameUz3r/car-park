package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AlertMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		fmt.Println("hello-alert")
	}
}
