package crud_errors

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrQueryExection   = errors.New("query execution error")
	ErrInvalidParam    = errors.New("invalid parameter")
	ErrNoContent       = errors.New("No input content")
	ErrRepoIsExist     = errors.New("repository is exist")
	ErrRepoIsNotExitst = errors.New("repository is not exitst")
)
