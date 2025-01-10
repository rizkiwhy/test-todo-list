package user

import (
	"errors"
	"rizkiwhy/test-todo-list/package/user/model"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ServiceImpl struct {
	Repository      Repository
	CacheRepository CacheRepository
}

type Service interface {
	Register(request model.RegisterRequest) (response model.RegisterResponse, err error)
	Login(request model.LoginRequest) (response model.LoginResponse, err error)
}

func NewService(repository Repository, cacheRepository CacheRepository) Service {
	return &ServiceImpl{
		Repository:      repository,
		CacheRepository: cacheRepository,
	}
}

func (s *ServiceImpl) Register(request model.RegisterRequest) (response model.RegisterResponse, err error) {
	isExists, err := s.Repository.IsExistsByEmail(request.Email)
	if err != nil {
		log.Error().Err(err).Str("email", request.Email).Msg("[UserService][Register] Failed to count user by email")
		return
	}

	if isExists {
		err = errors.New(model.ErrEmailAlreadyExists)
		log.Error().Err(err).Str("email", request.Email).Msg("[UserService][Register] Email already exists")
		return
	}

	err = request.HashPassword()
	if err != nil {
		log.Error().Err(err).Msg("[UserService][Register] Failed to hash password")
		return
	}

	user := request.ToUser()

	result, err := s.Repository.Create(user)
	if err != nil {
		log.Error().Err(err).Interface("user", request).Msg("[UserService][Register] Failed to create user")
	}

	return result.ToRegisterResponse(), nil
}

func (s *ServiceImpl) Login(request model.LoginRequest) (response model.LoginResponse, err error) {
	user, err := s.Repository.GetByEmail(request.Email)
	if err != nil {
		log.Error().Err(err).Str("email", request.Email).Msg("[UserService][Login] Failed to get user by email")
		return
	}

	if !user.IsExist() {
		err = errors.New(model.ErrInvalidEmailPassword)
		log.Error().Err(err).Str("email", request.Email).Msg("[UserService][Login] User not found")
		return
	}

	err = user.CompareHashPassword(request.Password)
	if err != nil {
		log.Error().Err(err).Str("email", request.Email).Msg("[UserService][Login] Failed to compare hash password")
		return
	}

	resToken, err := user.CreateToken()
	if err != nil {
		log.Error().Err(err).Str("email", request.Email).Msg("[UserService][Login] Failed to create token")
		return
	}

	go s.CacheRepository.SetJWTPayload(model.SetJWTPayloadRequest{
		UserID: user.ID,
		Email:  user.Email,
		JIT:    uuid.MustParse(resToken.JIT),
		Exp:    resToken.Exp,
	})

	response.Token = &resToken.Token

	return
}
