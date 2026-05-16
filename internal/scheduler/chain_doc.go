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
package scheduler
