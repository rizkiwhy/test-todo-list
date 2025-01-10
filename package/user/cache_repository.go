package user

import (
	"encoding/json"
	"rizkiwhy/test-todo-list/package/user/model"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type CacheRepository interface {
	SetJWTPayload(request model.SetJWTPayloadRequest) (err error)
	GetJWTPayload(request model.GetJWTPayloadRequest) (valueJWTPayload *model.ValueJWTPayload, err error)
}

type CacheRepositoryImpl struct {
	RedisClient *redis.Client
}

func NewCacheRepository(redisClient *redis.Client) CacheRepository {
	return &CacheRepositoryImpl{
		RedisClient: redisClient,
	}
}

func (c *CacheRepositoryImpl) SetJWTPayload(request model.SetJWTPayloadRequest) (err error) {
	request.KeyJWTPayload()
	request.ValueJWTPayload()

	valueJSON, err := json.Marshal(request.Value)
	if err != nil {
		log.Error().Err(err).Msg("[SetJWTPayload] Failed to marshal value to JSON")
		return
	}

	err = c.RedisClient.Set(c.RedisClient.Context(), request.Key, valueJSON, request.Exp).Err()
	if err != nil {
		log.Error().Err(err).Msg("[SetJWTPayload] Failed to set value in Redis")
	}

	return
}

func (c *CacheRepositoryImpl) GetJWTPayload(request model.GetJWTPayloadRequest) (valueJWTPayload *model.ValueJWTPayload, err error) {
	request.KeyJWTPayload()
	value, err := c.RedisClient.Get(c.RedisClient.Context(), request.Key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Debug().Str("key", request.Key).Msg("[GetJWTPayload] Key not found in Redis")
			return
		}
		log.Error().Err(err).Str("key", request.Key).Msg("[GetJWTPayload] Failed to get value from Redis")
		return
	}

	err = json.Unmarshal([]byte(value), &valueJWTPayload)
	if err != nil {
		log.Error().Err(err).Msg("[GetJWTPayload] Failed to unmarshal value from JSON")
	}

	return
}
