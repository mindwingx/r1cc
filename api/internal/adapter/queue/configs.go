package queue

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func kafkaTopicsConfigs() []kafka.TopicSpecification {
	return []kafka.TopicSpecification{
		{
			Topic:             ProdTopic,
			NumPartitions:     50,
			ReplicationFactor: 1,
			Config: map[string]string{
				"retention.bytes":   "107374182400", // ~100GB
				"segment.bytes":     "536870912",    // 512MB
				"max.message.bytes": "1000000",      // 1MB
			},
		},
		{
			Topic:             ExpressTopic,
			NumPartitions:     30,
			ReplicationFactor: 1,
			Config: map[string]string{
				"retention.bytes": "53687091200", // 50GB
				"segment.bytes":   "536870912",   // 512MB
			},
		},
		{
			Topic:             RetryTopic,
			NumPartitions:     20,
			ReplicationFactor: 1,
			Config: map[string]string{
				"retention.ms":  "86400000",  // 24h
				"segment.bytes": "536870912", // 512MB
			},
		},
		{
			Topic:             DlqTopic,
			NumPartitions:     10,
			ReplicationFactor: 1,
			Config: map[string]string{
				"retention.ms":  "2592000000", // 30d
				"segment.bytes": "536870912",  // 512MB
			},
		},
	}
}

func kafkaProducerConfig(q *queue) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": q.config.Host,
		"compression.type":  "lz4",   // high compression ratio
		"linger.ms":         5,       // batch waiting time
		"batch.size":        1000000, // 1mb batch size
		"acks":              "1",     // balanced durability/throughput
		"max.in.flight":     5,       // parallel requests
		"message.max.bytes": 1000000, // 1mb max message size
		"retries":           10,      // retry failed deliveries
		"retry.backoff.ms":  1000,    // wait 1s between retries
	}
}

func kafkaConsumerConfig(q *queue, gp string) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":         q.config.Host,
		"group.id":                  fmt.Sprintf("%s-group", gp),
		"auto.offset.reset":         "latest", // start from latest offset
		"fetch.min.bytes":           500000,   // 500kb min fetch size
		"fetch.max.bytes":           5000000,  // 5mb max fetch size
		"max.partition.fetch.bytes": 1000000,  // 1mb per partition
		"fetch.wait.max.ms":         500,      // reduce wait time
		"enable.auto.commit":        "true",   // auto commit offsets
		"auto.commit.interval.ms":   5000,     // commit every 5s
	}
}
