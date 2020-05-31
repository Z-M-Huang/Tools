package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisSet(t *testing.T) {
	assert.Empty(t, RedisSet("TestRedisSet", 1, 1*time.Second))
}

func TestRedisSetBytes(t *testing.T) {
	assert.Empty(t, RedisSetBytes("TestRedisSetBytes", &ShortLink{}, 1*time.Second))
}

func TestRedisSetString(t *testing.T) {
	assert.Empty(t, RedisSetString("TestRedisSetString", "TestRedisSetString", 1*time.Second))
}

func TestRedisIncr(t *testing.T) {
	assert.Empty(t, RedisSet("TestRedisIncr", 1, 1*time.Second))
	n, err := RedisIncr("TestRedisIncr")
	assert.Empty(t, err)
	assert.Equal(t, int64(2), n)
}

func TestRedisDecr(t *testing.T) {
	assert.Empty(t, RedisSet("TestRedisDecr", 2, 1*time.Second))
	n, err := RedisDecr("TestRedisDecr")
	assert.Empty(t, err)
	assert.Equal(t, int64(1), n)
}

func TestRedisDelete(t *testing.T) {
	assert.Empty(t, RedisSet("TestRedisDelete", 1, 1*time.Second))
	assert.Empty(t, RedisDelete("TestRedisDelete"))
}

func TestRedisExist(t *testing.T) {
	assert.Empty(t, RedisSet("TestRedisExist", 1, 1*time.Second))
	assert.True(t, RedisExist("TestRedisExist"))
}

func TestRedisExistFalse(t *testing.T) {
	assert.False(t, RedisExist("TestRedisExistFalse"))
}

func TestRedisGetInt(t *testing.T) {
	assert.Empty(t, RedisSet("TestRedisGetInt", 1, 1*time.Second))
	n, err := RedisGetInt("TestRedisGetInt")
	assert.Empty(t, err)
	assert.Equal(t, int64(1), n)
}

func TestRedisGetIntNotExistFail(t *testing.T) {
	_, err := RedisGetInt("TestRedisGetIntNotExistFail")
	assert.NotEmpty(t, err)
}

func TestRedisGetIntFail(t *testing.T) {
	assert.Empty(t, RedisSet("TestRedisGetIntFail", "a", 1*time.Second))
	_, err := RedisGetInt("TestRedisGetIntFail")
	assert.NotEmpty(t, err)
}

func TestRedisGetString(t *testing.T) {
	assert.Empty(t, RedisSet("TestRedisGetString", "a", 1*time.Second))
	n, err := RedisGetString("TestRedisGetString")
	assert.Empty(t, err)
	assert.Equal(t, "a", n)
}

func TestRedisGetStringNotExistFail(t *testing.T) {
	_, err := RedisGetString("TestRedisGetStringNotExistFail")
	assert.NotEmpty(t, err)
}

func TestRedisGet(t *testing.T) {
	assert.Empty(t, RedisSetBytes("TestRedisGet", &ShortLink{}, 1*time.Second))
	s := &ShortLink{}
	assert.Empty(t, RedisGet("TestRedisGet", &s))
}

func TestRedisGetNotExistFail(t *testing.T) {
	s := &ShortLink{}
	assert.NotEmpty(t, RedisGet("TestRedisGetNotExistFail", &s))
}
