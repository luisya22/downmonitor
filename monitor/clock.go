package monitor

import "time"

type Clock interface {
	Now() time.Time
}

type RealTime struct{}

func (rt *RealTime) Now() time.Time {
	return time.Now()
}

type MockTime struct {
	MockTime time.Time
}

func (mt *MockTime) Now() time.Time {
	return mt.MockTime
}
