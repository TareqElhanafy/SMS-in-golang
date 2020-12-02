package controller

import (
	"SMS/service"

	"github.com/gin-gonic/gin"
)

//ProfessorController inteface
type ProfessorController interface {
	StoreProf(ctx *gin.Context) error
}

type professorController struct {
	service service.ProfessorService
}

//NewProfessorController constructor function
func NewProfessorController(service service.ProfessorService) ProfessorController {
	return &professorController{
		service: service,
	}
}

func (controller *professorController) StoreProf(ctx *gin.Context) error {
	controller.service.StoreProf(ctx)
	return nil
}
