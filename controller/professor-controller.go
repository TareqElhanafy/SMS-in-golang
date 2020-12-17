package controller

import (
	"SMS/service"

	"github.com/gin-gonic/gin"
)

//ProfessorController inteface
type ProfessorController interface {
	StoreOrUpdateProf(ctx *gin.Context) error
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

func (controller *professorController) StoreOrUpdateProf(ctx *gin.Context) error {
	controller.service.StoreOrUpdateProf(ctx)
	return nil
}
func (controller *professorController) UpdateProf(ctx *gin.Context) error {
	controller.service.UpdateProf(ctx)
	return nil
}
