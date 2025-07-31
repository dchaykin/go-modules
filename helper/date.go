package helper

import (
	"errors"
	"strings"
	"time"

	"github.com/dchaykin/go-modules/log"
)

var knownDateLayouts = []string{
	time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
	time.RFC1123,          // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,         // "Mon, 02 Jan 2006 15:04:05 -0700"
	time.RFC822,           // "02 Jan 06 15:04 MST"
	time.RFC822Z,          // "02 Jan 06 15:04 -0700"
	time.RFC850,           // "Monday, 02-Jan-06 15:04:05 MST"
	"2006-01-02",          // ISO-Date only
	"02.01.2006",          // German format
	"2006-01-02 15:04:05", // MySQL-style datetime
	"02.01.2006 15:04:05", // German datetime
	"2006/01/02",          // Slashed format
	"02-Jan-2006",         // Common alt format
	"20060102",            // Compact YYYYMMDD
}

func ParseFlexibleDate(input string) (time.Time, error) {
	input = strings.TrimSpace(input)
	for _, layout := range knownDateLayouts {
		if t, err := time.Parse(layout, input); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("unrecognized date format: " + input)
}

func ParseDateTimeToMysql(input string) string {
	if input == "2006-01-02 15:04:05" {
		return input
	}
	if input == "" {
		return ""
	}
	t, err := ParseFlexibleDate(input)
	if err != nil {
		log.WrapError(err)
		return input
	}
	return t.Format("2006-01-02 15:04:05")
}
