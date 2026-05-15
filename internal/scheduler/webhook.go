package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookPayload is the JSON body sent to a webhook endpoint.
type WebhookPayload struct {
	Name      string    `json:"name"`
	FiredAt   time.Time `json:"fired_at"`
	Expression string   `json:"expression"`
	Status    string    `json:"status"`
	Error     string    `json:"error,omitempty"`
}

// WebhookNotifier sends an HTTP POST to a URL when a job fires.
type WebhookNotifier struct {
	URL     string
	Client  *http.Client
	Timeout time.Duration
}

// NewWebhookNotifier creates a WebhookNotifier with sensible defaults.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL:     url,
		Client:  &http.Client{},
		Timeout: 5 * time.Second,
	}
}

// Notify sends a webhook POST with the given payload.
func (w *WebhookNotifier) Notify(ctx context.Context, payload WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, w.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.Client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// WithWebhook returns a JobFunc middleware that fires a webhook after the job runs.
func WithWebhook(notifier *WebhookNotifier, name, expression string) func(JobFunc) JobFunc {
	return func(next JobFunc) JobFunc {
		return func(ctx context.Context) error {
			start := time.Now()
			jobErr := next(ctx)

			payload := WebhookPayload{
				Name:       name,
				FiredAt:    start,
				Expression: expression,
				Status:     "ok",
			}
			if jobErr != nil {
				payload.Status = "error"
				payload.Error = jobErr.Error()
			}

			_ = notifier.Notify(ctx, payload) // best-effort; don't override job error
			return jobErr
		}
	}
}
