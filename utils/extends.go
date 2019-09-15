package utils

import "time"

type RawTime string

func (t RawTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(t))
}