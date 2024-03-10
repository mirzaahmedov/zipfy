package handler

import (
	"zipfy/internal/store"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	store     store.Store
	jwtSecret string
	redis     *redis.Client
}

func NewHandler(store store.Store, jwtSecret string, redis *redis.Client) *Handler {
	return &Handler{
		store,
		jwtSecret,
		redis,
	}
}
