// Package humanizer provides utilities for converting parsed cron expressions
// into human-readable English descriptions.
//
// It works in conjunction with the parser package, accepting *parser.Expression
// values and producing clear, natural language summaries of when a cron job
// will run.
//
// Example usage:
//
//	expr, err := parser.Parse("0 9 * * 1-5")
//	if err != nil {
//		log.Fatal(err)
//	}
//	desc := humanizer.Humanize(expr)
//	fmt.Println(desc) // => "Minute 0, hour 9, day of week 1/2/3/4/5"
//
// The package also exposes MonthName and WeekdayName helpers for mapping
// numeric cron values to their string equivalents.
package humanizer
