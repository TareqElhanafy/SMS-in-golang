package repository

import (
	"SMS/model"
	"SMS/utils"
	"SMS/validator"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//SubjectRepository interface describes subject functions
type SubjectRepository interface {
	Create(ctx *gin.Context)
	GetSubjects(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

//Create subjects
func (db *database) Create(ctx *gin.Context) {
	rules := map[string][]string{
		"name": {"required"},
		"file": {"pdf"},
	}

	if msgs, err := validator.Validate(ctx, rules); err {
		ctx.JSON(http.StatusUnprocessableEntity, msgs)
		return
	}

	user := ctx.Value("user").(model.User)

	file, header, _ := ctx.Request.FormFile("file") //getting the file and the header to extract the name

	fileName := header.Filename //getting the file name
	parts := strings.Split(fileName, ".")
	extension := parts[1]
	fileName = strconv.Itoa(rand.Int()) + "." + extension

	var subject model.Subject
	subject = model.Subject{
		Name:     ctx.PostForm("name"),
		Material: fileName,
		UserID:   user.ID,
	}
	subject.Prepare()
	result := db.connection.Debug().Create(&subject)
	if result.Error != nil {
		ctx.JSON(http.StatusUnprocessableEntity, result.Error)
	}
	up, err := utils.UploadToS3(ctx, fileName, file) //uploading to S3
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":    "failed to upload to S3",
			"uploader": up,
		})
		return
	}
	db.connection.Debug().Preload("User").Find(&subject)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"subject": subject,
	})

}

//get all subjects for a specific professor
func (db *database) GetSubjects(ctx *gin.Context) {
	user := ctx.Value("user").(model.User)
	var subjects []model.Subject
	result := db.connection.Debug().Where("user_id = ?", user.ID).Find(&subjects)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, result.Error)
		return
	}
	if len(subjects) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    401,
			"status":  "Not Found",
			"message": "You do not have any subjects yet",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":     200,
		"status":   "success",
		"subjects": subjects,
	})
}

func (db *database) Delete(ctx *gin.Context) {
	var subject model.Subject
	SubjectID := ctx.Param("id")
	result := db.connection.Debug().Where("id = ?", SubjectID).First(&subject)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"status":  "Not Found",
			"message": "There is no subject with such an id",
		})
		return
	}
	result2 := db.connection.Debug().Delete(&subject)
	if result2.Error != nil {
		ctx.JSON(http.StatusInternalServerError, result.Error)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":            200,
		"status":          "success",
		"message":         "successfully Deleted",
		"deleted_subject": subject,
	})

}
