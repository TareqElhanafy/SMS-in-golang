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
}

func (db *database) Create(ctx *gin.Context) {
	rules := map[string][]string{
		"name": {"required"},
		"file": {"file", "pdf"},
	}

	if msgs, err := validator.Validate(ctx, rules); err {
		ctx.JSON(http.StatusUnprocessableEntity, msgs)
		return
	}

	user := ctx.Value("user").(model.User)

	file, header, err := ctx.Request.FormFile("file") //getting the file and the header to extract the name
	if err != nil {
		ctx.JSON(http.StatusNoContent, "No file")
		return
	}

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
			"error":    err,
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
