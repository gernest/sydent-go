package models

import "time"

// MS returns  time as milliseconds since epoch.
func MS(ts *time.Time) int64 {
	return ts.UnixNano() / int64(time.Millisecond)
}

// FromMS returns time.Time initialized from ms which is time in milliseconds
// since unix epoch.
func FromMS(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond))
}

// Time returns the current time in milliseconds.
func Time() int64 {
	now := time.Now()
	return MS(&now)
}
