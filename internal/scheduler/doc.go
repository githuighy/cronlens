// Package scheduler provides a high-level Schedule type that wraps a parsed
// cron expression together with a timezone location.
//
// It offers three primary operations:
//
//   - NextRun: compute the next scheduled time after a given instant.
//   - NextN: compute the next N scheduled times after a given instant.
//   - TimeSinceLast: compute how long ago the most recent scheduled run occurred.
//
// Example usage:
//
//	s, err := scheduler.New("0 8 * * MON-FRI", "America/New_York")
//	if err != nil {
//		log.Fatal(err)
//	}
//	next, err := s.NextRun(time.Now())
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Next run:", next)
package scheduler
