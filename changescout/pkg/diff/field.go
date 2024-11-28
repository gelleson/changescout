package diff

import "github.com/gelleson/changescout/changescout/internal/utils/transform"

func GetNewValueOrOldValue[T comparable](newValue *T, oldValue T) T {
	if newValue != nil && *newValue != oldValue {
		return *newValue
	}
	return oldValue
}

func SliceDiffChecker[T comparable](oldValue []T, newValue []T) []T {
	oldHash, newHash := transform.HashSlice(oldValue), transform.HashSlice(newValue)
	if oldHash != newHash && len(newValue) > 0 {
		return newValue
	}
	return oldValue
}
