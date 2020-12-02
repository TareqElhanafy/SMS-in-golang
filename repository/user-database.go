package repository

import (
	"SMS/model"
	"SMS/validator"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//UserDatabase interface that describes the user functions
type UserDatabase interface {
	StoreUser(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

//create new user function
func (db *database) StoreUser(ctx *gin.Context) {

	rules := map[string][]string{
		"name":     {"required"},
		"email":    {"required", "email"},
		"password": {"required", "minlength:6"},
		"role":     {"required"},
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
		Role:     ctx.PostForm("role"),
	}
	user.Prepare()                                                                            //Escaping and trimming the inputs
	result := db.connection.Debug().Select("Name", "Email", "Password", "Role").Create(&user) //specifing the columns which should be filled
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
	ctx.JSON(200, gin.H{
		"code":   200,
		"status": "success",
		"data":   user,
		"token":  token,
	})

}

// user login function
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
		ctx.JSON(http.StatusUnprocessableEntity, result2.Error)
		return
	}
	ctx.JSON(http.StatusFound, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Logged in !",
		"token":   token,
	})
}

func (db *database) Logout(ctx *gin.Context) {
	var token model.Token
	user := ctx.Value("user").(model.User)
	fmt.Print(user)
	result := db.connection.Debug().Where("user_id = ?", user.ID).Delete(&token)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, result.Error)
		return
	}
	ctx.JSON(http.StatusFound, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Logged out !",
	})
}
