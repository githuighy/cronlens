// Package validator provides fine-grained validation of cron expressions,
// returning structured errors per field so callers can surface precise
// feedback to users.
//
// Usage:
//
//	result := validator.Validate("*/5 * * * *")
//	if !result.Valid {
//		for _, e := range result.Errors {
//			fmt.Println(e)
//		}
//	}
//
// Supported syntax per field:
//
//	*         — wildcard (any value)
//	a         — literal integer
//	a-b       — inclusive range
//	a-b/n     — range with step
//	*/n       — wildcard with step
//	a,b,c     — comma-separated list (each item may be a range or step)
//
// Field order and bounds:
//
//	Position  Field          Min  Max
//	0         minute           0   59
//	1         hour             0   23
//	2         day-of-month     1   31
//	3         month            1   12
//	4         day-of-week      0    6
package validator
