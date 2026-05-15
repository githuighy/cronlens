package predictor

import "errors"

// ErrNoNextRun is returned when no matching run time can be found within the
// search window (e.g. an expression that can never be satisfied).
var ErrNoNextRun = errors.New("predictor: no next run time found within search window")
