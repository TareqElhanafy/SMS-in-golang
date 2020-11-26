package middleware

import (
	"SMS/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Auth function to determine whether the user is authenticated or not
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := ctx.GetHeader("Authorization")
		tokenString := authHeader[len(BearerSchema):]

		user, valid := utils.ValidateUser(tokenString)
		if valid {
			ctx.Set("user", user)
		} else {
			ctx.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		ctx.Next()
	}

}
