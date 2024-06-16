package app

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/aedobrynin/soa-hw/statistics/internal/grpcadapter"
)

const (
	DefaultShutdownTimeout = 20 * time.Second
	DefaultDebug           = true
	DefaultBasePath        = "/"
	DefaultDbAddr          = "statistics_clickhouse:9440"
	DefaultDbDatabase      = "statistics"
	DefaultDbUser          = "statistics"
	DefaultDbPassword      = "statistics"
	DefaultDbMigrationsDir = "file://clickhouse/statistics/migrations/"
	DefaultGRPCPort        = 8080
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
	App      AppConfig              `yaml:"app"`
	Database DatabaseConfig         `yaml:"database"`
	GRPC     grpcadapter.GRPCConfig `yaml:"grpc"`
}

func NewConfig(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cnf := Config{
		App: AppConfig{
			Debug:           DefaultDebug,
			ShutdownTimeout: DefaultShutdownTimeout,
		},
		Database: DatabaseConfig{
			Addr:          DefaultDbAddr,
			Database:      DefaultDbDatabase,
			User:          DefaultDbUser,
			Password:      DefaultDbPassword,
			MigrationsDir: DefaultDbMigrationsDir,
		},
		GRPC: grpcadapter.GRPCConfig{
			Port: DefaultGRPCPort,
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
