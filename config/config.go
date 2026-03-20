package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// add to Config struct
type Config struct {
	App       AppConfig
	HTTP      HTTPConfig
	Postgres  PostgresConfig
	Log       LogConfig
	JWT       JWTConfig
	Redis     RedisConfig
	RateLimit RateLimitConfig
}

type AppConfig struct {
	Name    string `env:"APP_NAME"    env-default:"todo-clean"`
	Version string `env:"APP_VERSION" env-default:"1.0.0"`
	Env     string `env:"APP_ENV"     env-default:"development"`
}

type HTTPConfig struct {
	Port string `env:"HTTP_PORT" env-default:"8080"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST"     env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT"     env-default:"5432"`
	User     string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"POSTGRES_DB"       env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE"  env-default:"disable"`
}

type LogConfig struct {
	Level string `env:"LOG_LEVEL" env-default:"info"`
}

type JWTConfig struct {
	Secret      string `env:"JWT_SECRET"       env-required:"true"`
	ExpiryHours int    `env:"JWT_EXPIRY_HOURS"  env-default:"24"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"     env-default:"localhost"`
	Port     string `env:"REDIS_PORT"     env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB"       env-default:"0"`
}

type RateLimitConfig struct {
	IPLimit   int `env:"RATE_LIMIT_IP"   env-default:"100"`
	UserLimit int `env:"RATE_LIMIT_USER" env-default:"1000"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		// .env file is optional — fall back to real env vars
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}
	return cfg, nil
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DBName, p.SSLMode,
	)
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
