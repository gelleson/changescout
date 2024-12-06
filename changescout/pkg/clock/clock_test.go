package clock

import (
	"testing"
	"time"
)

func TestClock_Now(t *testing.T) {
	fixedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		clock     *Clock
		wantFixed bool
	}{
		{"default", New(), false},
		{"fixed", NewFixedTime(fixedTime), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.clock.Now()
			if tt.wantFixed && !got.Equal(fixedTime) {
				t.Errorf("Clock.Now() = %v, want %v", got, fixedTime)
			}
			if !tt.wantFixed && got.Equal(fixedTime) {
				t.Errorf("Clock.Now() = %v, but should not be fixed", got)
			}
		})
	}
}
