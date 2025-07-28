package crud_errors

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrQueryExection       = errors.New("query execution error")
	ErrInvalidParam        = errors.New("invalid parameter")
	ErrNoContent           = errors.New("no input content")
	ErrRepoIsExist         = errors.New("repository is exist")
	ErrRepoIsNotExitst     = errors.New("repository is not exitst")
	ErrAddressIsExist      = errors.New("address is exists")
	ErrAddressIsEmpty      = errors.New("address is empty")
	ErrForeignKeyViolation = errors.New("foragin key violation, someone links is alive")
	ErrImageCorruption     = errors.New("image is corrupted or input data is not image")
	ErrConversionProblem   = errors.New("conversion problem, panic awoided")
)
