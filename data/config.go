package data

//Configuration app configuration
type Configuration struct {
	RedisConfig       *RedisConfiguration
	GoogleOauthConfig *GoogleOauthConfiguration

	Host   string
	JwtKey string
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
