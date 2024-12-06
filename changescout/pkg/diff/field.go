package diff

import (
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
)

// GetUpdatedValue returns the new value if it is not nil and differs from the old value; otherwise, it returns the old value.
func GetUpdatedValue[T comparable](newValue *T, oldValue T) T {
	if newValue != nil && *newValue != oldValue {
		return *newValue
	}
	return oldValue
}

// GetUpdatedValueWithPointer returns the new value if it is not nil and differs from the old value; otherwise, it returns the old value.
// Checks for nil pointers are included to prevent panic.
func GetUpdatedValueWithPointer[T comparable](newValue *T, oldValue *T) *T {
	// Return newValue if oldValue is nil, but newValue is not nil
	if oldValue == nil && newValue != nil {
		return newValue
	}
	// If both pointers are not nil and their values are different, return newValue
	if newValue != nil && oldValue != nil && *newValue != *oldValue {
		return newValue
	}
	// Return oldValue if it is not nil
	if oldValue != nil {
		return oldValue
	}
	// Return nil if both pointers are nil
	return nil
}

// CompareSlices returns the new slice if its hash differs from the old one's, and if the new slice is not empty.
func CompareSlices[T comparable](oldValue []T, newValue []T) []T {
	if len(oldValue) != len(newValue) || transform.HashSlice(oldValue) != transform.HashSlice(newValue) {
		if len(newValue) > 0 {
			return newValue
		}
	}
	return oldValue
}
