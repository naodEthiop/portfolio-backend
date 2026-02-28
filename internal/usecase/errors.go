package usecase

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrUpstream     = errors.New("upstream error")
	ErrRateLimited  = errors.New("rate limited")
)
