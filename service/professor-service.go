package service

import (
	"SMS/repository"

	"github.com/gin-gonic/gin"
)

//ProfessorService interface
type ProfessorService interface {
	StoreOrUpdateProf(ctx *gin.Context)
	UpdateProf(ctx *gin.Context)
}

type professorService struct {
	repository repository.ProfessorDatabase
}

//NewProfessorService constructor function
func NewProfessorService(repository repository.ProfessorDatabase) ProfessorService {
	return &professorService{
		repository: repository,
	}
}
func (service *professorService) StoreOrUpdateProf(ctx *gin.Context) {
	service.repository.StoreOrUpdateProf(ctx)
}
func (service *professorService) UpdateProf(ctx *gin.Context) {
	service.repository.UpdateProf(ctx)
}
