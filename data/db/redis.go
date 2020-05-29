package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v7"
)

var redisClient *redis.Client
var locker *redislock.Client

//InitRedis init Redis
func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     data.Config.RedisConfig.Addr,
		Password: data.Config.RedisConfig.Password,
		DB:       data.Config.RedisConfig.DB,
	})
	err := redisClient.Ping().Err()
	if err != nil {
		utils.Logger.Fatal("failed to connect to Redis")
	}

	locker = redislock.New(redisClient)
}

//RedisGet get value
func RedisGet(key string, out interface{}) error {
	val, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		out = nil
		return nil
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

//RedisGetString get string
func RedisGetString(key string) (string, error) {
	val, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		utils.Logger.Error(err.Error())
		return "", errors.New("Internal Error")
	}
	return val, nil
}

//RedisGetInt get int64
func RedisGetInt(key string) (int64, error) {
	val, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		return 0, errors.New("Not found")
	} else if err != nil {
		utils.Logger.Error(err.Error())
		return 0, errors.New("Internal Error")
	}

	ret, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		utils.Logger.Error(err.Error())
		return 0, fmt.Errorf("Failed to get %s as int", key)
	}

	return ret, nil
}

//RedisExist check key exists
func RedisExist(key string) bool {
	i, err := redisClient.Exists(key).Result()
	if err != nil || i < 1 {
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

//RedisSetString set string
func RedisSetString(key string, val string, expire time.Duration) error {
	bytes := []byte(val)
	return RedisSet(key, bytes, expire)
}

//RedisSetBytes set bytes
func RedisSetBytes(key string, val interface{}, expire time.Duration) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return RedisSet(key, bytes, expire)
}

//RedisIncr redis increase
func RedisIncr(key string) (int64, error) {
	result, err := redisClient.Incr(key).Result()
	if err != nil {
		utils.Logger.Error(err.Error())
		return 0, fmt.Errorf("Failed to incr %s", key)
	}
	return result, nil
}

//RedisDecr redis decrease
func RedisDecr(key string) (int64, error) {
	result, err := redisClient.Decr(key).Result()
	if err != nil {
		utils.Logger.Error(err.Error())
		return 0, fmt.Errorf("Failed to decr %s", key)
	}
	return result, nil
}

//RedisLock lock in redis
func RedisLock(key string, duration time.Duration) *redislock.Lock {
	lock, err := locker.Obtain(key, duration, nil)
	if err == redislock.ErrNotObtained {
		return nil
	} else if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	return lock
}
