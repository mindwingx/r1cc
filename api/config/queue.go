package config

type Queue struct {
	Host            string   `mapstructure:"QUEUE_HOST"`
	Topics          []string `mapstructure:"QUEUE_TOPICS"`
	RetryDelay      int      `mapstructure:"QUEUE_RETRY_DELAY_SEC"`
	ConsumerReadTtl int      `mapstructure:"QUEUE_CONSUMER_READ_TTL_MS"`
	FlushTtl        int      `mapstructure:"QUEUE_PRODUCER_FLUSH_TTL_MS"`
}
