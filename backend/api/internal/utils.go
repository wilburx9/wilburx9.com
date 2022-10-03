package internal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	oneDay   = 24 * time.Hour
	twoWeeks = oneDay * 14
)

// StringToTime parses timeStr with layout. Returns the current time if parsing fails
func StringToTime(layout, timeStr, ts string) time.Time {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		log.WithFields(log.Fields{
			"caller": ts,
			"source": timeStr,
			"layout": layout,
			"error":  err,
		}).Warning("Could not parse time")
		return time.Now()
	}
	return t
}

// MakeId generates an id
func MakeId(source string, currentId string) string {
	data := []byte(fmt.Sprintf("%v/%v", source, currentId))
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

// GetCacheControl returns Cache-Control instructions
func GetCacheControl(dataUpdatedAt time.Time) string {
	maxAge := oneDay.Seconds()
	now := time.Now()
	diff := now.Sub(dataUpdatedAt)

	// Check if the data is upto two weeks old yet
	if twoWeeks > diff {
		expiresIn := dataUpdatedAt.Add(twoWeeks).Sub(now) // Get how long until the data is two weeks old.
		maxAge = expiresIn.Seconds()
	}

	sMaxAge := maxAge / 2
	return fmt.Sprintf("public, max-age=%v, s-maxage=%v", int(maxAge), int(sMaxAge))
}
