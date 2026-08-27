package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `env:"ENV" env-default:"local"`
	DB       DBConfig
	TokenTTL time.Duration `env:"TOKEN_TTL" env-default:"1h"`
	GRPC     GRPCConfig
}

// LogValue keeps Config safe to log directly (e.g. slog.Any("cfg", cfg)) by
// routing through DBConfig's redacted LogValue instead of a raw JSON dump.
func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("env", c.Env),
		slog.Any("db", c.DB),
		slog.Duration("token_ttl", c.TokenTTL),
		slog.Int("grpc_port", c.GRPC.Port),
		slog.Duration("grpc_timeout", c.GRPC.Timeout),
	)
}

type DBConfig struct {
	Host     string `env:"DB__HOST" env-default:"localhost"`
	Port     int    `env:"DB__PORT" env-default:"5432"`
	User     string `env:"DB__USER" env-default:"sso_user"`
	Password string `env:"DB__PASSWORD" env-default:"sso_pass"`
	Name     string `env:"DB__NAME" env-default:"sso_db"`
	SSLMode  string `env:"DB__SSLMODE" env-default:"disable"`
}

// LogValue redacts the password so DBConfig can be logged safely via slog.
func (c DBConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", c.Host),
		slog.Int("port", c.Port),
		slog.String("user", c.User),
		slog.String("password", "***"),
		slog.String("name", c.Name),
		slog.String("sslmode", c.SSLMode),
	)
}

type GRPCConfig struct {
	Port    int           `env:"GRPC__PORT" env-default:"50051"`
	Timeout time.Duration `env:"GRPC__TIMEOUT" env-default:"10h"`
}

// DatabaseURL builds the Postgres DSN from the individual DB__* settings.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name, c.DB.SSLMode)
}

// MustLoad loads config from a .env file (if present) and the process
// environment. Environment variables already set take precedence over
// values from the .env file.
func MustLoad() *Config {
	_ = godotenv.Load()

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
