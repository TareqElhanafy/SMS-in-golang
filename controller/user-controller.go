package controller

import (
	"SMS/service"

	"github.com/gin-gonic/gin"
)

//UserController interface which describes all methods for users
type UserController interface {
	StoreUser(ctx *gin.Context) error
	DeleteUser(ctx *gin.Context) error
	Login(ctx *gin.Context) error
	Logout(ctx *gin.Context) error
}

type userController struct {
	service service.UserService
}

//NewUserController constructor function
func NewUserController(service service.UserService) UserController {
	return &userController{
		service: service,
	}
}

func (controller *userController) StoreUser(ctx *gin.Context) error {
	controller.service.StoreUser(ctx)
	return nil
}
func (controller *userController) Login(ctx *gin.Context) error {
	controller.service.Login(ctx)
	return nil
}
func (controller *userController) Logout(ctx *gin.Context) error {
	controller.service.Logout(ctx)
	return nil
}
func (controller *userController) DeleteUser(ctx *gin.Context) error {
	controller.service.DeleteUser(ctx)
	return nil
}
