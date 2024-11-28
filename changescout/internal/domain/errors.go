package domain

import "errors"

var (
	ErrCheckNotFound   = errors.New("check not found")
	ErrDiffFailed      = errors.New("diff failed")
	ErrRequestFailed   = errors.New("request failed")
	ErrWebsiteNotFound = errors.New("website not found")
)

func IsErrCheckNotFound(err error) bool {
	return errors.Is(err, ErrCheckNotFound)
}
