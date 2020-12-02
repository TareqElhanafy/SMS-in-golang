package repository

import (
	"SMS/model"
	"SMS/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

//ProfessorDatabase interface that descibes the professor functions
type ProfessorDatabase interface {
	StoreProf(ctx *gin.Context)
}

func (db *database) StoreProf(ctx *gin.Context) {
	rules := map[string][]string{
		"name":     {"required"},
		"email":    {"required", "email"},
		"password": {"required", "minlength:6"},
		"age":      {"required", "integer"},
		"phone":    {"required"},
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
		Role:     "professor",
	}
	result1 := db.connection.Debug().Create(&user)
	if result1.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result1.Error)
		return
	}

	token, err2 := user.GenerateToken(user.ID)
	if err2 != nil {
		panic("error to generate token")
	}
	var storingToken model.Token
	storingToken = model.Token{
		UserID: user.ID,
		Token:  token,
	}
	result2 := db.connection.Debug().Create(&storingToken)
	if result2.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result2.Error)
		return
	}
	var professor model.Professor
	professor = model.Professor{
		UserID: user.ID,
		Age:    ctx.PostForm("age"),
		Phone:  ctx.PostForm("phone"),
	}
	result3 := db.connection.Debug().Create(&professor)
	if result3.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result3.Error)
		return
	}
	db.connection.Debug().Preload("Professors").Find(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"code":   200,
		"status": "success",
		"user":   user,
		"token":  token,
	})
}
