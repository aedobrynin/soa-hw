package app

import (
	"os"

	"github.com/aedobrynin/soa-hw/posts/internal/grpcadapter"
	"gopkg.in/yaml.v3"
)

const (
	DefaultDebug         = false
	DefaultDSN           = "dsn://"
	DefaultMigrationsDir = "file://postgresql/posts/migrations/"
	DefaultGRPCPort      = 8080
)

type AppConfig struct {
	Debug bool `yaml:"debug"`
}

type DatabaseConfig struct {
	DSN           string `yaml:"dsn"`
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
			Debug: DefaultDebug,
		},
		Database: DatabaseConfig{
			DSN:           DefaultDSN,
			MigrationsDir: DefaultMigrationsDir,
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
