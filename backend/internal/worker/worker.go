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
	"webhook-service/internal/models"
	"webhook-service/internal/queue"
	"webhook-service/internal/store"
)

const (
	maxAttempts     = 3
	httpTimeout     = 30 * time.Second
	maxBodySize     = 1 << 20 // 1MB
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
	// Subscribe to delivery stream
	_, err := w.queue.SubscribeDelivery(func(job *queue.WebhookJob) error {
		return w.processDelivery(ctx, job)
	})
	if err != nil {
		return fmt.Errorf("subscribe delivery: %w", err)
	}

	// Subscribe to retry stream
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

	// Build request body
	body, _ := json.Marshal(map[string]interface{}{
		"event_id":   job.EventID.String(),
		"event_type": job.EventType,
		"payload":    json.RawMessage(job.Payload),
		"created_at": job.CreatedAt.Format(time.RFC3339),
	})

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", job.TargetURL, bytes.NewReader(body))
	if err != nil {
		return w.handleFailure(ctx, job, 0, "", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "WebhookService/1.0")
	req.Header.Set("X-Webhook-Event", job.EventType)
	req.Header.Set("X-Webhook-ID", job.EventID.String())

	// Sign payload if secret exists
	if job.Secret != "" {
		signature := signPayload(body, job.Secret)
		req.Header.Set("X-Webhook-Signature", "sha256="+signature)
	}

	// Execute request
	resp, err := w.client.Do(req)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		return w.handleFailure(ctx, job, 0, "", err)
	}
	defer resp.Body.Close()

	// Read response body (limited)
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))

	// Record delivery attempt
	attempt := &models.DeliveryAttempt{
		ID:             uuid.New(),
		EventID:        job.EventID,
		SubscriptionID: job.SubscriptionID,
		StatusCode:     resp.StatusCode,
		RequestBody:    string(body),
		ResponseBody:   string(respBody),
		DurationMs:     duration,
		AttemptNumber:  job.AttemptNumber,
		CreatedAt:      time.Now(),
	}

	// Success: 2xx status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		attempt.Status = "success"
		if err := w.store.CreateDeliveryAttempt(ctx, attempt); err != nil {
			w.logger.Error("Failed to record delivery", zap.Error(err))
		}
		w.logger.Info("Webhook delivered",
			zap.String("event_id", job.EventID.String()),
			zap.Int("status", resp.StatusCode),
			zap.Int64("duration_ms", duration),
		)
		return nil
	}

	// Non-2xx: treat as failure
	attempt.Status = "failed"
	w.store.CreateDeliveryAttempt(ctx, attempt)
	return w.handleFailure(ctx, job, resp.StatusCode, string(respBody), fmt.Errorf("status %d", resp.StatusCode))
}

func (w *Worker) handleFailure(ctx context.Context, job *queue.WebhookJob, statusCode int, respBody string, err error) error {
	w.logger.Warn("Webhook delivery failed",
		zap.String("event_id", job.EventID.String()),
		zap.Int("attempt", job.AttemptNumber),
		zap.Error(err),
	)

	// Record failed attempt
	attempt := &models.DeliveryAttempt{
		ID:             uuid.New(),
		EventID:        job.EventID,
		SubscriptionID: job.SubscriptionID,
		Status:         "failed",
		StatusCode:     statusCode,
		RequestBody:    "",
		ResponseBody:   respBody,
		DurationMs:     0,
		AttemptNumber:  job.AttemptNumber,
		CreatedAt:      time.Now(),
	}
	w.store.CreateDeliveryAttempt(ctx, attempt)

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

	// Max attempts reached: send to DLQ
	w.logger.Error("Max attempts reached, sending to DLQ",
		zap.String("event_id", job.EventID.String()),
		zap.String("subscription_id", job.SubscriptionID.String()),
	)

	// Mark as dead letter
	attempt.Status = "dead_letter"
	w.store.CreateDeliveryAttempt(ctx, attempt)

	return w.queue.PublishDLQ(ctx, job)
}

func signPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
