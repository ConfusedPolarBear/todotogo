package todo
import (
	"time"
	"log"
	"strings"
)

func ParseDates(raw string) string {
	n := time.Now()
	original := raw
	
	// Simple cases
	raw = replaceRelativeDate(raw, "due:today", n)
	raw = replaceRelativeDate(raw, "due:tomorrow", n.AddDate(0, 0, 1))
	raw = replaceRelativeDate(raw, "due:tom", n.AddDate(0, 0, 1))

	/*
		Relative dates are harder - there doesn't seem to be a way to convert the string "Monday" into a Time object
		To do this, we use a horrible kludge:
			Increment the current date by one day
			Check if the weekday starts with the provided string
			If it does, we found the Time object to use
		Pull requests to improve this are welcome
	*/
	prefixes := []string { "sun", "mon", "tue", "wed", "thu", "fri", "sat" }
	for _, day := range prefixes {
		needle := "due:" + day

		if strings.Contains(strings.ToLower(raw), needle) {
			for i := 1; i <= 7; i++ {
				added := n.AddDate(0, 0, i)
				found := added.Weekday().String()
				found = strings.ToLower(found)

				if strings.HasPrefix(found, day) {
					raw = replaceRelativeDate(raw, needle, added)
				}
			}
		}
	}

	if original != raw {
		log.Printf("Rewrote task from \"%s\" to \"%s\"", original, raw)
	}

	return raw
}

func replaceRelativeDate(haystack, needle string, date time.Time) string {
	if strings.Contains(haystack, needle) {
		return strings.ReplaceAll(haystack, needle, "due:" + formatYMD(date))
	}

	return haystack
}

func formatYMD(date time.Time) string {
	return date.Format("2006-01-02")
}