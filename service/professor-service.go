package service

import (
	"SMS/repository"

	"github.com/gin-gonic/gin"
)

//ProfessorService interface
type ProfessorService interface {
	StoreProf(ctx *gin.Context)
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
func (service *professorService) StoreProf(ctx *gin.Context) {
	service.repository.StoreProf(ctx)
}
