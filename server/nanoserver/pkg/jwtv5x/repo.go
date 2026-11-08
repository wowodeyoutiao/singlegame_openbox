package jwtv5x

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type tokenRepo struct {
	cache redis.UniversalClient
}

func NewTokenRepo(
	cache redis.UniversalClient,
) *tokenRepo {
	return &tokenRepo{
		cache: cache,
	}
}

func (r *tokenRepo) Set(ctx context.Context, key string, value string, duration time.Duration) error {
	return r.cache.Set(ctx, key, value, duration).Err()
}

func (r *tokenRepo) Exist(ctx context.Context, key string) (bool, error) {
	result, err := r.cache.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if result == 0 {
		return false, nil
	}
	return true, nil
}

func (r *tokenRepo) Get(ctx context.Context, key string) (string, error) {
	val, err := r.cache.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}
