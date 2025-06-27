package db

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// ConnectRedis connects to the Redis database using the given configuration.
func ConnectRedis() (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_DB_HOST")
	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_DB_PORT"))
	redisUsername := os.Getenv("REDIS_DB_USERNAME")
	redisPassword := os.Getenv("REDIS_DB_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Username: redisUsername,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return nil, err
	}

	fmt.Println("Redis connected")

	return client, nil
}

// DisconnectRedis disconnects from the Redis database.
func DisconnectRedis(client *redis.Client) error {
	err := client.Close()
	if err != nil {
		return err
	}

	fmt.Println("Redis disconnected")
	return nil
}
