# cronlens

> Human-readable cron expression parser and next-run predictor with timezone awareness

---

## Installation

```bash
go install github.com/yourusername/cronlens@latest
```

Or add it as a library:

```bash
go get github.com/yourusername/cronlens
```

---

## Usage

### CLI

```bash
# Parse a cron expression and show next 5 run times
cronlens "0 9 * * MON-FRI" --tz="America/New_York" --next=5
```

**Output:**
```
Expression : 0 9 * * MON-FRI
Description: At 09:00, Monday through Friday
Timezone   : America/New_York

Next runs:
  1. Mon, 14 Jul 2025 09:00:00 EDT
  2. Tue, 15 Jul 2025 09:00:00 EDT
  3. Wed, 16 Jul 2025 09:00:00 EDT
  4. Thu, 17 Jul 2025 09:00:00 EDT
  5. Fri, 18 Jul 2025 09:00:00 EDT
```

### Library

```go
import "github.com/yourusername/cronlens"

expr, err := cronlens.Parse("0 9 * * MON-FRI")
if err != nil {
    log.Fatal(err)
}

fmt.Println(expr.Describe())
// Output: "At 09:00, Monday through Friday"

loc, _ := time.LoadLocation("America/New_York")
next := expr.NextN(time.Now(), loc, 3)
for _, t := range next {
    fmt.Println(t)
}
```

---

## Features

- Translates cron expressions into plain English descriptions
- Predicts upcoming run times with full timezone support
- Supports standard 5-field and extended 6-field (with seconds) cron syntax
- Zero external dependencies

---

## License

[MIT](LICENSE)