package router

import (
	"rizkiwhy/test-todo-list/api/handler"
	"rizkiwhy/test-todo-list/api/middleware"
	pkgTodo "rizkiwhy/test-todo-list/package/todo"

	"github.com/gin-gonic/gin"
)

func SetupTodoRoutes(r *gin.Engine, authMiddleware *middleware.AuthMiddleware, service pkgTodo.Service) {
	todoHandler := handler.NewTodoHandler(service)
	todoRouter := r.Group("/checklist")

	todoRouter.Use(authMiddleware.AuthJWT())
	{
		todoRouter.POST("/", todoHandler.CreateTodo)
	}

}
