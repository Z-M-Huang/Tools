package utils

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var onceRedis sync.Once

//RedisClient redis instance
var RedisClient *redis.Client

//Logger global Logger
var Logger *zap.Logger

//Config application config
var Config *data.Configuration

func init() {
	onceRedis.Do(func() {
		initLogger()
		initConfig()
		initRedis()
	})
}

func initLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.OutputPaths = []string{"stdout"}
	Logger, _ = config.Build()
}

func initConfig() {
	redisDBNum, err := strconv.ParseInt(strings.TrimSpace(os.Getenv("REDIS_DB")), 10, 32)
	if err != nil {
		Logger.Sugar().Errorf("failed to parse redis db number, set to default 0 %s", err.Error())
		redisDBNum = 0
	}

	Config = &data.Configuration{
		RedisConfig: &data.RedisConfiguration{
			Addr:     strings.TrimSpace(os.Getenv("REDIS_ADDR")),
			Password: strings.TrimSpace(os.Getenv("REDIS_PASSWORD")),
			DB:       int(redisDBNum),
		},
		GoogleOauthConfig: &data.GoogleOauthConfiguration{
			ClientID:     strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
			ClientSecret: strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET")),
		},
		JwtKey: []byte(strings.TrimSpace(os.Getenv("JWT_KEY"))),
		Host:   strings.TrimSpace(os.Getenv("HOST")),
	}

	if Config.RedisConfig.Addr == "" {
		Logger.Fatal("REDIS_ADDR cannot be empty")
	} else if Config.GoogleOauthConfig.ClientID == "" {
		Logger.Fatal("GOOGLE_CLIENT_ID cannot be empty")
	} else if Config.GoogleOauthConfig.ClientSecret == "" {
		Logger.Fatal("GOOGLE_CLIENT_SECRET cannot be empty")
	} else if len(Config.JwtKey) == 0 {
		Logger.Fatal("JWT_KEY cannot be empty")
	} else if Config.Host == "" {
		Logger.Fatal("HOST cannot be empty")
	}
}

func initRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     Config.RedisConfig.Addr,
		Password: Config.RedisConfig.Password,
		DB:       Config.RedisConfig.DB,
	})
	err := RedisClient.Ping().Err()
	if err != nil {
		Logger.Fatal("failed to connect to Redis")
	}
}
