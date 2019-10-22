package rl

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type RedisLock struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
}

func (rl RedisLock) getClient() (*redis.Client, error) {

	if rl.Addr == "" {
		rl.Addr = "localhost:6379"
	}

	if rl.PoolSize <= 0 {
		rl.PoolSize = 10
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     rl.Addr,
		Password: rl.Password,
		DB:       rl.DB,
		PoolSize: rl.PoolSize,
	})

	err := redisClient.Ping().Err()

	return redisClient, err
}

// Lock returns whether lock the key succeed
func (rl RedisLock) Lock(key string, second int64) (bool, error) {

	if key == "" {
		return false, fmt.Errorf("key can not be null")
	}

	if second <= 0 {
		return false, fmt.Errorf("second should greater than 0")
	}

	client, err := rl.getClient()

	if err != nil {
		return false, fmt.Errorf("redis ping error %v", err)
	}

	d := time.Duration(second) * time.Second

	f, err := client.SetNX(key, 1, d).Result()

	if err != nil {
		return false, fmt.Errorf("redis setnx error %v", err)
	}

	return f, nil
}

// UnLock release a lock
func (rl RedisLock) UnLock(key string) (bool, error) {

	if key == "" {
		return false, fmt.Errorf("key can not be null")
	}

	client, err := rl.getClient()

	if err != nil {
		return false, fmt.Errorf("redis del error %v", err)
	}

	v, err := client.Del(key).Result()

	if err != nil {
		return false, fmt.Errorf("redis del error %v", err)
	}

	return v > 0, nil
}
