package database

import "errors"

var (
	ErrModeNotCorrect = errors.New("mode is not correct")
	ErrEntityNotFound = errors.New("entity not found")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrEntityNotFound)
}
