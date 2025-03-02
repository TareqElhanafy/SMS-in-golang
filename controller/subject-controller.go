package controller

import (
	"SMS/service"

	"github.com/gin-gonic/gin"
)

//SubjectController interfacee
type SubjectController interface {
	Create(ctx *gin.Context) error
	GetSubjects(ctx *gin.Context) error
	Delete(ctx *gin.Context) error
}

type subjectController struct {
	service service.SubjectService
}

//NewSubjectController contructor function
func NewSubjectController(service service.SubjectService) SubjectController {
	return &subjectController{
		service: service,
	}
}

func (controller *subjectController) Create(ctx *gin.Context) error {
	controller.service.Create(ctx)
	return nil
}

func (controller *subjectController) GetSubjects(ctx *gin.Context) error {
	controller.service.GetSubjects(ctx)
	return nil
}
func (controller *subjectController) Delete(ctx *gin.Context) error {
	controller.service.Delete(ctx)
	return nil
}
