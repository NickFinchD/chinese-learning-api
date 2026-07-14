package review

import "errors"

var (
	ErrWordNotFound     = errors.New("word not found")
	ErrWordNotReady     = errors.New("word is not ready for review")
	ErrWordAlreadyAdded = errors.New("word already added")
)
