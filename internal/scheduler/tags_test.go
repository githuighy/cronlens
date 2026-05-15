package scheduler

import (
	"testing"
)

func TestTags_SetAndGet(t *testing.T) {
	tags := make(Tags)

	if err := tags.Set("env", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, ok := tags.Get("env")
	if !ok || v != "production" {
		t.Errorf("expected env=production, got %q (found=%v)", v, ok)
	}
}

func TestTags_SetEmptyKey(t *testing.T) {
	tags := make(Tags)
	if err := tags.Set("", "value"); err == nil {
		t.Error("expected error for empty key, got nil")
	}
	if err := tags.Set("   ", "value"); err == nil {
		t.Error("expected error for whitespace-only key, got nil")
	}
}

func TestTags_Delete(t *testing.T) {
	tags := Tags{"owner": "alice"}
	tags.Delete("owner")
	if _, ok := tags.Get("owner"); ok {
		t.Error("expected owner to be deleted")
	}
	// deleting non-existent key should not panic
	tags.Delete("nonexistent")
}

func TestTags_Keys(t *testing.T) {
	tags := Tags{"z": "1", "a": "2", "m": "3"}
	keys := tags.Keys()
	expected := []string{"a", "m", "z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("keys[%d]: want %q, got %q", i, expected[i], k)
		}
	}
}

func TestTags_String(t *testing.T) {
	empty := make(Tags)
	if s := empty.String(); s != "(no tags)" {
		t.Errorf("empty tags string: want '(no tags)', got %q", s)
	}

	tags := Tags{"env": "prod", "team": "infra"}
	s := tags.String()
	if s != "env=prod, team=infra" {
		t.Errorf("unexpected string: %q", s)
	}
}

func TestTags_Clone(t *testing.T) {
	original := Tags{"key": "value"}
	cloned := original.Clone()
	cloned.Set("key", "modified")

	if v, _ := original.Get("key"); v != "value" {
		t.Errorf("original was mutated: got %q", v)
	}
}
