package main

import (
	"SMS/controller"
	"SMS/middleware"
	"SMS/repository"
	"SMS/service"
	"SMS/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	userRepository      repository.UserDatabase        = repository.NewDatabaseRepository()
	professorRepository repository.ProfessorDatabase   = repository.NewDatabaseRepository()
	subjectRepository   repository.SubjectRepository   = repository.NewDatabaseRepository()
	userService         service.UserService            = service.NewUserService(userRepository)
	professorService    service.ProfessorService       = service.NewProfessorService(professorRepository)
	subjectService      service.SubjectService         = service.NewSubjectService(subjectRepository)
	professorController controller.ProfessorController = controller.NewProfessorController(professorService)
	userController      controller.UserController      = controller.NewUserController(userService)
	subjectController   controller.SubjectController   = controller.NewSubjectController(subjectService)
)

func main() {

	session := utils.ConnectAws() // connecting to AWS by creating a session to be able to use SDK's service clients
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("sess", session)
		ctx.Next()
	})
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
		profsRoutes.POST("/new", middleware.Auth(), middleware.IsAdmin(), func(ctx *gin.Context) {
			err := professorController.StoreProf(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		})
		profsRoutes.POST("/create-subject", middleware.Auth(), middleware.IsAdmin(), func(ctx *gin.Context) {
			err := subjectController.Create(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		})
		profsRoutes.GET("/my-subjects", middleware.Auth(), middleware.IsAdmin(), func(ctx *gin.Context) {
			err := subjectController.GetSubjects(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		})
		profsRoutes.DELETE("/subject/:id/delete", middleware.Auth(), middleware.IsAdmin(), func(ctx *gin.Context) {
			err := subjectController.Delete(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		})
	}

	server.Run(":1002")
}
