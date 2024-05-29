package statisticsclient

import "time"

type StatisticsClientConfig struct {
	Addr         string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retries_count"`
}
