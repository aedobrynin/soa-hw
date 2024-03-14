package app

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"core/internal/httpadapter"
	"core/internal/service"
)

const (
	AppName                     = "auth"
	DefaultServeAddress         = "localhost:3000"
	DefaultShutdownTimeout      = 20 * time.Second
	DefaultBasePath             = "/"
	DefaultAccessTokenCookie    = "access_token"
	DefaultRefreshTokenCookie   = "refresh_token"
	DefaultSigningKey           = "qwerty"
	DefaultAccessTokenDuration  = 1 * time.Minute
	DefaultRefreshTokenDuration = 1 * time.Hour
	DefaultDSN                  = "dsn://"
	DefaultMigrationsDir        = "file://migrations/auth"
)

type AppConfig struct {
	Debug           bool          `yaml:"debug"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type Config struct {
	App      AppConfig          `yaml:"app"`
	Database DatabaseConfig     `yaml:"database"`
	HTTP     httpadapter.Config `yaml:"http"`

	Auth service.AuthConfig `yaml:"auth"`
}

func NewConfig(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cnf := Config{
		App: AppConfig{
			ShutdownTimeout: DefaultShutdownTimeout,
		},
		Database: DatabaseConfig{
			DSN:           DefaultDSN,
			MigrationsDir: DefaultMigrationsDir,
		},
		HTTP: httpadapter.Config{
			ServeAddress:       DefaultServeAddress,
			BasePath:           DefaultBasePath,
			AccessTokenCookie:  DefaultAccessTokenCookie,
			RefreshTokenCookie: DefaultRefreshTokenCookie,
		},
		Auth: service.AuthConfig{
			SigningKey:           DefaultSigningKey,
			AccessTokenDuration:  DefaultAccessTokenDuration,
			RefreshTokenDuration: DefaultRefreshTokenDuration,
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
