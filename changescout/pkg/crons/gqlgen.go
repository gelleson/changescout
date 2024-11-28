package crons

import (
	"fmt"
	"io"
)

// MarshalGQL implements the graphql.Marshaler interface
func (e CronExpression) MarshalGQL(w io.Writer) {
	// Quote the string value
	io.WriteString(w, `"`)
	io.WriteString(w, string(e))
	io.WriteString(w, `"`)
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (e *CronExpression) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("cron expression must be a string")
	}

	// Validate the expression
	scheduler := NewScheduler()
	if err := scheduler.Validate(CronExpression(str)); err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	*e = CronExpression(str)
	return nil
}
