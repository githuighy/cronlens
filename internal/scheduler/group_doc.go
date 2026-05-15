// Package scheduler provides scheduling primitives built on top of cronlens
// cron expression parsing.
//
// # Group
//
// Group lets you manage multiple named [Schedule] instances as a single unit.
// It is safe for concurrent use.
//
// Basic usage:
//
//	g := scheduler.NewGroup()
//
//	s1, _ := scheduler.New("0 9 * * 1-5", "America/New_York")
//	g.Add("standup", s1)
//
//	s2, _ := scheduler.New("30 17 * * 5", "America/New_York")
//	g.Add("weekly-report", s2)
//
//	// Find which job fires next.
//	entries := g.NextAll(time.Now())
//	for _, e := range entries {
//		fmt.Printf("%s → %s\n", e.Name, e.Next.Format(time.RFC3339))
//	}
//
// Names are unique within a Group; adding a duplicate returns an error.
// Use [Group.Remove] to deregister a schedule by name.
package scheduler
