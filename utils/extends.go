package utils

import "time"

type RawTime string

func (t RawTime) Time() (time.Time, error) {
	return time.Parse("10:01:01", string(t))
}