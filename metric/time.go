// lightmetric helps with easy manipulation of simple telegraf metrics
//  without telegraf libraries dependencies
//
// License: The MIT License (MIT)

package metric

import (
	"time"
)

// TimeWithPrecision returns the given time rounded to the given precision.
func TimeWithPrecision(t time.Time, precision time.Duration) time.Time {
	return t.Round(precision)
}
