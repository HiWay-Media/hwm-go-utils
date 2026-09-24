package redis

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var Brokers []string

// Init reads the cluster node addresses from REDIS_1 … REDIS_6, skipping unset ones.
func Init() {
	Brokers = nil
	for _, env := range []string{"REDIS_1", "REDIS_2", "REDIS_3", "REDIS_4", "REDIS_5", "REDIS_6"} {
		if addr := os.Getenv(env); addr != "" {
			Brokers = append(Brokers, addr)
		}
	}
}

// GetClient returns a cluster client for Brokers; REDIS_PASSWORD is used when set.
func GetClient() *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    Brokers,
		Password: os.Getenv("REDIS_PASSWORD"),
	})
}
