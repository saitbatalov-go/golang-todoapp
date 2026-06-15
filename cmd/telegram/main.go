package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// TaskCreatedEvent - структура события (должна совпадать с той, что отправляет producer)
type TaskCreatedEvent struct {
	TaskID       string    `json:"task_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	AuthorUserID string    `json:"author_user_id"`
	AuthorName   string    `json:"author_name"`
	Completed    bool      `json:"completed"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetKey - возвращает ключ для партиционирования
func (e TaskCreatedEvent) GetKey() string {
	return e.TaskID
}

// GetType - возвращает тип события
func (e TaskCreatedEvent) GetType() string {
	return "task.created"
}

// TelegramMessage - структура запроса к Telegram API
type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// TelegramConsumer - потребитель сообщений из Kafka
// Забирает события и отправляет их в Telegram
type TelegramConsumer struct {
	// Kafka consumer
	consumer sarama.Consumer

	// Настройки
	topic    string // имя топика (очереди)
	botToken string // токен Telegram бота
	chatID   string // ID чата куда отправлять

	// HTTP клиент для отправки в Telegram
	httpClient *http.Client

	logger *zap.Logger
}

func NewTelegramConsumer(
	brokers []string,
	topic string,
	botToken string,
	chatID string,
	logger *zap.Logger,
) (*TelegramConsumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V3_0_0_0

	// Читаем сообщения с начала (OffsetOldest)
	// Если потребитель упал, при перезапуске прочитает всё, что пропустил
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true 	// Возвращать ошибки в канал

	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	// Создаём потребителя
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("create consumer: %w", err)
	}

	logger.Info("Telegram consumer created",
		zap.Strings("brokers", brokers),
		zap.String("topic", topic))

	return &TelegramConsumer{
		consumer:   consumer,
		topic:      topic,
		botToken:   botToken,
		chatID:     chatID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		logger:     logger,
	}, nil
}

// Start - запускает потребителя (блокируется)
func (c *TelegramConsumer) Start(ctx context.Context) error {
	// Получаем список всех партиций топика
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return fmt.Errorf("get partitions: %w", err)
	}

	c.logger.Info("starting Telegram consumer",
		zap.String("topic", c.topic),
		zap.Int32("partitions", int32(len(partitions))))

	// Для каждой партиции запускаем отдельную горутину
	// Это позволяет обрабатывать сообщения параллельно
	for _, partition := range partitions {
		go c.consumePartition(ctx, partition)
	}

	// Ждём сигнала завершения
	<-ctx.Done()
	c.logger.Info("shutting down consumer...")

	return nil
}

// consumePartition - обрабатывает сообщения из одной партиции
func (c *TelegramConsumer) consumePartition(ctx context.Context, partition int32) {
	// Создаём потребителя для конкретной партиции
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetOldest)
	if err != nil {
		c.logger.Error("failed to consume partition",
			zap.Int32("partition", partition),
			zap.Error(err))
		return
	}
	defer partitionConsumer.Close()

	c.logger.Info("started partition consumer",
		zap.Int32("partition", partition))

	// Бесконечный цикл обработки сообщений
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			// Получили новое сообщение!
			c.handleMessage(msg)

		case err := <-partitionConsumer.Errors():
			// Ошибка при получении
			c.logger.Error("consumer error", zap.Error(err))

		case <-ctx.Done():
			// Пришёл сигнал завершения
			c.logger.Info("stopping partition consumer",
				zap.Int32("partition", partition))
			return
		}
	}
}

// handleMessage - обрабатывает одно сообщение из Kafka
func (c *TelegramConsumer) handleMessage(msg *sarama.ConsumerMessage) {
	c.logger.Debug("received message",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("key", string(msg.Key)))

	// 1. Парсим JSON в структуру события
	var event TaskCreatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		c.logger.Error("failed to unmarshal event",
			zap.Error(err),
			zap.ByteString("raw", msg.Value))
		return
	}

	// 2. Формируем красивое сообщение для Telegram
	text := c.formatTelegramMessage(event)

	// 3. Отправляем в Telegram
	if err := c.sendTelegramMessage(text); err != nil {
		c.logger.Error("failed to send Telegram message",
			zap.Error(err),
			zap.String("task_id", event.TaskID))
		return
	}

	c.logger.Info("message sent to Telegram",
		zap.String("task_id", event.TaskID),
		zap.String("title", event.Title))
}

// formatTelegramMessage - форматирует событие в текст для Telegram
func (c *TelegramConsumer) formatTelegramMessage(event TaskCreatedEvent) string {
	// Статус выполнения
	status := "🟡 In Progress"
	if event.Completed {
		status = "✅ Completed"
	}

	// Форматируем время
	createdAt := event.CreatedAt.Format("15:04:05 02.01.2006")

	// Строим сообщение с использованием MarkdownV2
	// Эмодзи не требуют экранирования, они безопасны
	message := fmt.Sprintf(
		"📋 *New Task Created!*\n\n"+
			"*Title:* %s\n"+
			"*Description:* %s\n"+
			"*Author:* %s\n"+
			"*Status:* %s\n"+
			"*Created:* %s\n\n"+
			"*ID:* `%s`",
		escapeMarkdown(event.Title),
		escapeMarkdown(event.Description),
		escapeMarkdown(event.AuthorName),
		status,
		createdAt,
		event.TaskID,
	)

	return message
}

// escapeMarkdown - экранирует спецсимволы для MarkdownV2
// Telegram требует экранировать: _ * [ ] ( ) ~ ` > # + - = | { } . !
func escapeMarkdown(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	result := text
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}
	return result
}

// sendTelegramMessage - отправляет сообщение в Telegram API
func (c *TelegramConsumer) sendTelegramMessage(text string) error {
	// Формируем URL
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.botToken)

	// Формируем тело запроса
	body := TelegramMessage{
		ChatID:    c.chatID,
		Text:      text,
		ParseMode: "MarkdownV2",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned %d", resp.StatusCode)
	}

	return nil
}

// Close - закрывает соединение с Kafka
func (c *TelegramConsumer) Close() error {
	return c.consumer.Close()
}

// ========== ТОЧКА ВХОДА ==========

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")

	if botToken == "" || chatID == "" {
		logger.Fatal("TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID are required")
	}

	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092" // по умолчанию
	}

	logger.Info("configuration loaded",
		zap.String("brokers", kafkaBrokers),
		zap.String("topic", "todoapp"),
		zap.String("chat_id", chatID))

	// 4. Создаём потребителя
	consumer, err := NewTelegramConsumer(
		[]string{kafkaBrokers},
		"todoapp", // имя топика
		botToken,
		chatID,
		logger,
	)
	if err != nil {
		logger.Fatal("failed to create consumer", zap.Error(err))
	}
	defer consumer.Close()

	// 5. Настраиваем graceful shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // kill
	)
	defer cancel()

	// 6. Запускаем потребителя (блокируется)
	if err := consumer.Start(ctx); err != nil {
		logger.Fatal("consumer failed", zap.Error(err))
	}

	logger.Info("consumer stopped")
}
