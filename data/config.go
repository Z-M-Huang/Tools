package data

//Configuration app configuration
type Configuration struct {
	DatabaseConfig    *DatabaseConfiguration
	RedisConfig       *RedisConfiguration
	GoogleOauthConfig *GoogleOauthConfiguration
	SitemapConfig     *SitemapConfiguration

	JwtKey          []byte
	Host            string
	ResourceVersion string

	IsDebug bool
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

//SitemapConfiguration sitemap generation config
type SitemapConfiguration struct {
	GenerateSitemap bool
	IsHTTPS         bool
}
