package core_kafka


import (
    "fmt"
    "time"

    "github.com/kelseyhightower/envconfig"
)

type Config struct {
    Brokers []string `envconfig:"BROKERS" required:"true" split_words:"true" default:"localhost:9092"`
    
  
    TopicPrefix string `envconfig:"TOPIC_PREFIX" default:"todoapp"`

    Producer ProducerConfig `ignored:"true"`
}

type ProducerConfig struct {
    // AckRequired - ждать подтверждения от брокера
    // true = надёжно, но медленнее
    // false = быстро, но можно потерять сообщения
    AckRequired bool `envconfig:"ACK_REQUIRED" default:"true"`
    
    // MaxRetries - сколько раз повторять отправку при ошибке
    MaxRetries int `envconfig:"MAX_RETRIES" default:"3"`
    
    // RetryBackoff - задержка между повторными попытками
    RetryBackoff time.Duration `envconfig:"RETRY_BACKOFF" default:"100ms"`
}


func NewConfig() (Config, error) {
    var cfg Config
    cfg.Producer = ProducerConfig{
        AckRequired:  true,
        MaxRetries:   3,
        RetryBackoff: 100 * time.Millisecond,
    }
    
    if err := envconfig.Process("KAFKA", &cfg); err != nil {
        return Config{}, fmt.Errorf("process kafka config: %w", err)
    }
    
    return cfg, nil
}

func NewConfigMust() Config {
    cfg, err := NewConfig()
    if err != nil {
        panic(fmt.Errorf("get kafka config: %w", err))
    }
    return cfg
}