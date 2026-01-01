package author

import "errors"

var (
	ErrNotFound  = errors.New("author not found")
	ErrDuplicate = errors.New("author already exists")
)
