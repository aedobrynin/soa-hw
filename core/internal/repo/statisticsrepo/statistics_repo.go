package statisticsrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aedobrynin/soa-hw/core/internal/config"
	"github.com/aedobrynin/soa-hw/core/internal/model"
	"github.com/aedobrynin/soa-hw/core/internal/repo"

	"github.com/segmentio/kafka-go"
)

var _ repo.Statistics = &StatisticsRepo{}

type StatisticsRepo struct {
	postsLikesWriter *kafka.Writer
	postsViewsWriter *kafka.Writer
}

func (r *StatisticsRepo) PushPostView(ctx context.Context, view model.PostView) error {
	raw, err := json.Marshal(view)
	if err != nil {
		return err
	}

	err = r.postsViewsWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(view.PostID),
		Value: raw,
	})
	if err != nil {
		return fmt.Errorf("error on pushing post view to Kafka: %v", err)
	}
	return nil
}

func (r *StatisticsRepo) PushPostLike(ctx context.Context, like model.PostLike) error {
	raw, err := json.Marshal(like)
	if err != nil {
		return err
	}

	err = r.postsLikesWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(like.PostID),
		Value: raw,
	})
	if err != nil {
		return fmt.Errorf("error on pushing post like to Kafka: %v", err)
	}
	return nil
}

func (r *StatisticsRepo) Stop(ctx context.Context) {
	if err := r.postsLikesWriter.Close(); err != nil {
		log.Printf("error on postsLikesWriter.Close(): %v", err)
	}
	if err := r.postsViewsWriter.Close(); err != nil {
		log.Printf("error on postsViewsWriter.Close(): %v", err)
	}
}

func New(cfg *config.KafkaConfig) repo.Statistics {
	postsLikesWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   []string{cfg.BrokerAddr},
		Topic:     cfg.PostsLikesTopicName,
		Async:     false, // TODO: myb use?
		BatchSize: 0,     // TODO: myb use?
	})

	postsViewsWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   []string{cfg.BrokerAddr},
		Topic:     cfg.PostsViewsTopicName,
		Async:     false, // TODO: myb use?
		BatchSize: 0,     // TODO: myb use?
	})
	return &StatisticsRepo{
		postsLikesWriter: postsLikesWriter,
		postsViewsWriter: postsViewsWriter,
	}
}
