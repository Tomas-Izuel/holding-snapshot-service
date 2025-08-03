package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// Connect establece la conexión con Redis
func Connect(redisURL string) error {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	RedisClient = redis.NewClient(opt)

	// Test de conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("✅ Conexión a Redis establecida")
	return nil
}

// GetClient retorna la instancia del cliente Redis
func GetClient() *redis.Client {
	return RedisClient
}

// Set almacena un valor en cache con TTL
func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return RedisClient.Set(ctx, key, value, ttl).Err()
}

// Get obtiene un valor del cache
func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// Delete elimina una clave del cache
func Delete(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}