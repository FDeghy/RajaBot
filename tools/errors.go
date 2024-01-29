package tools

import "errors"

var (
	ErrSubNotFound  = errors.New("subscription not found")
	ErrAlreadyTrial = errors.New("trial already used")
)
