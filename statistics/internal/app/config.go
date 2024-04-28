package app

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/aedobrynin/soa-hw/statistics/internal/httpadapter"
)

const (
	DefaultServeAddress    = "localhost:3000"
	DefaultShutdownTimeout = 20 * time.Second
	DefaultBasePath        = "/"
	DefaultDbAddr          = "statistics_clickhouse:9440"
	DefaultDbDatabase      = "statistics"
	DefaultDbUser          = "statistics"
	DefaultDbPassword      = "statistics"
	DefaultDbMigrationsDir = "file://clickhouse/statistics/migrations/"
)

type AppConfig struct {
	Debug           bool          `yaml:"debug"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Addr          string `yaml:"addr"`
	Database      string `yaml:"database"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type Config struct {
	App      AppConfig          `yaml:"app"`
	Database DatabaseConfig     `yaml:"database"`
	HTTP     httpadapter.Config `yaml:"http"`
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
			Addr:          DefaultDbAddr,
			Database:      DefaultDbDatabase,
			User:          DefaultDbUser,
			Password:      DefaultDbPassword,
			MigrationsDir: DefaultDbMigrationsDir,
		},
		HTTP: httpadapter.Config{
			ServeAddress: DefaultServeAddress,
			BasePath:     DefaultBasePath,
			UseTLS:       false,
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
