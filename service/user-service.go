package service

import (
	"SMS/repository"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	StoreUser(ctx *gin.Context)
}

type userService struct {
	repository repository.UserDatabase
}

func NewUserService(repository repository.UserDatabase) UserService {
	return &userService{
		repository: repository,
	}
}

func (service *userService) StoreUser(ctx *gin.Context) {
	service.repository.StoreUser(ctx)
}
