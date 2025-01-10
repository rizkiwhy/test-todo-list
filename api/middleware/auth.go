package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"rizkiwhy/test-todo-list/api/presenter"
	pkgUser "rizkiwhy/test-todo-list/package/user"
	mUser "rizkiwhy/test-todo-list/package/user/model"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AuthMiddleware struct {
	UserRepository      pkgUser.Repository
	UserCacheRepository pkgUser.CacheRepository
}

func NewAuthMiddleware(userRepository pkgUser.Repository, userCacheRepository pkgUser.CacheRepository) *AuthMiddleware {
	return &AuthMiddleware{
		UserRepository:      userRepository,
		UserCacheRepository: userCacheRepository,
	}
}

func (m *AuthMiddleware) AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Error().Msg("[AuthMiddleware][AuthJWT] Missing authorization header")
			abortWithUnauthorized(c, mUser.ErrUnauthorizedAccess, MissingAuthHeaderErrorMessage)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := parseJWTToken(tokenString)
		if err != nil {
			abortWithUnauthorized(c, mUser.ErrUnauthorizedAccess, ErrInvalidTokenMessage)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Error().Msg("[AuthMiddleware][AuthJWT] Invalid token claims or token is invalid")
			abortWithUnauthorized(c, mUser.ErrUnauthorizedAccess, ErrInvalidTokenClaimsMessage)
			return
		}

		payload, err := m.fetchJWTPayloadFromCache(claims["jit"])
		if err != nil {
			log.Error().Err(err).Msg("[AuthMiddleware][AuthJWT] Failed to fetch JWT payload from cache")
			abortWithUnauthorized(c, mUser.ErrUnauthorizedAccess, ErrInvalidTokenClaimsMessage)
			return
		}

		if err := payload.ValidateTokenClaims(claims); err != nil {
			log.Error().Err(err).Msg("[AuthMiddleware][AuthJWT] Invalid token claims compare with payload")
			abortWithUnauthorized(c, mUser.ErrUnauthorizedAccess, ErrInvalidTokenClaimsMessage)
			return
		}

		user, err := m.UserRepository.GetByEmail(payload.Email)
		if err != nil {
			log.Error().Err(err).Msg("[AuthMiddleware][AuthJWT] Failed to find user by email")
			abortWithUnauthorized(c, mUser.ErrUnauthorizedAccess, ErrInvalidTokenClaimsMessage)
			return
		}

		if !user.ValidateTokenClaimsSub(payload.UserID, claims["sub"].(float64)) {
			log.Error().Msg("[AuthMiddleware][AuthJWT] Invalid token claims compoare with user")
			abortWithUnauthorized(c, mUser.ErrUnauthorizedAccess, ErrInvalidTokenClaimsMessage)
			return
		}

		c.Set("user_id", payload.UserID)
		c.Set("email", payload.Email)

		c.Next()
	}
}

func parseJWTToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error().Msg("[AuthMiddleware][parseJWTToken] Invalid signing method")
			return nil, errors.New(ErrInvalidSigningMethodMessage)
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func (m *AuthMiddleware) fetchJWTPayloadFromCache(jit interface{}) (*mUser.ValueJWTPayload, error) {
	jitUUID, err := uuid.Parse(fmt.Sprintf("%v", jit))
	if err != nil {
		log.Error().Err(err).Msg("[AuthMiddleware][fetchJWTPayloadFromCache] Failed to parse JIT UUID")
		return nil, errors.New("invalid JIT")
	}

	payload, err := m.UserCacheRepository.GetJWTPayload(mUser.GetJWTPayloadRequest{JIT: jitUUID})
	if err != nil {
		log.Error().Err(err).Msg("[AuthMiddleware][fetchJWTPayloadFromCache] Failed to retrieve JWT payload from cache")
		return nil, err
	}

	return payload, nil
}

func abortWithUnauthorized(c *gin.Context, title, message string) {
	log.Error().Msgf("[AuthMiddleware] Unauthorized: %s - %s", title, message)
	c.AbortWithStatusJSON(http.StatusUnauthorized, presenter.FailureResponse(title, message))
}
