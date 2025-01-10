package model

import "time"

const (
	ErrEmailAlreadyExists   = "email already exists"
	ErrPasswordHashing      = "failed to hash password"
	ErrDatabase             = "failed to connect to database"
	ErrUserCreation         = "failed to create user"
	ErrInvalidRequest       = "invalid request"
	ErrUnauthorizedAccess   = "unauthorized access"
	ErrNotFound             = "record not found"
	ErrInternalError        = "internal server error"
	ErrInvalidEmailPassword = "invalid email or password"
	ErrMissingFields        = "missing required fields"
	ErrEmailFormat          = "invalid email format"
	ErrPasswordLength       = "password must be at least 8 characters long"
	ErrPasswordComplexity   = "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
	ErrTokenGeneration      = "failed to generate token"
	ErrInvalidToken         = "invalid token"
)

type RegisterResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TokenResponse struct {
	Token string
	JIT   string
	Exp   time.Duration
}

type LoginResponse struct {
	Token *string `json:"token,omitempty"`
}

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
