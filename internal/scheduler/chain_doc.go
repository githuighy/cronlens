// Package scheduler provides scheduling primitives for cron-based job execution.
//
// # JobChain
//
// JobChain allows you to compose multiple jobs into an ordered pipeline.
// Each step is identified by a name and executed sequentially. If any step
// returns an error, the chain halts immediately and returns a wrapped error
// identifying the failing step.
//
// Example usage:
//
//	chain := scheduler.NewJobChain(5 * time.Second).
//		Add("fetch", fetchData).
//		Add("process", processData).
//		Add("store", storeData)
//
//	if err := chain.Run(ctx); err != nil {
//		log.Printf("pipeline failed: %v", err)
//	}
//
// A per-step timeout can be set via NewJobChain; pass 0 to disable timeouts.
//
// # Error Handling
//
// Errors returned by chain steps are wrapped with the step name using fmt.Errorf,
// so callers can use errors.Is or errors.As to inspect the underlying cause:
//
//	if err := chain.Run(ctx); err != nil {
//		var stepErr *scheduler.StepError
//		if errors.As(err, &stepErr) {
//			log.Printf("step %q failed: %v", stepErr.Step, stepErr.Unwrap())
//		}
//	}
package scheduler
