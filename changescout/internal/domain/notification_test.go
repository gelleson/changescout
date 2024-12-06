package domain

import (
	"testing"
)

func TestNotificationType_Values(t *testing.T) {
	tests := []struct {
		name     string
		expected []string
	}{
		{
			name:     "BasicCase",
			expected: []string{"telegram"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notificationType := NotificationType("TelegramNotificationType")
			if values := notificationType.Values(); !equalStringSlices(values, tt.expected) {
				t.Errorf("Values() = %v, expected %v", values, tt.expected)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
