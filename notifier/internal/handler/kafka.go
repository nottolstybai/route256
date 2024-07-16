package handler

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"route256.ozon.ru/project/notifier/pkg/logger"
	"sync"
)

type KafkaConsumerGroup struct {
	sarama.ConsumerGroup
	bootstrapServer []string
	topics          []string
	id              string
	handler         sarama.ConsumerGroupHandler
}

func NewKafkaConsumerGroup(bootstrapServer []string, topics []string, id string, opts ...Option) (*KafkaConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Offsets.AutoCommit.Enable = false
	cg, err := sarama.NewConsumerGroup(bootstrapServer, id, config)
	if err != nil {
		return nil, fmt.Errorf("cant create consumer group: %w", err)
	}
	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &KafkaConsumerGroup{
		bootstrapServer: bootstrapServer,
		topics:          topics,
		id:              id,
		ConsumerGroup:   cg,
		handler:         NewConsumerGroupHandler(),
	}, nil
}

func (k *KafkaConsumerGroup) StartConsume(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("consumer-group run")
		for {
			select {
			case <-ctx.Done():
				logger.Info("consumer-group: ctx closed", zap.Error(ctx.Err()))
				return
			default:
				if err := k.Consume(ctx, k.topics, k.handler); err != nil {
					logger.Error("Error from consume", zap.Error(err))
				}
			}
		}
	}()
}
