package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
	"microservice/internal/domain"
	"microservice/pkg/utils"
	"time"
)

func (q *queue) productionTopicConsumer(ctx context.Context, handler chan struct{}) {
	c := q.consumers[ProdTopic]

	if err := c.Subscribe(ProdTopic, nil); err != nil {
		q.lgr.Error("queue.consumer.subscribe.prod", zap.Error(err))

		utils.PrintStd(utils.StdPanic, "queue.consumer.subscribe.prod: %s", err.Error())

		// todo: set alert with prometheus
		// todo: handle the db try counter and send to retry
	}

	for {
		select {
		case <-handler:
			if err := c.Close(); err != nil {
				q.lgr.Error("queue.consumer.close", zap.String("topic", ProdTopic), zap.Error(err))
				return
			}

			q.lgr.Info("queue.consumer.close", zap.String("topic", ProdTopic))
			return
		default:
			msg, err := c.ReadMessage(time.Duration(q.config.ConsumerReadTtl))
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok == true && kafkaErr.Code() != kafka.ErrTimedOut {
					q.lgr.Error("queue.consumer.prod.read", zap.Error(err))
					// todo: set alert with prometheus
					break
				}
			}

			if msg != nil {
				err = q.consumeAndSendMessage(ctx, ProdTopic, msg)
				if err = q.Produce(ctx, RetryTopic, string(msg.Key), msg.Value); err != nil {
					//todo: set grafana/prometheus alarm
					q.lgr.Error("queue.consumer.provider.reproduce",
						zap.String("topic", ProdTopic),
						zap.String("message.id", string(msg.Key)),
						zap.Error(err),
					)
				}
			}
		}
	}
}

func (q *queue) expressTopicConsumer(ctx context.Context, handler chan struct{}) {
	c := q.consumers[ExpressTopic]

	if err := c.Subscribe(ExpressTopic, nil); err != nil {
		q.lgr.Error("queue.consumer.subscribe.express", zap.Error(err))

		utils.PrintStd(utils.StdPanic, "queue.consumer.subscribe.express: %s", err.Error())

		// todo: set alert with prometheus
		// todo: handle the db try counter and send to retry
	}

	for {
		select {
		case <-handler:
			if err := c.Close(); err != nil {
				q.lgr.Error("queue.consumer.close", zap.String("topic", ExpressTopic), zap.Error(err))
				return
			}

			q.lgr.Info("queue.consumer.close", zap.String("topic", ExpressTopic))
			return
		default:
			msg, err := c.ReadMessage(time.Duration(q.config.ConsumerReadTtl))
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok == true && kafkaErr.Code() != kafka.ErrTimedOut {
					q.lgr.Error("queue.consumer.prod.read", zap.Error(err))
					// todo: set alert with prometheus
					break
				}
			}

			if msg != nil {
				err = q.consumeAndSendMessage(ctx, ExpressTopic, msg)
				if err = q.Produce(ctx, RetryTopic, string(msg.Key), msg.Value); err != nil {
					//todo: set grafana/prometheus alarm
					q.lgr.Error("queue.consumer.provider.reproduce",
						zap.String("topic", ExpressTopic),
						zap.String("message.id", string(msg.Key)),
						zap.Error(err),
					)
				}
			}
		}
	}
}

func (q *queue) retryTopicConsumer(ctx context.Context, handler chan struct{}) {
	c := q.consumers[RetryTopic]

	if err := c.Subscribe(RetryTopic, nil); err != nil {
		q.lgr.Error("queue.consumer.subscribe.retry", zap.Error(err))

		utils.PrintStd(utils.StdPanic, "queue.consumer.subscribe.retry: %s", err.Error())

		// todo: set alert with prometheus
		// todo: handle the db try counter and send to retry
	}
	for {
		select {
		case <-handler:
			if err := c.Close(); err != nil {
				q.lgr.Error("queue.consumer.close", zap.String("topic", RetryTopic), zap.Error(err))
				return
			}

			q.lgr.Info("queue.consumer.close", zap.String("topic", RetryTopic))
			return
		default:
			msg, err := c.ReadMessage(time.Duration(q.config.ConsumerReadTtl))
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok == true && kafkaErr.Code() != kafka.ErrTimedOut {
					q.lgr.Error("queue.consumer.prod.read", zap.Error(err))
					// todo: set alert with prometheus
					break
				}
			}

			if msg != nil {
				attempt := 0

				for {
					attempt++
					fmt.Println("tries ", attempt)

					if attempt == 3 {
						fmt.Println("to dql with tries of ", attempt)

						if err = q.Produce(ctx, DlqTopic, string(msg.Key), msg.Value); err != nil {
							//todo: set grafana/prometheus alarm
							q.lgr.Error("queue.consumer.provider.reproduce",
								zap.String("topic", RetryTopic),
								zap.String("message.id", string(msg.Key)),
								zap.Error(err),
							)
						}
						return
					}

					time.Sleep(time.Duration(attempt*3) * time.Second)

					if err = q.consumeAndSendMessage(ctx, RetryTopic, msg); err == nil {
						var value domain.OutboxMessage
						if err = json.Unmarshal(msg.Value, &value); err != nil {
							q.lgr.Error("queue.consumer.retry.parse", zap.String("topic", RetryTopic), zap.Error(err))
							return
						}

						if err = q.outbox.UpdateTryCount(ctx, value.OutboxId, attempt); err != nil {
							q.lgr.Error("queue.consumer.db.update.retry",
								zap.String("topic", RetryTopic),
								zap.String("message.id", string(msg.Key)),
								zap.Error(err),
							)
						}

						return
					}
				}
			}
		}
	}
}

func (q *queue) dlqTopicConsumer(ctx context.Context, handler chan struct{}) {
	c := q.consumers[DlqTopic]

	if err := c.Subscribe(DlqTopic, nil); err != nil {
		q.lgr.Error("queue.consumer.subscribe.dlq", zap.Error(err))

		utils.PrintStd(utils.StdPanic, "queue.consumer.subscribe.dlq: %s", err.Error())

		// todo: set alert with prometheus
		// todo: handle the db try counter and send to retry
	}

	for {
		select {
		case <-handler:
			if err := c.Close(); err != nil {
				q.lgr.Error("queue.consumer.close", zap.String("topic", DlqTopic), zap.Error(err))
				return
			}

			q.lgr.Info("queue.consumer.close", zap.String("topic", DlqTopic))
			return
		default:
			msg, err := c.ReadMessage(time.Duration(q.config.ConsumerReadTtl))
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok == true && kafkaErr.Code() != kafka.ErrTimedOut {
					q.lgr.Error("queue.consumer.prod.read", zap.Error(err))
					// todo: set alert with prometheus
					break
				}
			}

			if msg != nil {
				var value domain.OutboxMessage
				if err = json.Unmarshal(msg.Value, &value); err != nil {
					q.lgr.Error("queue.consumer.parse", zap.String("topic", DlqTopic), zap.Error(err))
					//todo set alarm
					break
				}

				err = q.updateStatus(ctx, value.MessageId, value.OutboxId, domain.MsgFailed, domain.OutboxFailed)
				if err != nil {
					q.lgr.Error("queue.consumer.db.status", zap.String("topic", DlqTopic), zap.Error(err))
					return
				}
			}
		}
	}
}

//

func (q *queue) consumeAndSendMessage(ctx context.Context, t string, msg *kafka.Message) (err error) {
	var value domain.OutboxMessage
	if err = json.Unmarshal(msg.Value, &value); err != nil {
		q.lgr.Error("queue.consumer.parse", zap.String("topic", t), zap.Error(err))
		return
	}

	err = q.updateStatus(ctx, value.MessageId, value.OutboxId, domain.MsgSending, domain.OutboxPublishing)
	if err != nil {
		if err = q.Produce(ctx, RetryTopic, string(msg.Key), msg.Value); err != nil {
			//todo: set grafana/prometheus alarm
			q.lgr.Error("queue.consumer.provider.reproduce",
				zap.String("topic", t),
				zap.String("message.id", string(msg.Key)),
				zap.Error(err),
			)
		}

		return
	}

	if _, err = q.sms.Send(value.Mobile, value.MessageText); err != nil {
		q.lgr.Error("queue.consumer.provider.send",
			zap.String("topic", t),
			zap.String("message.id", string(msg.Key)),
			zap.Error(err),
		)

		return
	}

	err = q.updateStatus(ctx, value.MessageId, value.OutboxId, domain.MsgSent, domain.OutboxPublished)
	return
}

func (q *queue) updateStatus(ctx context.Context, msgId, outboxId uint, msgSt domain.MessageStatus, outboxSt domain.OutboxStatus) (err error) {
	db := q.sql.Tx()
	tx := db.Begin()

	if err = q.message.UpdateStatus(ctx, msgId, string(msgSt)); err != nil {
		_ = tx.Rollback()
		//todo: prometheus/grafana alarm for sent but not updated status
		q.lgr.Error("queue.consumer.prod.db.message",
			zap.String("staus", string(msgSt)),
			zap.Uint("message.id", msgId), zap.Error(err))
		return
	}

	if err = q.outbox.UpdateStatus(ctx, outboxId, string(outboxSt)); err != nil {
		_ = tx.Rollback()
		//todo: prometheus/grafana alarm for sent but not updated status
		q.lgr.Error("queue.consumer.prod.db.outbox",
			zap.String("staus", string(outboxSt)),
			zap.Uint("outbox.id", outboxId), zap.Error(err))
		return
	}

	_ = tx.Commit()
	return
}
