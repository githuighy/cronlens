package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"
)

// WithPipelineLogging wraps each stage in the pipeline with duration logging.
// It returns a new Pipeline whose stages emit log lines on start and finish.
func WithPipelineLogging(p *Pipeline) *Pipeline {
	logged := NewPipeline()
	for _, stage := range p.stages {
		s := stage // capture
		logged.Add(s.Name, func(ctx context.Context) error {
			log.Printf("[pipeline] starting stage %q", s.Name)
			start := time.Now()
			err := s.Job(ctx)
			dur := time.Since(start)
			if err != nil {
				log.Printf("[pipeline] stage %q failed after %s: %v", s.Name, dur, err)
			} else {
				log.Printf("[pipeline] stage %q completed in %s", s.Name, dur)
			}
			return err
		})
	}
	return logged
}

// AsPipelineJob converts a Pipeline into a single job function suitable for
// use with Chain, retry policies, or any other scheduler primitive.
func AsPipelineJob(p *Pipeline, name string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if err := p.Run(ctx); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		return nil
	}
}
