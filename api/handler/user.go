package handler

import (
	"net/http"
	"rizkiwhy/test-todo-list/api/presenter"
	pkgUser "rizkiwhy/test-todo-list/package/user"
	"rizkiwhy/test-todo-list/package/user/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	Service pkgUser.Service
}

func NewUserHandler(service pkgUser.Service) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var request model.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("[UserHandler][Register] Failed to bind json")
		c.JSON(http.StatusBadRequest, presenter.FailureResponse(presenter.RegisterInvalidRequestMessage, err.Error()))
		return
	}

	response, err := h.Service.Register(request)
	if err != nil {
		log.Error().Err(err).Msg("[UserHandler][Register] Failed to register user")
		presenter.HandleError(c, err, presenter.RegisterStatusCodeMap, presenter.RegisterFailureMessage)
		return
	}

	c.JSON(http.StatusCreated, presenter.SuccessResponse(presenter.RegisterSuccessMessage, response))
}

func (h *UserHandler) Login(c *gin.Context) {
	var request model.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("[UserHandler][Login] Failed to bind json")
		c.JSON(http.StatusBadRequest, presenter.FailureResponse(presenter.LoginInvalidCredentialsMessage, err.Error()))
		return
	}

	response, err := h.Service.Login(request)
	if err != nil {
		log.Error().Err(err).Msg("[UserHandler][Login] Failed to login user")
		presenter.HandleError(c, err, presenter.LoginStatusCodeMap, presenter.LoginFailureMessage)
		return
	}

	c.JSON(http.StatusOK, presenter.SuccessResponse(presenter.LoginSuccessMessage, response))
}
