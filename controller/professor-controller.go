package controller

import (
	"SMS/service"

	"github.com/gin-gonic/gin"
)

//ProfessorController inteface
type ProfessorController interface {
	StoreProf(ctx *gin.Context) error
	UpdateProf(ctx *gin.Context) error
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
func (controller *professorController) UpdateProf(ctx *gin.Context) error {
	controller.service.UpdateProf(ctx)
	return nil
}
