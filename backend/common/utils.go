package common

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// StringToTime parses timeStr with layout. Returns the current time if parsing fails
func StringToTime(layout string, timeStr string) time.Time {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		log.WithFields(log.Fields{
			"source": timeStr,
			"layout": layout,
			"error":  err,
		}).Warning("Could not parse time")
		return time.Now()
	}
	return t
}

// GetFirstNCodePoints Returns the first n code points of string. E.g FirstNCodePoints("世界 Hello", 1) == "世"
func GetFirstNCodePoints(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}
