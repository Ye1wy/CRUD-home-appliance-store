package postgres

import "errors"

var (
	ErrProductNotFound = errors.New("Product not found in database")
	ErrClientNotFound  = errors.New("client not found")
	ErrQueryExection   = errors.New("query execution error")
)
