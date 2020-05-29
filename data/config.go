package data

import (
	"os"
	"strconv"
	"strings"

	"github.com/Z-M-Huang/Tools/utils"
)

//Config global config
var Config *Configuration

//Configuration app configuration
type Configuration struct {
	DatabaseConfig    *DatabaseConfiguration
	RedisConfig       *RedisConfiguration
	GoogleOauthConfig *GoogleOauthConfiguration
	EmailConfig       *EmailConfiguration

	JwtKey          []byte
	Host            string
	ResourceVersion string

	HTTPS         bool
	EnableSitemap bool
	IsDebug       bool
}

//RedisConfiguration redis config
type RedisConfiguration struct {
	Addr     string
	Password string
	DB       int
}

//GoogleOauthConfiguration google oauth2
type GoogleOauthConfiguration struct {
	ClientID     string
	ClientSecret string
}

//DatabaseConfiguration gorm
type DatabaseConfiguration struct {
	ConnectionString string
	Driver           string
}

//EmailConfiguration email configuration for email sms
type EmailConfiguration struct {
	SMTPServer   string
	IMAPServer   string
	EmailAddress string
	Password     string
}

//LoadProductionConfig load config
func LoadProductionConfig() {
	redisDBNum, err := strconv.ParseInt(strings.TrimSpace(os.Getenv("REDIS_DB")), 10, 32)
	if err != nil {
		utils.Logger.Sugar().Warnf("failed to parse redis db number, set to default 0 %s", err.Error())
		redisDBNum = 0
	}

	enableSitemap := false
	enableSitemapStr := strings.TrimSpace(os.Getenv("ENABLE_SITEMAP"))
	if enableSitemapStr == "" {
		utils.Logger.Warn("ENABLE_SITEMAP is empty, set to default: false")
	} else {
		enableSitemap, err = strconv.ParseBool(enableSitemapStr)
		if err != nil {
			utils.Logger.Error("Failed to parse ENABLE_SITEMAP to boolean, set to default: false")
		}
	}

	isHTTPS := false
	isHTTPSStr := strings.TrimSpace(os.Getenv("HTTPS"))
	if isHTTPSStr == "" {
		utils.Logger.Warn("HTTPS is empty, set to default: false")
	} else {
		isHTTPS, err = strconv.ParseBool(isHTTPSStr)
		if err != nil {
			utils.Logger.Error("Failed to parse SITEMAP_HTTPS to boolean, set to default: false")
		}
	}

	Config = &Configuration{
		DatabaseConfig: &DatabaseConfiguration{
			ConnectionString: strings.TrimSpace(os.Getenv("CONNECTION_STRING")),
			Driver:           strings.TrimSpace(os.Getenv("DB_DRIVER")),
		},
		RedisConfig: &RedisConfiguration{
			Addr:     strings.TrimSpace(os.Getenv("REDIS_ADDR")),
			Password: strings.TrimSpace(os.Getenv("REDIS_PASSWORD")),
			DB:       int(redisDBNum),
		},
		GoogleOauthConfig: &GoogleOauthConfiguration{
			ClientID:     strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
			ClientSecret: strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET")),
		},
		EmailConfig: &EmailConfiguration{
			SMTPServer:   strings.TrimSpace(os.Getenv("SMTP_SERVER")),
			IMAPServer:   strings.TrimSpace(os.Getenv("IMAP_SERVER")),
			EmailAddress: strings.TrimSpace(os.Getenv("EMAIL_ADDRESS")),
			Password:     strings.TrimSpace(os.Getenv("EMAIL_PASSWORD")),
		},
		JwtKey:          []byte(strings.TrimSpace(os.Getenv("JWT_KEY"))),
		Host:            strings.TrimSuffix(strings.TrimSpace(os.Getenv("HOST")), "/"),
		ResourceVersion: strings.TrimSpace(os.Getenv("RESOURCE_VERSION")),
		IsDebug:         os.Getenv("DEBUG") != "",
		HTTPS:           isHTTPS,
		EnableSitemap:   enableSitemap,
	}

	if Config.RedisConfig.Addr == "" {
		utils.Logger.Fatal("REDIS_ADDR cannot be empty")
	} else if Config.GoogleOauthConfig.ClientID == "" {
		utils.Logger.Fatal("GOOGLE_CLIENT_ID cannot be empty")
	} else if Config.GoogleOauthConfig.ClientSecret == "" {
		utils.Logger.Fatal("GOOGLE_CLIENT_SECRET cannot be empty")
	} else if len(Config.JwtKey) == 0 {
		utils.Logger.Fatal("JWT_KEY cannot be empty")
	} else if Config.Host == "" {
		utils.Logger.Fatal("HOST cannot be empty")
	}

	if Config.EmailConfig.SMTPServer == "" {
		utils.Logger.Error("SMTP_SERVER is empty. Some feature may not work...")
	} else if Config.EmailConfig.EmailAddress == "" {
		utils.Logger.Error("EMAIL_ADDRESS is empty. Some feature may not work...")
	} else if Config.EmailConfig.Password == "" {
		utils.Logger.Error("EMAIL_PASSWORD is empty. Some feature may not work...")
	} else if Config.EmailConfig.IMAPServer == "" {
		utils.Logger.Error("IMAP_SERVER is empty. Some feature may not work...")
	}
}
