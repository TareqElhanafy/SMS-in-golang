package repository

import (
	"SMS/model"
	"SMS/validator"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//UserDatabase interface that describes the user functions
type UserDatabase interface {
	StoreUser(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
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
	}
	user.Prepare()                                                                    //Escaping and trimming the inputs
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
	ctx.JSON(200, gin.H{
		"code":   200,
		"status": "success",
		"data":   user,
		"token":  token,
	})

}

//delete user function
func (db *database) DeleteUser(ctx *gin.Context) {
	UserID, err := strconv.ParseUint(ctx.Param("ID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "Bad Request",
		})
		return
	}

	var user model.User
	result := db.connection.Debug().Where("id=?", UserID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":  404,
			"error": "There is no user with such an id",
		})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result.Error)
		return
	}
	result2 := db.connection.Debug().Delete(&user)
	if result2.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result2.Error)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Successfully Deleted",
	})
}

// user login function
func (db *database) Login(ctx *gin.Context) {
	rules := map[string][]string{
		"email":    {"required", "email"},
		"password": {"required"},
	}

	if msgs, err := validator.Validate(ctx, rules); err {
		ctx.JSON(http.StatusUnprocessableEntity, msgs)
		return
	}
	var user model.User
	result := db.connection.Debug().Where("email = ?", ctx.PostForm("email")).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":  404,
			"error": "There is no user with such an email",
		})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result.Error)
		return
	}

	isMatched := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ctx.PostForm("password")))
	if isMatched != nil {
		ctx.JSON(http.StatusUnauthorized, "Unauthorized, please check that your password is correct")
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

//logout function
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
