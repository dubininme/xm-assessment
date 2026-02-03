package config

import (
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Port            string `envconfig:"PORT" default:"8080"`
	Db              DbConfig
	Kafka           KafkaConfig
	Outbox          OutboxConfig
	ShutdownTimeout int    `envconfig:"SHUTDOWN_TIMEOUT" default:"5"`
	JWTSecret       string `envconfig:"JWT_SECRET"`
}

type DbConfig struct {
	DBHost            string `envconfig:"DB_HOST"`
	DBPort            string `envconfig:"DB_PORT"`
	DBUser            string `envconfig:"DB_USER"`
	DBPassword        string `envconfig:"DB_PASSWORD"`
	DBName            string `envconfig:"DB_NAME"`
	DBMaxOpenConns    int    `envconfig:"DB_MAX_OPEN_CONNS" default:"10"`
	DBMaxIdleConns    int    `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	DBConnMaxLifetime int    `envconfig:"DB_CONN_MAX_LIFETIME" default:"300"`
}

type KafkaConfig struct {
	Brokers string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
	Topic   string `envconfig:"KAFKA_TOPIC" default:"company-events"`
}

func (k *KafkaConfig) BrokersList() []string {
	return strings.Split(k.Brokers, ",")
}

type OutboxConfig struct {
	BatchSize int           `envconfig:"OUTBOX_BATCH_SIZE" default:"100"`
	Interval  time.Duration `envconfig:"OUTBOX_INTERVAL" default:"5s"`
}

func InitConfig() (*AppConfig, error) {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
