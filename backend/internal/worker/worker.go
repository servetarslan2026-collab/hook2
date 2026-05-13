package worker

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"webhook-service/internal/api/middleware"
	"webhook-service/internal/queue"
	"webhook-service/internal/store"
)

const (
	maxAttempts = 3
	httpTimeout = 30 * time.Second
	maxBodySize = 1 << 20 // 1MB
)

var retryDelays = []time.Duration{
	1 * time.Second,
	30 * time.Second,
	5 * time.Minute,
}

type Worker struct {
	queue  *queue.Queue
	store  *store.Store
	client *http.Client
	logger *zap.Logger
}

func New(q *queue.Queue, s *store.Store, logger *zap.Logger) *Worker {
	return &Worker{
		queue: q,
		store: s,
		client: &http.Client{
			Timeout: httpTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		logger: logger,
	}
}

func (w *Worker) Start(ctx context.Context, count int) error {
	_, err := w.queue.SubscribeDelivery(func(job *queue.WebhookJob) error {
		return w.processDelivery(ctx, job)
	})
	if err != nil {
		return fmt.Errorf("subscribe delivery: %w", err)
	}

	_, err = w.queue.SubscribeRetry(func(job *queue.WebhookJob) error {
		return w.processDelivery(ctx, job)
	})
	if err != nil {
		return fmt.Errorf("subscribe retry: %w", err)
	}

	w.logger.Info("Worker started", zap.Int("count", count))
	return nil
}

func (w *Worker) processDelivery(ctx context.Context, job *queue.WebhookJob) error {
	start := time.Now()

	// Build payload
	payload, _ := json.Marshal(map[string]interface{}{
		"event_id":   job.EventID.String(),
		"event_type": job.EventType,
		"payload":    json.RawMessage(job.Payload),
		"created_at": job.CreatedAt.Format(time.RFC3339),
	})

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", job.TargetURL, bytes.NewReader(payload))
	if err != nil {
		return w.handleFailure(ctx, job, 0, "", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "WebhookService/1.0")
	req.Header.Set("X-Webhook-Event", job.EventType)
	req.Header.Set("X-Webhook-ID", job.EventID.String())

	// Sign payload if secret exists
	if job.Secret != "" {
		signature := signPayload(payload, job.Secret)
		req.Header.Set("X-Webhook-Signature", "sha256="+signature)
	}

	// Execute request
	resp, err := w.client.Do(req)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		return w.handleFailure(ctx, job, 0, "", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))

	// Success: 2xx
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		w.store.CreateDeliveryAttempt(ctx, job.EventID, job.SubscriptionID, "success", resp.StatusCode, string(payload), string(respBody), duration, job.AttemptNumber)
		middleware.IncWebhookDelivery("success")
		w.logger.Info("Webhook delivered",
			zap.String("event_id", job.EventID.String()),
			zap.Int("status", resp.StatusCode),
			zap.Int64("duration_ms", duration),
		)
		return nil
	}

	// Non-2xx: failure
	w.store.CreateDeliveryAttempt(ctx, job.EventID, job.SubscriptionID, "failed", resp.StatusCode, string(payload), string(respBody), duration, job.AttemptNumber)
	middleware.IncWebhookDelivery("failed")
	return w.handleFailure(ctx, job, resp.StatusCode, string(respBody), fmt.Errorf("status %d", resp.StatusCode))
}

func (w *Worker) handleFailure(ctx context.Context, job *queue.WebhookJob, statusCode int, respBody string, err error) error {
	w.logger.Warn("Webhook delivery failed",
		zap.String("event_id", job.EventID.String()),
		zap.Int("attempt", job.AttemptNumber),
		zap.Error(err),
	)

	// Retry if under max attempts
	if job.AttemptNumber < maxAttempts {
		job.AttemptNumber++
		delay := retryDelays[job.AttemptNumber-1]

		w.logger.Info("Scheduling retry",
			zap.String("event_id", job.EventID.String()),
			zap.Int("attempt", job.AttemptNumber),
			zap.Duration("delay", delay),
		)

		time.AfterFunc(delay, func() {
			if err := w.queue.PublishRetry(ctx, job); err != nil {
				w.logger.Error("Failed to publish retry", zap.Error(err))
			}
		})
		return nil
	}

	// Max attempts: dead letter
	w.logger.Error("Max attempts reached, sending to DLQ",
		zap.String("event_id", job.EventID.String()),
		zap.String("subscription_id", job.SubscriptionID.String()),
	)

	w.store.CreateDeliveryAttempt(ctx, job.EventID, job.SubscriptionID, "dead_letter", statusCode, "", respBody, 0, job.AttemptNumber)
	middleware.IncWebhookDelivery("dead_letter")
	return w.queue.PublishDLQ(ctx, job)
}

func signPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
