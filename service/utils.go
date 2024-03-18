package service

import "time"

func Now() time.Time {
	return time.Now().Truncate(time.Millisecond)
}

func UTCNow() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}
