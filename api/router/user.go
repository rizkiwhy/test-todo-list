package router

import (
	"rizkiwhy/test-todo-list/api/handler"
	pkgUser "rizkiwhy/test-todo-list/package/user"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, service pkgUser.Service) {
	userHandler := handler.NewUserHandler(service)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
}
