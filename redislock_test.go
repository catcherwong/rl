package rl

import (
	"testing"
	"time"
)

func TestRedisLock_Lock(t *testing.T) {
	key := "lock_key:1"
	var second int64 = 1

	rl := RedisLock{}

	v1, _ := rl.Lock(key, second)
	v2, _ := rl.Lock(key, second)

	if !v1 {
		t.Error("Lock should succeed, but fail")
	}

	if v2 {
		t.Error("Lock should fail, but succeed")
	}

	time.Sleep(time.Duration(second) * time.Second)

	v3, _ := rl.Lock(key, second)

	if !v3 {
		t.Error("After some seconds, it should be locked")
	}
}

func TestRedisLock_Lock2(t *testing.T) {
	key := ""
	var second int64 = 1

	rl := RedisLock{Addr: "localhost:6379"}

	v, err := rl.Lock(key, second)

	if err == nil {
		t.Error("key should not be empty")
	}

	if v {
		t.Error("result should be false")
	}
}

func TestRedisLock_Lock3(t *testing.T) {
	key := "lock_key:2"
	var second int64 = -1

	rl := RedisLock{}

	v, err := rl.Lock(key, second)

	if err == nil {
		t.Error("second should greater than 0")
	}

	if v {
		t.Error("result should be false")
	}
}

func TestRedisLock_UnLock(t *testing.T) {
	key := "unlock_key:1"
	var second int64 = 5

	rl := RedisLock{}

	v, _ := rl.Lock(key, second)

	if !v {
		t.Error("result should be true")
	}

	v, _ = rl.UnLock(key)

	if !v {
		t.Error("result should be true")
	}

	v, _ = rl.Lock(key, second)

	if !v {
		t.Error("result should be true after unlock")
	}
}

func TestRedisLock_UnLock2(t *testing.T) {
	key := ""

	rl := RedisLock{}

	v, err := rl.UnLock(key)

	if err == nil {
		t.Error("err should not be nil")
	}

	if v {
		t.Error("result should be false")
	}
}
