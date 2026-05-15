// Package scheduler provides scheduling primitives built on top of cronlens
// cron expression parsing.
//
// # Snapshot
//
// A Snapshot captures the observable state of a named schedule at a specific
// moment in time. It is intended for diagnostics, audit trails, and dashboard
// integrations where a lightweight, serialisable record of schedule state is
// needed without holding a reference to the live Schedule object.
//
// Fields captured per snapshot:
//
//   - Name       — identifier matching the Schedule name
//   - Expression — the raw cron expression
//   - Timezone   — IANA timezone string (e.g. "America/New_York")
//   - Tags       — copy of key/value metadata at capture time
//   - NextRun    — the next scheduled execution time
//   - LastRun    — the most recent execution time (nil if never run)
//   - CapturedAt — UTC timestamp set automatically by Save
//
// # SnapshotStore
//
// SnapshotStore is a thread-safe, in-memory registry of Snapshot values keyed
// by schedule name. It supports save, get, delete, and enumeration operations.
// Saving a snapshot with an existing name overwrites the previous entry.
package scheduler
