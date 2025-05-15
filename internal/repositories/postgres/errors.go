package postgres

import "errors"

var (
	ErrNotFound      = errors.New("item not found")
	ErrQueryExection = errors.New("query execution error")
)
