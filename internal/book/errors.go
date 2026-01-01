package book

import "errors"

var (
	ErrNotFound        = errors.New("book not found")
	ErrDuplicate       = errors.New("book already exists")
	ErrInvalidAuthorId = errors.New("invalid author ID")
)
