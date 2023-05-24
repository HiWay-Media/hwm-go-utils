package redis


import (
	"github.com/go-redis/redis/v8"
	"os"
)

var Brokers []string

func Init() {
	Brokers = []string{
		os.Getenv("REDIS_1"),
		os.Getenv("REDIS_2"),
		os.Getenv("REDIS_3"),
		os.Getenv("REDIS_4"),
		os.Getenv("REDIS_5"),
		os.Getenv("REDIS_6"),
	}
}

func GetClient() *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: Brokers,
	})
}