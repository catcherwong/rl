package rl

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
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

// LockWithId returns whether lock the key succeed
func (rl RedisLock) LockWithId(key string, id string, second int64) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key can not be null")
	}

	if id == "" {
		return false, fmt.Errorf("id can not be null")
	}

	if second <= 0 {
		return false, fmt.Errorf("second should greater than 0")
	}

	client, err := rl.getClient()

	if err != nil {
		return false, fmt.Errorf("redis ping error %v", err)
	}

	d := time.Duration(second) * time.Second

	f, err := client.SetNX(key, id, d).Result()

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

// UnLockWithId release a lock
func (rl RedisLock) UnLockWithId(key string, id string) (bool, error) {

	if key == "" {
		return false, fmt.Errorf("key can not be null")
	}

	if id == "" {
		return false, fmt.Errorf("id can not be null")
	}

	client, err := rl.getClient()

	if err != nil {
		return false, fmt.Errorf("redis del error %v", err)
	}

	script := redis.NewScript(`
	if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end
	`)

	v, err := script.Run(client, []string{key}, id).Result()

	if err != nil {
		return false, fmt.Errorf("redis del error %v", err)
	}

	return v > 0, nil
}
