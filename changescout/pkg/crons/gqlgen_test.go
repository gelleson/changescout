package crons

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalGQL(t *testing.T) {
	cronExpr := CronExpression("0 0 * * *")
	var buf bytes.Buffer

	cronExpr.MarshalGQL(&buf)

	assert.Equal(t, `"0 0 * * *"`, buf.String(), "the marshalled string should be correctly quoted")
}

func TestUnmarshalGQL(t *testing.T) {
	var expr CronExpression

	t.Run("valid cron expression", func(t *testing.T) {
		err := expr.UnmarshalGQL("0 0 * * *")
		assert.NoError(t, err, "valid cron expression should not return an error")
		assert.Equal(t, CronExpression("0 0 * * *"), expr, "the expression should be correctly set")
	})

	t.Run("invalid cron expression", func(t *testing.T) {
		err := expr.UnmarshalGQL("invalid")
		assert.Error(t, err, "invalid cron expression should return an error")
		assert.Contains(t, err.Error(), "invalid cron expression", "the error message should be correct")
	})

	t.Run("non-string value", func(t *testing.T) {
		err := expr.UnmarshalGQL(123)
		assert.Error(t, err, "non-string input should return an error")
		assert.Contains(t, err.Error(), "cron expression must be a string", "the error message should be correct")
	})
}
