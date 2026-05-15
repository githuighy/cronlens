package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhookNotifier_Success(t *testing.T) {
	var received WebhookPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewWebhookNotifier(srv.URL)
	payload := WebhookPayload{
		Name:       "test-job",
		Expression: "* * * * *",
		FiredAt:    time.Now(),
		Status:     "ok",
	}
	if err := n.Notify(context.Background(), payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Name != "test-job" {
		t.Errorf("name: got %q, want %q", received.Name, "test-job")
	}
}

func TestWebhookNotifier_NonSuccessStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := NewWebhookNotifier(srv.URL)
	err := n.Notify(context.Background(), WebhookPayload{Status: "ok"})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestWebhookNotifier_InvalidURL(t *testing.T) {
	n := NewWebhookNotifier("http://127.0.0.1:0") // nothing listening
	err := n.Notify(context.Background(), WebhookPayload{})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestWithWebhook_SuccessJob(t *testing.T) {
	var received WebhookPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewWebhookNotifier(srv.URL)
	job := func(ctx context.Context) error { return nil }
	wrapped := Chain(job, WithWebhook(n, "my-job", "0 * * * *"))

	if err := wrapped(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Status != "ok" {
		t.Errorf("status: got %q, want %q", received.Status, "ok")
	}
}

func TestWithWebhook_FailingJob(t *testing.T) {
	var received WebhookPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewWebhookNotifier(srv.URL)
	jobErr := errors.New("something went wrong")
	job := func(ctx context.Context) error { return jobErr }
	wrapped := Chain(job, WithWebhook(n, "bad-job", "0 * * * *"))

	err := wrapped(context.Background())
	if !errors.Is(err, jobErr) {
		t.Fatalf("expected original job error, got %v", err)
	}
	if received.Status != "error" {
		t.Errorf("status: got %q, want %q", received.Status, "error")
	}
	if received.Error != jobErr.Error() {
		t.Errorf("error field: got %q, want %q", received.Error, jobErr.Error())
	}
}
