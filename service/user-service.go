package service

import (
	"SMS/repository"

	"github.com/gin-gonic/gin"
)

//UserService interface
type UserService interface {
	StoreUser(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type userService struct {
	repository repository.UserDatabase
}

//NewUserService constructor function
func NewUserService(repository repository.UserDatabase) UserService {
	return &userService{
		repository: repository,
	}
}

func (service *userService) StoreUser(ctx *gin.Context) {
	service.repository.StoreUser(ctx)
}
func (service *userService) Login(ctx *gin.Context) {
	service.repository.Login(ctx)
}
func (service *userService) Logout(ctx *gin.Context) {
	service.repository.Logout(ctx)
}
