package util

import "time"

// Date is a shorthand for calling time.Date with time properties set to 0 and with UTC timezone.
func Date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
