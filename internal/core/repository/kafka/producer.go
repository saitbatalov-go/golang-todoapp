package core_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	"go.uber.org/zap"
)


type MessageProducer struct {
    // producer - внутренний клиент Sarama
    producer sarama.SyncProducer
    
    // topic - имя топика (очереди) куда отправляем сообщения
    topic string
    
    logger *core_logger.Logger
}

func NewMessageProducer(ctx context.Context, cfg Config, logger *core_logger.Logger) (*MessageProducer, error) {
    saramaConfig := sarama.NewConfig()
    saramaConfig.Version = sarama.V3_0_0_0 
    
    if cfg.Producer.AckRequired {
        // Ждём подтверждения от всех реплик (максимальная надёжность)
        saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
    } else {
        // Не ждём подтверждения (максимальная скорость)
        saramaConfig.Producer.RequiredAcks = sarama.NoResponse
    }
    
    // Количество повторных попыток
    saramaConfig.Producer.Retry.Max = cfg.Producer.MaxRetries
    
    // Ждать ответа от брокера
    saramaConfig.Producer.Return.Successes = true
    
    // Сжимаем сообщения для экономии места
    saramaConfig.Producer.Compression = sarama.CompressionSnappy
    
    // Создаём синхронного производителя (ждёт подтверждения)
    producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
    if err != nil {
        return nil, fmt.Errorf("create producer: %w", err)
    }
    
    logger.Info("Kafka producer created",
        zap.Strings("brokers", cfg.Brokers),
        zap.String("topic", cfg.TopicPrefix))
    
    return &MessageProducer{
        producer: producer,
        topic:    cfg.TopicPrefix,
        logger:   logger,
    }, nil
}

func (p *MessageProducer) Publish(ctx context.Context, key string, event Event) error {
    start := time.Now()
    
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("marshal event: %w", err)
    }
    
    msg := &sarama.ProducerMessage{
        Topic: p.topic,                  // имя очереди
        Key:   sarama.StringEncoder(key), // ключ (для партиционирования)
        Value: sarama.ByteEncoder(data),  // тело сообщения (JSON)
        Headers: []sarama.RecordHeader{
            {Key: []byte("event_type"), Value: []byte(event.GetType())},
            {Key: []byte("timestamp"), Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
        },
    }
    
    partition, offset, err := p.producer.SendMessage(msg)
    if err != nil {
        return fmt.Errorf("send message: %w", err)
    }
    
    elapsed := time.Since(start)
    p.logger.Debug("message published",
        zap.String("topic", p.topic),
        zap.String("key", key),
        zap.String("event_type", event.GetType()),
        zap.Int32("partition", partition),
        zap.Int64("offset", offset),
        zap.Duration("elapsed", elapsed))
    
    return nil
}

func (p *MessageProducer) Close() error {
    return p.producer.Close()
}