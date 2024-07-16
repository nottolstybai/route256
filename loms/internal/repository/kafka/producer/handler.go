package producer

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"route256.ozon.ru/project/loms/internal/entity"
	"route256.ozon.ru/project/loms/pkg/logger"
	"time"
)

// Handler type that will handle sending messages to kafka
type Handler struct {
	producer sarama.SyncProducer
	topic    string
}

func NewHandler(producer sarama.SyncProducer, topic string) *Handler {
	return &Handler{producer: producer, topic: topic}
}

// SendMessage sends event to specified topic
func (h *Handler) SendMessage(event entity.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: h.topic,
		//Key:   sarama.StringEncoder(strconv.Itoa(int(event.ID))),
		Value: sarama.ByteEncoder(bytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("loms"),
				Value: []byte("sync-producer"),
			},
		},
		Timestamp: time.Now(),
	}

	partition, offset, err := h.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	logger.Info("msg sent",
		zap.Int32("partition", partition),
		zap.Int32("eventID", event.ID),
		zap.Int64("offset", offset))
	return nil
}

func (h *Handler) CloseConnections() error {
	err := h.producer.Close()
	if err != nil {
		return err
	}
	return nil
}
