package duration

import (
	"math"
	"time"
)

const millisecondsRoundFactor = 1000

func GetDurationInMilliseconds(start time.Time) float64 {
	duration := time.Since(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := math.Round(milliseconds*millisecondsRoundFactor) / millisecondsRoundFactor

	return rounded
}
