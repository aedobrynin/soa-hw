package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/aedobrynin/soa-hw/core/internal/clients/postsclient"
	"github.com/aedobrynin/soa-hw/core/internal/clients/statisticsclient"
	"github.com/aedobrynin/soa-hw/core/internal/httpadapter"
	"github.com/aedobrynin/soa-hw/core/internal/service"
)

const (
	// AppConfig defaults
	DefaultDebug           = false
	DefaultShutdownTimeout = 20 * time.Second

	// DatabaseConfig defaults
	DefaultDSN           = "dsn://"
	DefaultMigrationsDir = "file://postgresql/core/migrations/"

	// HTTP defaults (TODO: myb move to httpadapter module?)
	DefaultServeAddress         = "localhost:3000"
	DefaultBasePath             = "/"
	DefaultAccessTokenCookie    = "access_token"
	DefaultRefreshTokenCookie   = "refresh_token"
	DefaultSigningKey           = "qwerty"
	DefaultAccessTokenDuration  = 1 * time.Minute
	DefaultRefreshTokenDuration = 1 * time.Hour
	DefaultUseTLS               = false

	// PostsConfig defaults (TODO: myb move to postsclient module?)
	DefaultPostsAddr         = "posts_service:8080"
	DefaultPostsTimeout      = 5 * time.Second
	DefaultPostsRetriesCount = 3

	// StatisticsConfig defaults (TODO: myb move to statisticsclient module?)
	DefaultStatisticsAddr         = "statistics_service:8080"
	DefaultStatisticsTimeout      = 5 * time.Second
	DefaultStatisticsRetriesCount = 3

	// KafkaConfig defaults
	DefaultKafkaBrokerAddr          = "kafka:9092"
	DefaultKafkaPostsViewsTopicName = "posts_views"
	DefaultKafkaPostsLikesTopicName = "posts_likes"
)

type AppConfig struct {
	Debug           bool          `yaml:"debug"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type KafkaConfig struct {
	BrokerAddr          string
	PostsViewsTopicName string
	PostsLikesTopicName string
}

type Config struct {
	App      AppConfig          `yaml:"app"`
	Database DatabaseConfig     `yaml:"database"`
	HTTP     httpadapter.Config `yaml:"http"`
	Kafka    KafkaConfig        `yaml:"kafka"`

	Auth       service.AuthConfig                      `yaml:"auth"`
	Posts      postsclient.PostsClientConfig           `yaml:"posts_client"`
	Statistics statisticsclient.StatisticsClientConfig `yaml:"statistics_client"`
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
			DSN:           DefaultDSN,
			MigrationsDir: DefaultMigrationsDir,
		},
		HTTP: httpadapter.Config{
			ServeAddress:       DefaultServeAddress,
			BasePath:           DefaultBasePath,
			UseTLS:             DefaultUseTLS,
			AccessTokenCookie:  DefaultAccessTokenCookie,
			RefreshTokenCookie: DefaultRefreshTokenCookie,
		},
		Kafka: KafkaConfig{
			BrokerAddr:          DefaultKafkaBrokerAddr,
			PostsViewsTopicName: DefaultKafkaPostsViewsTopicName,
			PostsLikesTopicName: DefaultKafkaPostsLikesTopicName,
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
		Statistics: statisticsclient.StatisticsClientConfig{
			Addr:         DefaultStatisticsAddr,
			Timeout:      DefaultStatisticsTimeout,
			RetriesCount: DefaultStatisticsRetriesCount,
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
