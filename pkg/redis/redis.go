package redis

import (
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func New(addr, password string, db int) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	client := Redis{rdb}

	return &client
}
