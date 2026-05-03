package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func userRouter(rg *gin.RouterGroup, repos *repository.Repos) *gin.RouterGroup {
	service := service.NewUserService(repos)
	handler := handler.NewUserHandler(service)
	//make route with url param
	user := rg.Group("/user")
	{
		user.POST("/register", handler.CreateUser)
		user.POST("/login", handler.LoginUser)
		user.GET("/profile", auth.AuthRequired(true), handler.GetProfile)
		user.POST("/send-otp", handler.SendOTP)
		user.POST("/verify-otp", handler.VerifyOTP)
	}

	role := rg.Group("/role")
	{
		role.GET("/list", handler.ListRoles)
	}

	return rg
}
