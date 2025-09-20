package queue

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"microservice/config"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/orm"
	"microservice/internal/adapter/provider/sms"
	"microservice/internal/adapter/registry"
	"microservice/internal/adapter/trace"
	"microservice/internal/modules/port"
	"microservice/pkg/meta"
	"microservice/pkg/meta/status"
	"microservice/pkg/utils"
	"time"
)

type (
	QFx struct {
		fx.In
		Logger      logger.ILogger
		Trace       trace.ITracer
		SmsProvider sms.ISmsProvider
		Sql         orm.ISqlTx
		Message     port.IMessageRepository
		Outbox      port.IOutboxRepository
	}
	queue struct {
		config    config.Queue
		producer  *kafka.Producer
		consumers map[string]*kafka.Consumer
		lgr       logger.ILogger
		trc       trace.ITracer
		sms       sms.ISmsProvider
		sql       orm.ISqlTx
		message   port.IMessageRepository
		outbox    port.IOutboxRepository
	}
)

func New(registry registry.IRegistry) IQueue {
	q := new(queue)
	if err := registry.Parse(&q.config); err != nil {
		utils.PrintStd(utils.StdPanic, "queue", "config parse err: %s", err)
	}

	return q
}

func (q *queue) Init() {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers":           q.config.Host,
		"socket.timeout.ms":           30000, // 30s timeout for network operations
		"metadata.request.timeout.ms": 10000, // 10s timeout for metadata fetch
		//"debug":                       "broker,admin", // Enable debug logs for connection issue
	})

	if err != nil {
		utils.PrintStd(utils.StdPanic, "queue", "kafka admin init err: %s", err)
	}

	//

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	metadata, err := admin.GetMetadata(nil, false, int(10*time.Second.Milliseconds()))
	if err != nil {
		utils.PrintStd(utils.StdPanic, "queue", "failed to fetch metadata: %s", err)
		return
	}

	errCount := 0
	newTopicsConf := make([]kafka.TopicSpecification, 0)
	topics := make([]kafka.TopicResult, 0)

	if len(metadata.Topics) > 0 {
		existingTopics := make(map[string]bool)

		for _, t := range metadata.Topics {
			existingTopics[t.Topic] = true
		}

		for _, t := range kafkaTopicsConfigs() {
			b := existingTopics[t.Topic]
			if b == false {
				newTopicsConf = append(newTopicsConf, t)
			}
		}
	} else {
		newTopicsConf = kafkaTopicsConfigs()
	}

	if len(newTopicsConf) > 0 {
		topics, err = admin.CreateTopics(ctx, newTopicsConf)
		if err != nil {
			utils.PrintStd(utils.StdPanic, "queue", "admin create topics err: %s", err)
		}

		//

		for _, t := range topics {
			if t.Error.Code() != kafka.ErrNoError {
				utils.PrintStd(utils.StdLog, "queue", "failed to create topic %s: %s", t.Topic, t.Error.Error())
				errCount++
			}
		}

		if errCount > 0 {
			utils.PrintStd(utils.StdPanic, "queue", "admin create topics failed")
		}
	}

	admin.Close()

	//

	time.Sleep(2 * time.Second)

	producer, err := kafka.NewProducer(kafkaProducerConfig(q))
	if err != nil {
		utils.PrintStd(utils.StdPanic, "queue", "producer init err: %s", err)
	}

	q.producer = producer

	//

	q.consumers = make(map[string]*kafka.Consumer)

	for _, gp := range q.config.Topics {
		consumer, cErr := kafka.NewConsumer(kafkaConsumerConfig(q, gp))
		if cErr != nil {
			utils.PrintStd(utils.StdLog, "queue", "failed to create consumer: %s", cErr.Error())
			errCount++
			break
		}

		q.consumers[gp] = consumer
	}

	if errCount > 0 {
		utils.PrintStd(utils.StdPanic, "queue", "consumer init failed")
	}
}

func (q *queue) Produce(ctx context.Context, topic, key string, value []byte) (err error) {
	sp, _ := q.trc.SpanByCtx(ctx, "queue.produce", "adapter.client")
	defer sp.End()

	txErr := q.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, nil)

	if txErr != nil {
		q.lgr.Error("queue.produce", zap.String(topic, key), zap.Error(txErr))
		sp.RecordError(txErr)

		err = meta.ServiceErr(status.Failed, txErr)
		return
	}

	q.producer.Flush(q.config.FlushTtl)
	return
}

func (q *queue) Fx(lc fx.Lifecycle, qfx QFx) IQueue {
	topicHdl := make(chan struct{})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			q.lgr = qfx.Logger
			q.trc = qfx.Trace
			q.sql = qfx.Sql
			q.sms = qfx.SmsProvider
			q.message = qfx.Message
			q.outbox = qfx.Outbox

			utils.PrintStd(utils.StdLog, "queue", "initiated")

			if c, ok := q.consumers[ProdTopic]; ok && c != nil {
				go q.productionTopicConsumer(ctx, topicHdl)
			}

			if c, ok := q.consumers[ExpressTopic]; ok && c != nil {
				go q.expressTopicConsumer(ctx, topicHdl)
			}

			if c, ok := q.consumers[RetryTopic]; ok && c != nil {
				go q.retryTopicConsumer(ctx, topicHdl)
			}

			if c, ok := q.consumers[DlqTopic]; ok && c != nil {
				go q.dlqTopicConsumer(ctx, topicHdl)
			}

			return
		},
		OnStop: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "queue", "stopping...")

			close(topicHdl)

			if q.producer != nil {
				q.producer.Flush(q.config.FlushTtl)
				q.producer.Close()
			}

			utils.PrintStd(utils.StdLog, "queue", "stopped")
			return
		},
	})

	return q
}
