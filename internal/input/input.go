package input

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var dateLayouts = []string{
	"01-02-2006",
	"01/02/2006",
	"1/2/2006",
}

func ParseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	for _, layout := range dateLayouts {
		t, err := time.ParseInLocation(layout, s, time.Local)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date %q: use MM-DD-YYYY, MM/DD/YYYY, or M/D/YYYY", s)
}

func ParseTime(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	if len(s) < 4 {
		return 0, fmt.Errorf("invalid time %q: use format like 3:00AM or 9:30PM", s)
	}

	period := s[len(s)-2:]
	if period != "AM" && period != "PM" {
		return 0, fmt.Errorf("invalid time %q: must end with AM or PM", s)
	}

	timePart := s[:len(s)-2]
	parts := strings.SplitN(timePart, ":", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time %q: use format like 3:00AM or 9:30PM", s)
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 1 || hour > 12 {
		return 0, fmt.Errorf("invalid hour in %q: must be 1-12", s)
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return 0, fmt.Errorf("invalid minute in %q: must be 0-59", s)
	}

	if period == "AM" && hour == 12 {
		hour = 0
	} else if period == "PM" && hour != 12 {
		hour += 12
	}

	return time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute, nil
}

func Levenshtein(a, b string) int {
	la := len(a)
	lb := len(b)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)

	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min3(
				prev[j]+1,
				curr[j-1]+1,
				prev[j-1]+cost,
			)
		}
		prev, curr = curr, prev
	}

	return prev[lb]
}

func FuzzyMatch(query, target string, threshold float64) bool {
	query = strings.ToLower(strings.TrimSpace(query))
	target = strings.ToLower(strings.TrimSpace(target))

	if query == "" {
		return false
	}
	if query == target {
		return true
	}

	dist := Levenshtein(query, target)
	maxLen := len(query)
	if len(target) > maxLen {
		maxLen = len(target)
	}
	if maxLen == 0 {
		return false
	}

	similarity := 1.0 - float64(dist)/float64(maxLen)
	return similarity >= threshold
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
