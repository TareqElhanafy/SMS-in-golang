package repository

import (
	"SMS/model"
	"SMS/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//UserDatabase interface that describe
type UserDatabase interface {
	StoreUser(ctx *gin.Context)
	Login(ctx *gin.Context)
}

//create new user
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
	var tokens []model.Token
	user = model.User{
		Name:     ctx.PostForm("name"),
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
	}
	user.Prepare()
	result := db.connection.Debug().Select("Name", "Email", "Password").Create(&user) //specifing the columns which should be filled
	if result.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result.Error)
		return
	}
	token, err := user.GenerateToken(user.ID)
	if err != nil {
		panic("error to generate token")
	}
	storingToken := model.Token{
		UserID: user.ID,
		Token:  token,
	}
	result2 := db.connection.Debug().Create(&storingToken)
	if result2.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result2.Error)
		return
	}
	db.connection.Debug().Model(&user).Where("user_id = ?", user.ID).Association("Tokens").Find(&tokens)

	ctx.JSON(200, gin.H{
		"code":   200,
		"status": "success",
		"data":   user,
		"token":  token,
	})

}

// user login
func (db *database) Login(ctx *gin.Context) {
	rules := map[string][]string{
		"email":    {"required", "email"},
		"password": {"required", "minlength:6"},
	}

	if msgs, err := validator.Validate(ctx, rules); err {
		ctx.JSON(http.StatusUnprocessableEntity, msgs)
		return
	}
	var user model.User
	result := db.connection.Debug().Where("email = ?", ctx.PostForm("email")).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, "Unauthorized, please check your email and password are correct")
		return
	}

	isMatched := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ctx.PostForm("password")))
	if isMatched != nil {
		ctx.JSON(http.StatusUnauthorized, "Unauthorized, please check your email and password are correct")
		return
	}
	token, err := user.GenerateToken(user.ID)
	if err != nil {
		panic("error to generate token")
	}
	storingToken := &model.Token{
		UserID: user.ID,
		Token:  token,
	}
	result2 := db.connection.Debug().Create(&storingToken)
	if result2.Error != nil {
		ctx.JSON(http.StatusNotFound, result2.Error)
		return
	}
	ctx.JSON(http.StatusFound, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Logged in !",
		"token":   token,
	})
}
