package app

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/aedobrynin/soa-hw/core/internal/clients/postsclient"
	"github.com/aedobrynin/soa-hw/core/internal/httpadapter"
	"github.com/aedobrynin/soa-hw/core/internal/service"
)

const (
	DefaultServeAddress         = "localhost:3000"
	DefaultShutdownTimeout      = 20 * time.Second
	DefaultBasePath             = "/"
	DefaultAccessTokenCookie    = "access_token"
	DefaultRefreshTokenCookie   = "refresh_token"
	DefaultSigningKey           = "qwerty"
	DefaultAccessTokenDuration  = 1 * time.Minute
	DefaultRefreshTokenDuration = 1 * time.Hour
	DefaultDSN                  = "dsn://"
	DefaultMigrationsDir        = "file://postgresql/core/migrations/"
	DefaultPostsAddr            = "posts_service:8080"
	DefaultPostsTimeout         = 5 * time.Second
	DefaultPostsRetriesCount    = 3
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

	Auth  service.AuthConfig            `yaml:"auth"`
	Posts postsclient.PostsClientConfig `yaml:"posts_client"`
}

func NewConfig(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cnf := Config{
		App: AppConfig{
			Debug:           false,
			ShutdownTimeout: DefaultShutdownTimeout,
		},
		Database: DatabaseConfig{
			DSN:           DefaultDSN,
			MigrationsDir: DefaultMigrationsDir,
		},
		HTTP: httpadapter.Config{
			ServeAddress:       DefaultServeAddress,
			BasePath:           DefaultBasePath,
			UseTLS:             false,
			AccessTokenCookie:  DefaultAccessTokenCookie,
			RefreshTokenCookie: DefaultRefreshTokenCookie,
		},
		Auth: service.AuthConfig{
			SigningKey:           DefaultSigningKey,
			AccessTokenDuration:  DefaultAccessTokenDuration,
			RefreshTokenDuration: DefaultRefreshTokenDuration,
		},
		Posts: postsclient.PostsClientConfig{
			Addr:         DefaultPostsAddr,
			Timeout:      DefaultPostsTimeout,
			RetriesCount: DefaultPostsRetriesCount,
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
