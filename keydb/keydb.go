package keydb

import (
	"github.com/go-redis/redis/v8"
)

// GetKeyDBClient returns a client for addr without authentication, DB 0.
func GetKeyDBClient(addr string) *redis.Client {
	return NewKeyDBClient(addr, "", 0)
}

// NewKeyDBClient returns a client for addr with the given password and DB.
func NewKeyDBClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
