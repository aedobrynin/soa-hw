package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type (
	ClickHouseContainer struct {
		testcontainers.Container
		Config ClickHouseContainerConfig
	}

	ClickHouseContainerOption func(c *ClickHouseContainerConfig)

	ClickHouseContainerConfig struct {
		ImageTag   string
		User       string
		Password   string
		MappedPort string
		Database   string
		Host       string
	}
)

func (c ClickHouseContainer) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Config.Host, c.Config.MappedPort)
}

func NewClickHouseContainer(ctx context.Context, opts ...ClickHouseContainerOption) (*ClickHouseContainer, error) {
	const (
		clickHouseImage = "clickhouse/clickhouse-server"
		clickHousePort  = "9000"
	)

	config := ClickHouseContainerConfig{
		ImageTag: "24.4",
		User:     "statistics",
		Password: "statistics",
		Database: "statistics",
	}
	for _, opt := range opts {
		opt(&config)
	}

	containerPort := clickHousePort + "/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Env: map[string]string{
				"CLICKHOUSE_DB":       config.Database,
				"CLICKHOUSE_USER":     config.User,
				"CLICKHOUSE_PASSWORD": config.Password,
			},
			ExposedPorts: []string{
				containerPort,
				"8123/tcp",
			},
			Image: fmt.Sprintf("%s:%s", clickHouseImage, config.ImageTag),
			WaitingFor: wait.ForHTTP("/ping").WithPort("8123/tcp").WithStatusCodeMatcher(
				func(status int) bool {
					return status == http.StatusOK
				},
			),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("getting request provider: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting host for: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return nil, fmt.Errorf("getting mapped port for (%s): %w", containerPort, err)
	}
	config.MappedPort = mappedPort.Port()
	config.Host = host

	fmt.Println("Host:", config.Host, config.MappedPort)

	return &ClickHouseContainer{
		Container: container,
		Config:    config,
	}, nil
}
