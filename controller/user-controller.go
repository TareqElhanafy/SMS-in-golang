package controller

import (
	"SMS/service"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	StoreUser(ctx *gin.Context) error
}

type userController struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return &userController{
		service: service,
	}
}

func (controller *userController) StoreUser(ctx *gin.Context) error {
	controller.service.StoreUser(ctx)
	return nil
}
