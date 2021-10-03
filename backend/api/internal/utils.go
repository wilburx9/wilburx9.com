package internal

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/configs"
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

// GetDataCollection prepends the running environment passed string and returns it
func GetDataCollection(coll string) string {
	return fmt.Sprintf("%s_%s", configs.Config.Env, coll)
}