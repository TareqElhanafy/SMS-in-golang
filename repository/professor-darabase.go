package repository

import (
	"SMS/model"
	"SMS/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//ProfessorDatabase interface that descibes the professor functions
type ProfessorDatabase interface {
	StoreProf(ctx *gin.Context)
	UpdateProf(ctx *gin.Context)
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

	tx := db.connection.Begin() //Begining a transaction to avoid saving wrong related data in DB

	result1 := tx.Debug().Create(&user)
	if result1.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusUnprocessableEntity, result1.Error)
		return
	}

	var professor model.Professor
	professor = model.Professor{
		UserID: user.ID,
		Age:    ctx.PostForm("age"),
		Phone:  ctx.PostForm("phone"),
	}
	result3 := tx.Debug().Create(&professor)
	if result3.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusUnprocessableEntity, result3.Error)
		return
	}
	tx.Commit()
	db.connection.Debug().Preload("Professor").Find(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"code":   200,
		"status": "success",
		"user":   user,
	})

}

func (db *database) UpdateProf(ctx *gin.Context) {
	rules := map[string][]string{
		"age":      {"min:35"},
		"password": {"minlength:8"},
		"email":    {"email"},
	}

	if msgs, err := validator.Validate(ctx, rules); err {
		ctx.JSON(http.StatusUnprocessableEntity, msgs)
		return
	}
	UserID, err := strconv.ParseUint(ctx.Param("ID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "Bad Request",
		})
		return
	}
	var user model.User
	result := db.connection.Debug().Where("id = ?", UserID, "role = ?", "professor").First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Not user with such an id",
		})
		return
	}
	user = model.User{
		Name:  ctx.PostForm("name"),
		Email: ctx.PostForm("email"),
	}
	pass := ctx.PostForm("password")
	confirmedPassowrd := ctx.PostForm("confirm_password")
	if pass != "" {
		if confirmedPassowrd == "" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "confirm_password is required!",
			})
			return
		} else if pass != confirmedPassowrd {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Unmatched Password !",
			})
			return
		}
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		pass = string(hashedPassword)
		user.Password = pass
		if result2 := db.connection.Debug().Where("id = ?", UserID).Updates(&user); result2.Error != nil {
			ctx.JSON(http.StatusUnprocessableEntity, result2.Error)
			return
		}

	}
	var professor model.Professor
	professor = model.Professor{
		Age:   ctx.PostForm("age"),
		Phone: ctx.PostForm("phone"),
	}

	result4 := db.connection.Debug().Where("user_id = ?", UserID).Updates(&professor)
	if result4.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result4.Error)
		return
	}
	result5 := db.connection.Debug().Where("id = ?", UserID, "role = ?", "professor").First(&user)
	if result5.Error != nil {
		ctx.JSON(http.StatusNotFound, result5.Error)
		return
	}
	db.connection.Debug().Preload("Professor").Find(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"code":   200,
		"status": "success",
		"user":   user,
	})
}
