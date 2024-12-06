package clock

import (
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
	"time"
)

type Clock struct {
	fixedTime *time.Time
}

func (c Clock) Now() time.Time {
	if c.fixedTime != nil {
		return *c.fixedTime
	}
	return time.Now()
}

func New() *Clock {
	return &Clock{}
}

func NewFixedTime(t time.Time) *Clock {
	return &Clock{
		fixedTime: transform.ToPtr(t),
	}
}
