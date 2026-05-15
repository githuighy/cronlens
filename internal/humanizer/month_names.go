package humanizer

var monthNames = map[int]string{
	1:  "January",
	2:  "February",
	3:  "March",
	4:  "April",
	5:  "May",
	6:  "June",
	7:  "July",
	8:  "August",
	9:  "September",
	10: "October",
	11: "November",
	12: "December",
}

var weekdayNames = map[int]string{
	0: "Sunday",
	1: "Monday",
	2: "Tuesday",
	3: "Wednesday",
	4: "Thursday",
	5: "Friday",
	6: "Saturday",
}

// MonthName returns the full name of a month given its number (1-12).
func MonthName(m int) string {
	if name, ok := monthNames[m]; ok {
		return name
	}
	return ""
}

// WeekdayName returns the full name of a weekday given its number (0=Sunday).
func WeekdayName(d int) string {
	if name, ok := weekdayNames[d]; ok {
		return name
	}
	return ""
}

// MonthNamesToList converts a slice of month numbers to a human-readable list.
func MonthNamesToList(months []int) string {
	names := make([]string, 0, len(months))
	for _, m := range months {
		if name := MonthName(m); name != "" {
			names = append(names, name)
		}
	}
	return joinWithAnd(names)
}

// WeekdayNamesToList converts a slice of weekday numbers to a human-readable list.
func WeekdayNamesToList(days []int) string {
	names := make([]string, 0, len(days))
	for _, d := range days {
		if name := WeekdayName(d); name != "" {
			names = append(names, name)
		}
	}
	return joinWithAnd(names)
}

func joinWithAnd(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return items[0] + " and " + items[1]
	default:
		return strings.Join(items[:len(items)-1], ", ") + ", and " + items[len(items)-1]
	}
}
