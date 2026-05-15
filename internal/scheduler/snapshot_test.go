package scheduler

import (
	"testing"
	"time"
)

func baseSnapshot(name string) Snapshot {
	now := time.Now().UTC()
	return Snapshot{
		Name:       name,
		Expression: "*/5 * * * *",
		Timezone:   "UTC",
		Tags:       map[string]string{"env": "test"},
		NextRun:    now.Add(5 * time.Minute),
		LastRun:    &now,
	}
}

func TestSnapshotStore_SaveAndGet(t *testing.T) {
	store := NewSnapshotStore()
	snap := baseSnapshot("job-a")
	store.Save(snap)

	got, ok := store.Get("job-a")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if got.Name != "job-a" {
		t.Errorf("expected name job-a, got %s", got.Name)
	}
	if got.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestSnapshotStore_GetMissing(t *testing.T) {
	store := NewSnapshotStore()
	_, ok := store.Get("nonexistent")
	if ok {
		t.Error("expected false for missing snapshot")
	}
}

func TestSnapshotStore_Overwrite(t *testing.T) {
	store := NewSnapshotStore()
	store.Save(baseSnapshot("job-b"))

	updated := baseSnapshot("job-b")
	updated.Expression = "0 * * * *"
	store.Save(updated)

	got, _ := store.Get("job-b")
	if got.Expression != "0 * * * *" {
		t.Errorf("expected updated expression, got %s", got.Expression)
	}
	if store.Count() != 1 {
		t.Errorf("expected count 1, got %d", store.Count())
	}
}

func TestSnapshotStore_Delete(t *testing.T) {
	store := NewSnapshotStore()
	store.Save(baseSnapshot("job-c"))
	store.Delete("job-c")

	_, ok := store.Get("job-c")
	if ok {
		t.Error("expected snapshot to be deleted")
	}
}

func TestSnapshotStore_All(t *testing.T) {
	store := NewSnapshotStore()
	store.Save(baseSnapshot("job-1"))
	store.Save(baseSnapshot("job-2"))
	store.Save(baseSnapshot("job-3"))

	all := store.All()
	if len(all) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(all))
	}
}

func TestSnapshotStore_Count(t *testing.T) {
	store := NewSnapshotStore()
	if store.Count() != 0 {
		t.Error("expected empty store")
	}
	store.Save(baseSnapshot("x"))
	if store.Count() != 1 {
		t.Errorf("expected count 1, got %d", store.Count())
	}
}
