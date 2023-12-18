package keydb

import (
	"github.com/go-redis/redis/v8"
)


func GetKeyDBClient(addr string) *redis.ClusterClient {
	return redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: "", // no password set
        DB:       0,  // use default DB
    })
}