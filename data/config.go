package data

import (
	"os"
	"strings"

	"github.com/Z-M-Huang/Tools/utils"
	"github.com/caarlos0/env/v6"
)

//Config global config
var (
	Config            *Configuration
	DatabaseConfig    *DatabaseConfiguration
	RedisConfig       *RedisConfiguration
	GoogleOauthConfig *GoogleOauthConfiguration
	EmailConfig       *EmailConfiguration
)

//Configuration app configuration
type Configuration struct {
	RapidAPIKey string `env:"RAPIDAPI_KEY"`

	JwtKey          []byte
	Host            string `env:"HOST"`
	ResourceVersion string `env:"RESOURCE_VERSION"`

	HTTPS         bool `env:"HTTPS" envDefault:"false"`
	EnableSitemap bool `env:"ENABLE_SITEMAP" envDefault:"false"`
	IsDebug       bool `env:"DEBUG" envDefault:"false"`
}

//RedisConfiguration redis config
type RedisConfiguration struct {
	Addr     string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

//GoogleOauthConfiguration google oauth2
type GoogleOauthConfiguration struct {
	ClientID     string `env:"GOOGLE_CLIENT_ID"`
	ClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
}

//DatabaseConfiguration gorm
type DatabaseConfiguration struct {
	ConnectionString string `env:"CONNECTION_STRING" envDefault:"db.db"`
	Driver           string `env:"DB_DRIVER" envDefault:"sqlite3"`
}

//EmailConfiguration email configuration for email sms
type EmailConfiguration struct {
	SMTPServer   string `env:"SMTP_SERVER"`
	EmailAddress string `env:"EMAIL_ADDRESS"`
	Password     string `env:"EMAIL_PASSWORD"`
}

//LoadProductionConfig load config
func LoadProductionConfig() {
	Config = &Configuration{}
	if err := env.Parse(Config); err != nil {
		utils.Logger.Fatal(err.Error())
	}
	DatabaseConfig = &DatabaseConfiguration{}
	if err := env.Parse(DatabaseConfig); err != nil {
		utils.Logger.Fatal(err.Error())
	}
	RedisConfig = &RedisConfiguration{}
	if err := env.Parse(RedisConfig); err != nil {
		utils.Logger.Fatal(err.Error())
	}
	GoogleOauthConfig = &GoogleOauthConfiguration{}
	if err := env.Parse(GoogleOauthConfig); err != nil {
		utils.Logger.Fatal(err.Error())
	}
	EmailConfig = &EmailConfiguration{}
	if err := env.Parse(EmailConfig); err != nil {
		utils.Logger.Fatal(err.Error())
	}

	Config.JwtKey = []byte(strings.TrimSpace(os.Getenv("JWT_KEY")))

	if RedisConfig.Addr == "" {
		utils.Logger.Fatal("REDIS_ADDR cannot be empty")
	} else if GoogleOauthConfig.ClientID == "" {
		utils.Logger.Fatal("GOOGLE_CLIENT_ID cannot be empty")
	} else if GoogleOauthConfig.ClientSecret == "" {
		utils.Logger.Fatal("GOOGLE_CLIENT_SECRET cannot be empty")
	} else if len(Config.JwtKey) == 0 {
		utils.Logger.Fatal("JWT_KEY cannot be empty")
	} else if Config.Host == "" {
		utils.Logger.Fatal("HOST cannot be empty")
	}

	if EmailConfig.SMTPServer == "" {
		utils.Logger.Error("SMTP_SERVER is empty. Some feature may not work...")
	} else if EmailConfig.EmailAddress == "" {
		utils.Logger.Error("EMAIL_ADDRESS is empty. Some feature may not work...")
	} else if EmailConfig.Password == "" {
		utils.Logger.Error("EMAIL_PASSWORD is empty. Some feature may not work...")
	}

	if Config.RapidAPIKey == "" {
		utils.Logger.Error("RAPIDAPI_KEY is empty. Some feature may not work...")
	}
}
