package main

import (
	"SMS/controller"
	"SMS/middleware"
	"SMS/repository"
	"SMS/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	databaseRepository repository.DatabaseRepository = repository.NewDatabaseRepository()
	userRepository     repository.UserDatabase       = repository.NewDatabaseRepository()
	userService        service.UserService           = service.NewUserService(userRepository)
	userController     controller.UserController     = controller.NewUserController(userService)
)

func main() {

	server := gin.New()
	server.POST("/users", func(ctx *gin.Context) {
		err := userController.StoreUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	server.GET("/users", middleware.Auth(), func(ctx *gin.Context) {
		ctx.JSON(200, "hi")
	})
	server.Run(":1001")
}
