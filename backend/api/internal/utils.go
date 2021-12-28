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

// MakeId generates an id
func MakeId(source string, currentId string) string {
	data := []byte(fmt.Sprintf("%v/%v", source, currentId))
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}


// GetCacheControl returns Cache-Control instructions
func GetCacheControl(updatedAt time.Time) string {
	maxAge := oneDay.Seconds()
	now := time.Now()
	diff := now.Sub(updatedAt)

	if twoWeeks > diff { // Confirm that the data is less than two weeks old
		expiresIn := updatedAt.Add(twoWeeks).Sub(now)
		maxAge =  expiresIn.Seconds()
	}

	sMaxAge := maxAge / 2
	return fmt.Sprintf("public, max-age=%v, s-maxage=%v", int(maxAge), int(sMaxAge))
}
