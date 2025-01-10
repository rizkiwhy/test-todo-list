package handler

import (
	"net/http"
	"rizkiwhy/test-todo-list/api/presenter"
	pkgTodo "rizkiwhy/test-todo-list/package/todo"
	"rizkiwhy/test-todo-list/package/todo/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type TodoHandler struct {
	Service pkgTodo.Service
}

func NewTodoHandler(service pkgTodo.Service) *TodoHandler {
	return &TodoHandler{
		Service: service,
	}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var request model.CreateTodoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("[TodoHandler][CreateTodo] Failed to bind json")
		c.JSON(http.StatusBadRequest, presenter.FailureResponse(presenter.CreateTodoInvalidRequestMessage, err.Error()))
		return
	}

	value, exists := c.Get("user_id")
	if !exists || value.(int64) == 0 {
		log.Error().Msg("[TodoHandler][CreateTodo] Failed to get user id")
		c.JSON(http.StatusBadRequest, presenter.FailureResponse(presenter.CreateTodoUnauthorizedErrorMessage, presenter.CreateTodoUnauthorizedErrorMessage))
		return
	}
	request.UserID = value.(int64)

	response, err := h.Service.Create(request.ToTodo())
	if err != nil {
		log.Error().Err(err).Msg("[TodoHandler][CreateTodo] Failed to create todo")
		c.JSON(http.StatusInternalServerError, presenter.FailureResponse(presenter.CreateTodoInternalServerErrorMessage, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, presenter.SuccessResponse(presenter.CreateTodoSuccessMessage, response))
}
