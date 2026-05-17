package scheduler

import (
	"context"
	"fmt"
	"sync"
)

// PipelineStage represents a single stage in a pipeline.
type PipelineStage struct {
	Name string
	Job  func(ctx context.Context) error
}

// Pipeline executes a sequence of stages, passing a shared context.
// If any stage fails, execution stops and the error is returned.
type Pipeline struct {
	mu     sync.RWMutex
	stages []PipelineStage
}

// NewPipeline creates an empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Add appends a stage to the pipeline.
func (p *Pipeline) Add(name string, job func(ctx context.Context) error) *Pipeline {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stages = append(p.stages, PipelineStage{Name: name, Job: job})
	return p
}

// Len returns the number of stages.
func (p *Pipeline) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.stages)
}

// Names returns the names of all stages in order.
func (p *Pipeline) Names() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	names := make([]string, len(p.stages))
	for i, s := range p.stages {
		names[i] = s.Name
	}
	return names
}

// Run executes all stages sequentially. Stops on first error.
func (p *Pipeline) Run(ctx context.Context) error {
	p.mu.RLock()
	stages := make([]PipelineStage, len(p.stages))
	copy(stages, p.stages)
	p.mu.RUnlock()

	for _, stage := range stages {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("pipeline aborted before stage %q: %w", stage.Name, err)
		}
		if err := stage.Job(ctx); err != nil {
			return fmt.Errorf("pipeline stage %q failed: %w", stage.Name, err)
		}
	}
	return nil
}
