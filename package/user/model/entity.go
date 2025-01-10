package model

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64      `gorm:"primaryKey;autoIncrement;not null"`
	Name         string     `gorm:"size:255;not null"`
	Email        string     `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string     `gorm:"size:255;not null"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;not null"`
	UpdatedAt    *time.Time `gorm:"autoUpdateTime"`
}

func (u *User) ValidateTokenClaimsSub(sub int64, userID float64) bool {
	return u.ID == sub && u.ID == int64(userID)
}

func (u *User) ToRegisterResponse() RegisterResponse {
	return RegisterResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func (u *User) IsExist() bool {
	return u != nil
}

func (u *User) CompareHashPassword(password string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		log.Error().Err(err).Msg("[UserEntity][CompareHashPassword] Failed to compare hash password")
		err = errors.New(ErrInvalidEmailPassword)
	}

	return
}

func (u *User) CreateToken() (*TokenResponse, error) {
	jit := uuid.New().String()
	expTime, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		log.Error().Err(err).Msg("[UserEntity][CreateToken] Failed to convert JWT_EXPIRATION to int")

	}
	exp := time.Hour * time.Duration(expTime)

	now := time.Now()

	claims := jwt.MapClaims{
		"sub":   u.ID,
		"email": u.Email,
		"iat":   now.Unix(),
		"exp":   now.Add(exp).Unix(),
		"jit":   jit,
		"iss":   os.Getenv("JWT_ISSUER"),
	}

	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Error().Err(err).Msg("[UserEntity][CreateToken] Failed to sign token")
		return nil, err
	}

	return &TokenResponse{
		Token: signedToken,
		JIT:   jit,
		Exp:   exp,
	}, nil
}

func (u *User) SetJWTPayloadRequest(request TokenResponse) SetJWTPayloadRequest {
	return SetJWTPayloadRequest{
		UserID: u.ID,
		Email:  u.Email,
		JIT:    uuid.MustParse(request.JIT),
		Exp:    request.Exp,
	}
}
