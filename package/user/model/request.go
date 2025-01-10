package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserDomain          = "user"
	PrefixKeyJWTPayload = "%s-jit:%s"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (r *RegisterRequest) HashPassword() (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("[UserRequest][HashPassword] Failed to hash password")
		return
	}

	r.Password = string(hashedPassword)

	return
}

func (r *RegisterRequest) ToUser() User {
	return User{
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: r.Password,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SetJWTPayloadRequest struct {
	UserID int64
	Email  string
	JIT    uuid.UUID
	Exp    time.Duration
	Key    string
	Value  ValueJWTPayload
}

func (request *SetJWTPayloadRequest) KeyJWTPayload() {
	request.Key = fmt.Sprintf(PrefixKeyJWTPayload, UserDomain, request.JIT.String())
}

type ValueJWTPayload struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}

func (request *ValueJWTPayload) ValidateTokenClaims(claims jwt.MapClaims) (err error) {
	if userID, ok := claims["sub"].(float64); !ok || userID != float64(request.UserID) {
		err = errors.New(ErrInvalidToken)
		log.Error().Err(err).Msg("[ValueJWTPayload][ValidateTokenClaims] Invalid token claims")
		return
	}
	if email, ok := claims["email"].(string); !ok || email != request.Email {
		err = errors.New(ErrInvalidToken)
		log.Error().Err(err).Msg("[ValueJWTPayload][ValidateTokenClaims] Invalid token claims")
		return
	}

	return nil
}

func (request *SetJWTPayloadRequest) ValueJWTPayload() {
	request.Value.Email = request.Email
	request.Value.UserID = request.UserID
}

type GetJWTPayloadRequest struct {
	JIT uuid.UUID
	Key string
}

func (request *GetJWTPayloadRequest) KeyJWTPayload() {
	request.Key = fmt.Sprintf(PrefixKeyJWTPayload, UserDomain, request.JIT.String())
}
