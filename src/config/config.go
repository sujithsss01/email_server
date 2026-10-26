package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	KafkaConfig     KafkaConfig     `envPrefix:"KAFKA_"`
	SMTPConfig      SMTPConfig      `envPrefix:"SMTP_"`
	CassandraConfig CassandraConfig `envPrefix:"CASSANDRA_"`
	Environment     string          `env:"ENVIRONMENT" envDefault:"local"`
}
type KafkaConfig struct {
	Url string `env:"BROKER_URL" envDefault:"localhost:9092"`
}
type SMTPConfig struct {
	Url string `env:"RECEIVER_URL" envDefault:"localhost:1025"`
}
type CassandraConfig struct {
	Url      string `env:"SERVER_URL" envDefault:"localhost:9042"`
	Keyspace string `env:"KEYSPACE" envDefault:"public"`
}

func (c *Config) isProduction() bool {
	return c.Environment == "production"
}

func initiateConfig() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load env files")
	}
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		panic("Failed to parse env file")
	}
	instance = &cfg
}

func GetConfig() *Config {
	once.Do(initiateConfig)
	return instance
}
