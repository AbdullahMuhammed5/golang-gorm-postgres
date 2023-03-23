package middleware

import (
	"net/http"
	"strings"

	"github.com/abdullahmuhammed5/golang-gorm-postgres/models"
	"github.com/gin-gonic/gin"
)

func OnlyAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := ctx.MustGet("currentUser").(models.User)

		if !strings.Contains(currentUser.Role, "admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not authorized!"})
			return
		}
		ctx.Next()
	}
}
