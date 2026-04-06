package repository

import "errors"

var (
	ErrNotFound               = errors.New("not found")
	ErrShortLinkAlreadyExists = errors.New("short link already exists")
	ErrLongLinkAlreadyExists  = errors.New("long link already exists")
)
