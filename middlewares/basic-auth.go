package middlewares

import (
	"github.com/gin-gonic/gin"
)

func BasicAuth(accounts gin.Accounts) gin.HandlerFunc {
	return gin.BasicAuth(accounts)
}
