package utils

import "time"

type RealClock struct{}

func (c *RealClock) NowUnix() int64 {
	return time.Now().Unix()
}
