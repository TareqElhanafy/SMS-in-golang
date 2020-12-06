package service

import (
	"SMS/repository"

	"github.com/gin-gonic/gin"
)

//SubjectService interface which describes the subject functions
type SubjectService interface {
	Create(ctx *gin.Context)
}

type subjectService struct {
	repository repository.SubjectRepository
}

//NewSubjectService contructor function
func NewSubjectService(repository repository.SubjectRepository) SubjectService {
	return &subjectService{
		repository: repository,
	}
}

func (service *subjectService) Create(ctx *gin.Context) {
	service.repository.Create(ctx)
}
