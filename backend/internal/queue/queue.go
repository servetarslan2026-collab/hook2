package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type WebhookJob struct {
	ApplicationID  uuid.UUID `json:"application_id"`
	EventID        uuid.UUID `json:"event_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	TargetURL      string    `json:"target_url"`
	Secret         string    `json:"secret"`
	EventType      string    `json:"event_type"`
	Payload        []byte    `json:"payload"`
	AttemptNumber  int       `json:"attempt_number"`
	MaxAttempts    int       `json:"max_attempts"`
	RetryDelays    string    `json:"retry_delays"`
	CreatedAt      time.Time `json:"created_at"`
}

type Queue struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	logger *zap.Logger
}

func Connect(natsURL string, logger *zap.Logger) (*Queue, error) {
	nc, err := nats.Connect(natsURL,
		nats.ReconnectWait(time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Warn("NATS disconnected", zap.Error(err))
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS reconnected")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("connect NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("JetStream context: %w", err)
	}

	q := &Queue{conn: nc, js: js, logger: logger}

	if err := q.setupStreams(); err != nil {
		nc.Close()
		return nil, fmt.Errorf("setup streams: %w", err)
	}

	return q, nil
}

func (q *Queue) setupStreams() error {
	streams := []struct {
		name     string
		subjects []string
	}{
		{"WEBHOOK_DELIVERY", []string{"webhook.delivery"}},
		{"WEBHOOK_RETRY", []string{"webhook.retry"}},
		{"WEBHOOK_DLQ", []string{"webhook.dlq"}},
	}

	for _, s := range streams {
		_, err := q.js.StreamInfo(s.name)
		if err != nil {
			_, err = q.js.AddStream(&nats.StreamConfig{
				Name:       s.name,
				Subjects:   s.subjects,
				Storage:    nats.FileStorage,
				Retention:  nats.WorkQueuePolicy,
				MaxAge:     7 * 24 * time.Hour,
				MaxMsgs:    1_000_000,
				Discard:    nats.DiscardOld,
				Replicas:   1,
			})
			if err != nil {
				return fmt.Errorf("create stream %s: %w", s.name, err)
			}
			q.logger.Info("Created NATS stream", zap.String("stream", s.name))
		}
	}
	return nil
}

func (q *Queue) PublishEvent(ctx context.Context, job *WebhookJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}

	_, err = q.js.Publish("webhook.delivery", data)
	if err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	q.logger.Debug("Published webhook job",
		zap.String("event_id", job.EventID.String()),
		zap.String("subscription_id", job.SubscriptionID.String()),
	)
	return nil
}

func (q *Queue) PublishRetry(ctx context.Context, job *WebhookJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal retry job: %w", err)
	}

	_, err = q.js.Publish("webhook.retry", data)
	if err != nil {
		return fmt.Errorf("publish retry: %w", err)
	}

	q.logger.Debug("Published retry job",
		zap.String("event_id", job.EventID.String()),
		zap.Int("attempt", job.AttemptNumber),
	)
	return nil
}

func (q *Queue) PublishDLQ(ctx context.Context, job *WebhookJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal DLQ job: %w", err)
	}

	_, err = q.js.Publish("webhook.dlq", data)
	if err != nil {
		return fmt.Errorf("publish DLQ: %w", err)
	}

	q.logger.Warn("Sent to dead letter queue",
		zap.String("event_id", job.EventID.String()),
		zap.String("subscription_id", job.SubscriptionID.String()),
	)
	return nil
}

func (q *Queue) SubscribeDelivery(handler func(*WebhookJob) error) (*nats.Subscription, error) {
	return q.js.QueueSubscribe("webhook.delivery", "webhook-workers", func(msg *nats.Msg) {
		var job WebhookJob
		if err := json.Unmarshal(msg.Data, &job); err != nil {
			q.logger.Error("Failed to unmarshal delivery job", zap.Error(err))
			msg.Nak()
			return
		}

		if err := handler(&job); err != nil {
			q.logger.Error("Failed to process delivery job", zap.Error(err),
				zap.String("event_id", job.EventID.String()))
			msg.Nak()
			return
		}

		msg.Ack()
	}, nats.ManualAck(), nats.MaxAckPending(100))
}

func (q *Queue) SubscribeRetry(handler func(*WebhookJob) error) (*nats.Subscription, error) {
	return q.js.QueueSubscribe("webhook.retry", "webhook-retry-workers", func(msg *nats.Msg) {
		var job WebhookJob
		if err := json.Unmarshal(msg.Data, &job); err != nil {
			q.logger.Error("Failed to unmarshal retry job", zap.Error(err))
			msg.Nak()
			return
		}

		if err := handler(&job); err != nil {
			q.logger.Error("Failed to process retry job", zap.Error(err))
			msg.Nak()
			return
		}

		msg.Ack()
	}, nats.ManualAck(), nats.MaxAckPending(50))
}

func (q *Queue) Close() {
	if q.conn != nil {
		q.conn.Close()
	}
}
