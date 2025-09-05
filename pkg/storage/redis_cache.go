package storage

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements the Cache interface using Redis as backend.
type RedisCache[T any] struct {
	client    *redis.Client
	keyPrefix string
}

// TODO: Make it to use connection pool initialized on server instance like s.datastore
func NewRedisCache[T any](address, password, keyPrefix string) *RedisCache[T] {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: "default",
		Password: password,
		DB:       0,
	})

	status := client.Conn().Ping(context.Background())

	fmt.Println("DEBUG: Redis connection status:", status.String())

	return &RedisCache[T]{client: client, keyPrefix: keyPrefix}
}

func (r *RedisCache[T]) Get(key string) T {
	var zero T
	ctx := context.Background()

	fmt.Println("DEBUG: Getting cache for key:", r.keyPrefix+key)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		fmt.Println("DEBUG: Error getting cache:", err)
		return zero
	}
	fmt.Println("DEBUG: retrieved value:", data)
	var value T
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&value); err != nil {
		fmt.Println("DEBUG: Error decoding cache value:", err)
		return zero
	}

	return value
}

func (r *RedisCache[T]) Set(key string, value T, ttl time.Duration) {
	fmt.Println("DEBUG: Setting cache for key:", r.keyPrefix+key)
	if ttl < 0 {
		fmt.Println("DEBUG: Invalid TTL, must be positive?")
		return
	}
	ctx := context.Background()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		fmt.Println("DEBUG: Error encoding cache value:", err)
		fmt.Println(value)
		return
	}

	if err := r.client.Set(ctx, key, buf.Bytes(), ttl).Err(); err != nil {
		fmt.Println("Error setting cache:", err) // TODO inject logger
		return
	}
}

func (r *RedisCache[T]) Delete(key string) {
	ctx := context.Background()
	_ = r.client.Del(ctx, key).Err()
}

func (r *RedisCache[T]) Stop() {
	_ = r.client.Close()
}

var _ Cache[any] = (*RedisCache[any])(nil)
