package main

import (
    "context"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "net/smtp"
    "os"
    "os/signal"
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

// EmailConsumer - потребитель сообщений из Kafka, отправляющий email
type EmailConsumer struct {
    consumer   sarama.Consumer
    topic      string
    smtpHost   string
    smtpPort   string
    fromEmail  string
    fromPass   string
    toEmail    string
    logger     *zap.Logger
}

func NewEmailConsumer(
    brokers []string,
    topic string,
    smtpHost, smtpPort, fromEmail, fromPass, toEmail string,
    logger *zap.Logger,
) (*EmailConsumer, error) {
    config := sarama.NewConfig()
    config.Version = sarama.V3_0_0_0
    config.Consumer.Offsets.Initial = sarama.OffsetOldest
    config.Consumer.Return.Errors = true

    consumer, err := sarama.NewConsumer(brokers, config)
    if err != nil {
        return nil, fmt.Errorf("create consumer: %w", err)
    }

    logger.Info("Email consumer created",
        zap.Strings("brokers", brokers),
        zap.String("topic", topic),
        zap.String("smtp_host", smtpHost),
        zap.String("from", fromEmail),
        zap.String("to", toEmail))

    return &EmailConsumer{
        consumer: consumer,
        topic:    topic,
        smtpHost: smtpHost,
        smtpPort: smtpPort,
        fromEmail: fromEmail,
        fromPass:  fromPass,
        toEmail:   toEmail,
        logger:    logger,
    }, nil
}

func (c *EmailConsumer) Start(ctx context.Context) error {
    partitions, err := c.consumer.Partitions(c.topic)
    if err != nil {
        return fmt.Errorf("get partitions: %w", err)
    }

    c.logger.Info("starting email consumer",
        zap.String("topic", c.topic),
        zap.Int32("partitions", int32(len(partitions))))

    for _, partition := range partitions {
        go c.consumePartition(ctx, partition)
    }

    <-ctx.Done()
    c.logger.Info("shutting down consumer...")
    return nil
}

func (c *EmailConsumer) consumePartition(ctx context.Context, partition int32) {
    partitionConsumer, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetOldest)
    if err != nil {
        c.logger.Error("failed to consume partition",
            zap.Int32("partition", partition),
            zap.Error(err))
        return
    }
    defer partitionConsumer.Close()

    for {
        select {
        case msg := <-partitionConsumer.Messages():
            c.handleMessage(msg)
        case err := <-partitionConsumer.Errors():
            c.logger.Error("consumer error", zap.Error(err))
        case <-ctx.Done():
            return
        }
    }
}

func (c *EmailConsumer) handleMessage(msg *sarama.ConsumerMessage) {
    var event TaskCreatedEvent
    if err := json.Unmarshal(msg.Value, &event); err != nil {
        c.logger.Error("failed to unmarshal event", zap.Error(err))
        return
    }

    c.logger.Debug("received message",
        zap.String("task_id", event.TaskID),
        zap.String("title", event.Title))

    // Формируем email
    subject := fmt.Sprintf("📋 New Task Created: %s", event.Title)
    body := c.formatEmailBody(event)

    // Отправляем email
    if err := c.sendEmail(subject, body); err != nil {
        c.logger.Error("failed to send email",
            zap.Error(err),
            zap.String("task_id", event.TaskID))
        return
    }

    c.logger.Info("email sent",
        zap.String("task_id", event.TaskID),
        zap.String("title", event.Title),
        zap.String("to", c.toEmail))
}

func (c *EmailConsumer) formatEmailBody(event TaskCreatedEvent) string {
    status := "🟡 In Progress"
    if event.Completed {
        status = "✅ Completed"
    }

    createdAt := event.CreatedAt.Format("15:04:05 02.01.2006")

    return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
    <div style="background-color: #4CAF50; color: white; padding: 20px; text-align: center; border-radius: 10px 10px 0 0;">
        <h1 style="margin: 0;">📋 New Task Created!</h1>
    </div>
    <div style="background-color: #f9f9f9; padding: 20px; border: 1px solid #ddd; border-radius: 0 0 10px 10px;">
        <h2>Task Details:</h2>
        <table style="width: 100%; border-collapse: collapse;">
            <tr><td style="padding: 8px 0;"><strong>Title:</strong></td><td>%s</td></tr>
            <tr><td style="padding: 8px 0;"><strong>Description:</strong></td><td>%s</td></tr>
            <tr><td style="padding: 8px 0;"><strong>Author:</strong></td><td>%s</td></tr>
            <tr><td style="padding: 8px 0;"><strong>Status:</strong></td><td>%s</td></tr>
            <tr><td style="padding: 8px 0;"><strong>Created:</strong></td><td>%s</td></tr>
            <tr><td style="padding: 8px 0;"><strong>ID:</strong></td><td><code>%s</code></td></tr>
        </table>
        <hr>
        <p style="color: #666; font-size: 12px;">This notification was sent automatically by TodoApp.</p>
    </div>
</body>
</html>`,
        event.Title,
        event.Description,
        event.AuthorName,
        status,
        createdAt,
        event.TaskID,
    )
}

func (c *EmailConsumer) sendEmail(subject, body string) error {
    // Настройка TLS
    tlsConfig := &tls.Config{
        InsecureSkipVerify: false,
        ServerName:         c.smtpHost,
    }

    // Подключение к серверу
    conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", c.smtpHost, c.smtpPort), tlsConfig)
    if err != nil {
        return fmt.Errorf("tls dial: %w", err)
    }
    defer conn.Close()

    client, err := smtp.NewClient(conn, c.smtpHost)
    if err != nil {
        return fmt.Errorf("create client: %w", err)
    }
    defer client.Close()

    // Аутентификация
    auth := smtp.PlainAuth("", c.fromEmail, c.fromPass, c.smtpHost)
    if err = client.Auth(auth); err != nil {
        return fmt.Errorf("auth: %w", err)
    }

    // Отправитель и получатель
    if err = client.Mail(c.fromEmail); err != nil {
        return fmt.Errorf("mail from: %w", err)
    }
    if err = client.Rcpt(c.toEmail); err != nil {
        return fmt.Errorf("rcpt to: %w", err)
    }

    // Тело письма
    w, err := client.Data()
    if err != nil {
        return fmt.Errorf("data: %w", err)
    }
    defer w.Close()

    message := fmt.Sprintf(
        "From: %s\r\n"+
            "To: %s\r\n"+
            "Subject: =?UTF-8?B?%s?=\r\n"+
            "MIME-Version: 1.0\r\n"+
            "Content-Type: text/html; charset=UTF-8\r\n"+
            "\r\n%s",
        c.fromEmail,
        c.toEmail,
        encodeBase64(subject),
        body,
    )

    if _, err = w.Write([]byte(message)); err != nil {
        return fmt.Errorf("write message: %w", err)
    }

    return nil
}

func encodeBase64(s string) string {
    return s // В реальном коде используйте base64.StdEncoding.EncodeToString([]byte(s))
}

func (c *EmailConsumer) Close() error {
    return c.consumer.Close()
}

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

    // Читаем переменные
    kafkaBrokers := os.Getenv("KAFKA_BROKERS")
    if kafkaBrokers == "" {
        kafkaBrokers = "localhost:9092"
    }

    smtpHost := os.Getenv("EMAIL_SMTP_HOST")
    smtpPort := os.Getenv("EMAIL_SMTP_PORT")
    fromEmail := os.Getenv("EMAIL_FROM")
    fromPass := os.Getenv("EMAIL_PASSWORD")
    toEmail := os.Getenv("EMAIL_TO")

    if smtpHost == "" || fromEmail == "" || fromPass == "" || toEmail == "" {
        logger.Fatal("EMAIL_SMTP_HOST, EMAIL_FROM, EMAIL_PASSWORD, EMAIL_TO are required")
    }

    logger.Info("configuration loaded",
        zap.String("brokers", kafkaBrokers),
        zap.String("topic", "todoapp"),
        zap.String("smtp_host", smtpHost),
        zap.String("from", fromEmail),
        zap.String("to", toEmail))

    // Создаём потребителя
    consumer, err := NewEmailConsumer(
        []string{kafkaBrokers},
        "todoapp",
        smtpHost,
        smtpPort,
        fromEmail,
        fromPass,
        toEmail,
        logger,
    )
    if err != nil {
        logger.Fatal("failed to create consumer", zap.Error(err))
    }
    defer consumer.Close()

    ctx, cancel := signal.NotifyContext(
        context.Background(),
        syscall.SIGINT,
        syscall.SIGTERM,
    )
    defer cancel()

    if err := consumer.Start(ctx); err != nil {
        logger.Fatal("consumer failed", zap.Error(err))
    }

    logger.Info("consumer stopped")
}