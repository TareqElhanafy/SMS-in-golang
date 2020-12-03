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
	databaseRepository  repository.DatabaseRepository  = repository.NewDatabaseRepository()
	userRepository      repository.UserDatabase        = repository.NewDatabaseRepository()
	professorRepository repository.ProfessorDatabase   = repository.NewDatabaseRepository()
	userService         service.UserService            = service.NewUserService(userRepository)
	professorService    service.ProfessorService       = service.NewProfessorService(professorRepository)
	professorController controller.ProfessorController = controller.NewProfessorController(professorService)
	userController      controller.UserController      = controller.NewUserController(userService)
)

func main() {

	server := gin.New()
	usersRoutes := server.Group("/users")
	{
		usersRoutes.POST("/", func(ctx *gin.Context) {
			err := userController.StoreUser(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		})
		usersRoutes.POST("/login", func(ctx *gin.Context) {
			err := userController.Login(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		})
		usersRoutes.GET("/logout", middleware.Auth(), func(ctx *gin.Context) {
			err := userController.Logout(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		})
		usersRoutes.GET("/all", middleware.Auth(), middleware.IsAdmin(), func(ctx *gin.Context) {
			ctx.JSON(200, "hi")
		})

	}
	profsRoutes := server.Group("/profs")
	{
		profsRoutes.POST("/", middleware.Auth(), middleware.IsAdmin(), func(ctx *gin.Context) {
			err := professorController.StoreProf(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		})
	}

	server.Run(":1002")
}
