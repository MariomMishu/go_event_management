package config

import "time"

type AppConfig struct {
	Name            string
	Port            string
	NumberOfWorkers int
}
type DbConfig struct {
	Host            string
	Port            string
	User            string
	Pass            string
	Schema          string
	MaxIdleConn     int
	MaxOpenConn     int
	MaxConnLifetime time.Duration
	Debug           bool
}
type RedisConfig struct {
	Host               string
	Port               string
	Pass               string
	Db                 int
	MandatoryPrefix    string
	AccessUuidPrefix   string
	RefreshUuidPrefix  string
	UserPrefix         string
	PermissionPrefix   string
	UserCacheTTL       time.Duration
	PermissionCacheTTL time.Duration
}
type JwtConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}
type EmailConfig struct {
	Url     string
	Timeout time.Duration
}
type Config struct {
	App   *AppConfig
	DB    *DbConfig
	Redis *RedisConfig
	Jwt   *JwtConfig
	Email *EmailConfig
}

var config Config

func GetAll() Config {
	return config
}
func App() *AppConfig {
	return config.App
}
func Db() *DbConfig {
	return config.DB
}
func Redis() *RedisConfig {
	return config.Redis
}
func Jwt() *JwtConfig {
	return config.Jwt
}
func Email() *EmailConfig { return config.Email }

func LoadConfig() {
	// Set defaults or load from env
	// config.App = &AppConfig{
	// 	Port: "8080", // or os.Getenv("PORT")
	// }
	setDefaultConfig()
}
func setDefaultConfig() {
	config.App = &AppConfig{
		Name:            "EMS",
		Port:            "8080",
		NumberOfWorkers: 5,
	}
	config.DB = &DbConfig{
		Host:            "127.0.0.1",
		Port:            "3306",
		User:            "mariom",
		Pass:            "password",
		Schema:          "event_management",
		MaxIdleConn:     1,
		MaxOpenConn:     2,
		MaxConnLifetime: 30,
		Debug:           true,
	}
	config.Redis = &RedisConfig{
		Host:            "127.0.0.1",
		Port:            "6379",
		Pass:            "",
		Db:              0,
		MandatoryPrefix: "event_managment_",
	}
	// Initialize JwtConfig here:
	config.Jwt = &JwtConfig{
		AccessTokenSecret:  "your-access-secret",
		RefreshTokenSecret: "your-refresh-secret",
		// expiry in seconds or whatever unit you prefer
		AccessTokenExpiry:  time.Hour,      // e.g. 1 hour
		RefreshTokenExpiry: 24 * time.Hour, // e.g. 24 hours
	}
}
