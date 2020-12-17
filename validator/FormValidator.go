package validator

import (
	"SMS/utils"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Validate function for form validation
func Validate(ctx *gin.Context, rules map[string][]string) (map[string][]string, bool) {

	msgs := make(map[string][]string)

	for field, fieldRules := range rules {
		for _, rule := range fieldRules {
			parts := strings.Split(rule, ":")
			switch parts[0] {
			case "required":
				_, singleInput := ctx.GetPostForm(field)         //returns value and bool(true or false) for the every field like username ,password,email
				_, multipleInputs := ctx.GetPostFormArray(field) //returns value and bool(true or false) for every field if it takes array
				if !singleInput && !multipleInputs {
					msgs[field] = append(msgs[field], field+" is required")
				}

			case "email":
				var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
				if len(ctx.PostForm(field)) > 254 || !rxEmail.MatchString(ctx.PostForm(field)) {
					msgs[field] = append(msgs[field], field+" must be valid")
				}
			case "file":
				_, _, err := ctx.Request.FormFile("file") //takes the name of the field of files if exists
				if err != nil {
					msgs[field] = append(msgs[field], "there is no file")
				}
			case "minlength":
				number, _ := strconv.Atoi(parts[1])
				fieldNumber := len(ctx.PostForm(field))
				if ctx.PostForm(field) != "" && fieldNumber < number {
					msgs[field] = append(msgs[field], field+" must be at least "+strconv.Itoa(number)+" characters")

				}
			case "min":
				fieldNum, _ := strconv.Atoi(parts[1])
				num, _ := strconv.Atoi(ctx.PostForm(field))
				if ctx.PostForm(field) != "" && num < fieldNum {
					msgs[field] = append(msgs[field], field+"must be at least "+strconv.Itoa(fieldNum))
				}

			case "pdf":
				_, header, err := ctx.Request.FormFile(field)
				if err != nil {
					msgs[field] = append(msgs[field], "there is no file")
				} else if !utils.IsPdf(header) {
					msgs[field] = append(msgs[field], "only pdf files are allowed")
				}
			}

		}

	}
	return msgs, len(msgs) > 0

}
