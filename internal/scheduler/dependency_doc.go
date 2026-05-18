// Package scheduler provides scheduling primitives for cronlens.
//
// # Dependency
//
// DependencyStore allows jobs to declare dependencies on other jobs.
// A job will only execute when all of its declared dependencies have
// been marked as done via MarkDone.
//
// Basic usage:
//
//	store := scheduler.NewDependencyStore()
//
//	// jobB depends on jobA
//	store.Declare("jobB", "jobA")
//
//	// wrap jobB so it checks deps before running
//	protected := scheduler.WithDependency(store, "jobB", jobBFunc)
//
//	// after jobA succeeds, mark it done
//	store.MarkDone("jobA")
//
//	// now jobB will run
//	protected(ctx)
//
// Circular dependencies are detected eagerly at Declare time and
// return ErrCircularDependency. If a dependency is not yet satisfied
// at run time, WithDependency returns ErrDependencyNotMet without
// invoking the wrapped function.
package scheduler
