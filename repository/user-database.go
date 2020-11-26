package repository

import (
	"SMS/model"
	"SMS/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

//UserDatabase interface that describe
type UserDatabase interface {
	StoreUser(ctx *gin.Context)
}

func (db *database) StoreUser(ctx *gin.Context) {

	rules := map[string][]string{
		"name":     {"required"},
		"email":    {"required", "email"},
		"password": {"required", "minlength:6"},
	}

	if msgs, err := validator.Validate(ctx, rules); err {
		ctx.JSON(http.StatusUnprocessableEntity, msgs)
		return
	}

	var user model.User
	user = model.User{
		Name:     ctx.PostForm("name"),
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
	}
	user.Prepare()
	result := db.connection.Create(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result.Error)
		return
	}
	token, err := user.GenerateToken(user.ID)
	if err != nil {
		panic("error to generate token")
	}
	user.Tokens = append(user.Tokens, token)
	db.connection.Save(&user)
	ctx.JSON(200, gin.H{
		"code":   200,
		"status": "success",
		"data":   user,
	})

}
