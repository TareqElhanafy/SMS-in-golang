package middleware

import (
	"SMS/model"
	"SMS/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//IsAdmin middleware to check if the user is authorized to do some actions or not
func IsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Value("user").(model.User)
		fmt.Print(user)
		db := utils.DB()
		result := db.Debug().Where(&model.User{ID: user.ID, Role: "superAdmin"}).First(&user)
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
	}

}
