package main

import (
	"fmt"
	"net/http"

	"rizkiwhy/test-todo-list/api/middleware"
	"rizkiwhy/test-todo-list/api/router"

	pkgTodo "rizkiwhy/test-todo-list/package/todo"
	mTodo "rizkiwhy/test-todo-list/package/todo/model"
	pkgUser "rizkiwhy/test-todo-list/package/user"
	mUser "rizkiwhy/test-todo-list/package/user/model"
	"rizkiwhy/test-todo-list/util/config"
	"rizkiwhy/test-todo-list/util/database"
	"rizkiwhy/test-todo-list/util/logger"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog/log"
)

func main() {
	g := gin.Default()

	logger.InitLogger()

	router.SetupPingRoutes(g)
	db, err := database.MySQLConnection()
	if err != nil {
		log.Fatal().Err(err).Msg("[main] Failed to connect to MySQL")
	}
	db.AutoMigrate(&mUser.User{}, &mTodo.Todo{})

	redisClient, err := database.RedisConnection()
	if err != nil {
		log.Fatal().Err(err).Msg("[main] Failed to connect to Redis")
	}

	userRepository := pkgUser.NewRepository(db)
	userCacheRepository := pkgUser.NewCacheRepository(redisClient)
	userService := pkgUser.NewService(userRepository, userCacheRepository)
	router.SetupUserRoutes(g, userService)

	authMiddleware := middleware.NewAuthMiddleware(userRepository, userCacheRepository)

	todoRepository := pkgTodo.NewRepository(db)
	todoService := pkgTodo.NewService(todoRepository)
	router.SetupTodoRoutes(g, authMiddleware, todoService)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Config("SERVICE_PORT")),
		Handler: g,
	}

	server.ListenAndServe()
}
