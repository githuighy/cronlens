// Package scheduler provides scheduling primitives for cron-based job execution.
//
// # Pipeline
//
// Pipeline executes a sequence of named stages in order. Unlike JobChain, which
// wraps independent jobs, Pipeline is designed for multi-step workflows where
// each stage is part of a single logical operation.
//
// Basic usage:
//
//	p := scheduler.NewPipeline().
//		Add("fetch",   fetchData).
//		Add("process", processData).
//		Add("store",   storeResults)
//
//	if err := p.Run(ctx); err != nil {
//		log.Println("pipeline failed:", err)
//	}
//
// Logging middleware:
//
//	logged := scheduler.WithPipelineLogging(p)
//	logged.Run(ctx)
//
// Converting to a single job:
//
//	job := scheduler.AsPipelineJob(p, "my-pipeline")
package scheduler
