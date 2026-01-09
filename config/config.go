package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DBString                    string `env:"DB_STRING" env-default:""`
	RabbitMQUser                string `env:"RABBIMQ_USER" env-default:""`
	RabbitMQPassword            string `env:"RABBIMQ_PASSWORD" env-default:""`
	RabbitMQHost                string `env:"RABBIMQ_HOST" env-default:""`
	RabbitMQPort                int    `env:"RABBIMQ_PORT" env-default:""`
	RabbitMQQueueName           string `env:"RABBIMQ_QUEUE_NAME" env-default:""`
	RabbitMQConsumeWorkersCount int    `env:"RABBIMQ_CONSUME_WORKERS_COUNT" env-default:""`
}

func Load() (*Config, error) {
	var config Config

	if err := cleanenv.ReadConfig(".env", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
