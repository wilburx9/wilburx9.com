package common

import (
	"fmt"
	"time"
)

// StringToTime parses timeStr with layout. Returns the current time if parsing fails
func StringToTime(layout string, timeStr string) time.Time {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not parse time: \"%s\" with layout:\"%s\" :: \"%v\"", timeStr, layout, err))
		return time.Now()
	}
	return t
}
