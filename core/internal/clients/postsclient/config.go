package postsclient

import "time"

type PostsClientConfig struct {
	Addr         string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retries_count"`
}
