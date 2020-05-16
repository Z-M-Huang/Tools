package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func init() {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     data.Config.RedisConfig.Addr,
		Password: data.Config.RedisConfig.Password,
		DB:       data.Config.RedisConfig.DB,
	})
	err := redisClient.Ping().Err()
	if err != nil {
		utils.Logger.Fatal("failed to connect to Redis")
	}
}

//RedisGet get value
func RedisGet(key string, out interface{}) error {
	val, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		return errors.New("Not found")
	} else if err != nil {
		utils.Logger.Error(err.Error())
		return errors.New("Internal Error")
	}
	err = json.Unmarshal([]byte(val), &out)
	if err != nil {
		utils.Logger.Error(err.Error())
		return errors.New("Internal Error")
	}
	return nil
}

//RedisExist check key exists
func RedisExist(key string) bool {
	err := redisClient.Exists(key).Err()
	if err != nil {
		return false
	}
	return true
}

//RedisDelete delete key
func RedisDelete(key string) error {
	err := redisClient.Exists(key).Err()
	if err != nil {
		utils.Logger.Error(err.Error())
		return fmt.Errorf("Failed to delete %s in Redis", key)
	}
	return nil
}

//RedisSet set value
func RedisSet(key string, val interface{}, expire time.Duration) error {
	err := redisClient.Set(key, val, expire).Err()
	if err != nil {
		utils.Logger.Error(err.Error())
		return fmt.Errorf("Failed to set key %s in Redis", key)
	}
	return nil
}

//RedisSetBytes set bytes
func RedisSetBytes(key string, val interface{}, expire time.Duration) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return RedisSet(key, bytes, expire)
}
