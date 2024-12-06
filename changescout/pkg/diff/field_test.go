package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DiffTestSuite struct {
	suite.Suite
}

func (suite *DiffTestSuite) TestGetUpdatedValue() {
	tests := []struct {
		name          string
		newValue      *int
		oldValue      int
		expectedValue int
	}{
		{
			name:          "New value is nil",
			newValue:      nil,
			oldValue:      5,
			expectedValue: 5,
		},
		{
			name:          "Different new value",
			newValue:      intPointer(10),
			oldValue:      5,
			expectedValue: 10,
		},
		{
			name:          "Same new value",
			newValue:      intPointer(5),
			oldValue:      5,
			expectedValue: 5,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := GetUpdatedValue(tt.newValue, tt.oldValue)
			assert.Equal(suite.T(), tt.expectedValue, result, "they should be equal")
		})
	}
}

func (suite *DiffTestSuite) TestGetUpdatedValueWithPointer() {
	tests := []struct {
		name          string
		newValue      *int
		oldValue      *int
		expectedValue *int
	}{
		{
			name:          "Both non-nil and different",
			newValue:      intPointer(10),
			oldValue:      intPointer(5),
			expectedValue: intPointer(10),
		},
		{
			name:          "Both non-nil and same",
			newValue:      intPointer(5),
			oldValue:      intPointer(5),
			expectedValue: intPointer(5),
		},
		{
			name:          "newValue nil",
			newValue:      nil,
			oldValue:      intPointer(5),
			expectedValue: intPointer(5),
		},
		{
			name:          "oldValue nil",
			newValue:      intPointer(10),
			oldValue:      nil,
			expectedValue: intPointer(10),
		},
		{
			name:          "Both nil",
			newValue:      nil,
			oldValue:      nil,
			expectedValue: nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := GetUpdatedValueWithPointer(tt.newValue, tt.oldValue)
			assert.Equal(suite.T(), tt.expectedValue, result, "they should be equal")
		})
	}
}

func (suite *DiffTestSuite) TestCompareSlices() {
	tests := []struct {
		name          string
		oldValue      []int
		newValue      []int
		expectedValue []int
	}{
		{
			name:          "Different slices",
			oldValue:      []int{1, 2, 3},
			newValue:      []int{4, 5, 6},
			expectedValue: []int{4, 5, 6},
		},
		{
			name:          "Same slices",
			oldValue:      []int{1, 2, 3},
			newValue:      []int{1, 2, 3},
			expectedValue: []int{1, 2, 3},
		},
		{
			name:          "New slice is empty",
			oldValue:      []int{1, 2, 3},
			newValue:      []int{},
			expectedValue: []int{1, 2, 3},
		},
		{
			name:          "Old slice is empty",
			oldValue:      []int{},
			newValue:      []int{4, 5, 6},
			expectedValue: []int{4, 5, 6},
		},
		{
			name:          "Both slices empty",
			oldValue:      []int{},
			newValue:      []int{},
			expectedValue: []int{},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := CompareSlices(tt.oldValue, tt.newValue)
			assert.Equal(suite.T(), tt.expectedValue, result, "they should be equal")
		})
	}
}

func TestDiffTestSuite(t *testing.T) {
	suite.Run(t, new(DiffTestSuite))
}

// Вспомогательная функция
func intPointer(i int) *int {
	return &i
}
