package internal

import (
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"reflect"
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

// ToTimeHookFunc is a custom hook for parsing Time. Source:: https://github.com/mitchellh/mapstructure/issues/159#issuecomment-482201507
func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}
