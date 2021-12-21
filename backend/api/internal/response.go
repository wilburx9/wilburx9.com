package internal

import (
	"fmt"
	"math"
	"time"
)

const (
	twoWeeks = 24 * time.Hour * 14
)

// Result encapsulates the fields and methods needed for saving data from Fetcher(s) and writing to http requests
type Result struct {
	UpdatedAt time.Time   `json:"updated_at" firestore:"updated_at" mapstructure:"updated_at"`
}

// AverageCacheControl returns Cache-Control instructions
func AverageCacheControl(updatedAts []time.Time) string {
	maxAge, sMaxAge := averageMaxAges(updatedAts)
	return fmt.Sprintf("public, max-age=%v, s-maxage=%v", maxAge, sMaxAge)
}

func averageMaxAges(updatedAts []time.Time) (int, int) {
	sumMaxAge := 0.0
	for _, t := range updatedAts {
		sumMaxAge += maxAge(t)
	}
	maxAge := int(math.Ceil(sumMaxAge / float64(len(updatedAts))))
	return maxAge, maxAge / 2
}

// maxAge returns how long (in seconds) this response should be cached by client.
func maxAge(updatedAt time.Time) float64 {
	now := time.Now()
	diff := now.Sub(updatedAt)

	if twoWeeks > diff { // Confirm that the data is less than two weeks old
		expiresIn := updatedAt.Add(twoWeeks).Sub(now)
		return expiresIn.Seconds()
	}

	// The data hasn't been updated for two weeks. Return 0 to force the client to fetch a fresh data
	return 0
}