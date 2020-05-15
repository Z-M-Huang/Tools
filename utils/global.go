package utils

import (
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql" //supporting packages
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var onceRedis sync.Once

//Templates page templates
var Templates *template.Template

//DB database connection
var DB *gorm.DB

//RedisClient redis instance
var RedisClient *redis.Client

//Logger global Logger
var Logger *zap.Logger

//Config application config
var Config *data.Configuration

//AppList in home page
var AppList []*webdata.AppCategory

const (
	//ClaimCtxKey claim context key
	ClaimCtxKey string = "claim_context_key"
	//ResponseCtxKey page data context key
	ResponseCtxKey string = "response_context_key"
	//SessionTokenKey auth token in session key
	SessionTokenKey string = "session_token"
	//UsedTokenKey contains apps used. This is a slice of app titles
	UsedTokenKey string = "used_token"
	//PageStyleCookieKey allows users to choose the style wants to use for page
	PageStyleCookieKey string = "page_style"
)

func init() {
	onceRedis.Do(func() {
		initLogger()
		initTemplates()
		initConfig()
		initDB()
		initRedis()
	})
}

func initLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.OutputPaths = []string{"stdout"}
	Logger, _ = config.Build()
}

func initTemplates() {
	var err error
	Templates = template.New("")
	getTemplateFuncs()
	Templates, err = Templates.ParseFiles(getAlltemplates("templates/")...)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

func getAlltemplates(inputPath string) []string {
	var ret []string
	filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if path != inputPath && info.IsDir() {
			ret = append(ret, getAlltemplates(path)...)
		} else if strings.Contains(info.Name(), ".gohtml") {
			ret = append(ret, path)
		}
		return nil
	})
	return ret
}

func getTemplateFuncs() {
	Templates.Funcs(template.FuncMap{"add": func(i, j int) int { return i + j }})
	Templates.Funcs(template.FuncMap{"mod": func(i, j int) int { return i % j }})
	Templates.Funcs(template.FuncMap{"nospace": func(i string) string {
		return strings.ReplaceAll(i, " ", "")
	}})
}

func initConfig() {
	redisDBNum, err := strconv.ParseInt(strings.TrimSpace(os.Getenv("REDIS_DB")), 10, 32)
	if err != nil {
		Logger.Sugar().Warnf("failed to parse redis db number, set to default 0 %s", err.Error())
		redisDBNum = 0
	}

	genSitemap := false
	genSitemapStr := strings.TrimSpace(os.Getenv("ENABLE_SITEMAP"))
	if genSitemapStr == "" {
		Logger.Warn("ENABLE_SITEMAP is empty, set to default: false")
	} else {
		genSitemap, err = strconv.ParseBool(genSitemapStr)
		if err != nil {
			Logger.Error("Failed to parse ENABLE_SITEMAP to boolean, set to default: false")
		}
	}

	isHTTPS := false
	isHTTPSStr := strings.TrimSpace(os.Getenv("HTTPS"))
	if isHTTPSStr == "" {
		Logger.Warn("HTTPS is empty, set to default: false")
	} else {
		isHTTPS, err = strconv.ParseBool(isHTTPSStr)
		if err != nil {
			Logger.Error("Failed to parse SITEMAP_HTTPS to boolean, set to default: false")
		}
	}

	Config = &data.Configuration{
		DatabaseConfig: &data.DatabaseConfiguration{
			ConnectionString: strings.TrimSpace(os.Getenv("CONNECTION_STRING")),
			Driver:           strings.TrimSpace(os.Getenv("DB_DRIVER")),
		},
		RedisConfig: &data.RedisConfiguration{
			Addr:     strings.TrimSpace(os.Getenv("REDIS_ADDR")),
			Password: strings.TrimSpace(os.Getenv("REDIS_PASSWORD")),
			DB:       int(redisDBNum),
		},
		GoogleOauthConfig: &data.GoogleOauthConfiguration{
			ClientID:     strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
			ClientSecret: strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET")),
		},
		JwtKey:          []byte(strings.TrimSpace(os.Getenv("JWT_KEY"))),
		Host:            strings.TrimSuffix(strings.TrimSpace(os.Getenv("HOST")), "/"),
		ResourceVersion: strings.TrimSpace(os.Getenv("RESOURCE_VERSION")),
		IsDebug:         os.Getenv("DEBUG") != "",
		IsHTTPS:         isHTTPS,
		SitemapConfig: &data.SitemapConfiguration{
			GenerateSitemap: genSitemap,
		},
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

func initDB() {
	var err error
	DB, err = gorm.Open(Config.DatabaseConfig.Driver, Config.DatabaseConfig.ConnectionString)
	if err != nil {
		Logger.Sugar().Fatalf("failed to open database %s", err.Error())
	}
	DB.AutoMigrate(&dbentity.User{}, &dbentity.Application{})
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
