package database

import "errors"

const (
	// databases
	DATABASE = "Store"
	// collections
	CLIENTS   = "clients"
	PRODUCTS  = "products"
	SUPPLIERS = "suppliers"
	IMAGES    = "images"
	ADDRESSES = "addresses"
)

var (
	ErrURLNotFound    = errors.New("[ERROR] Url not found")
	ErrURLExist       = errors.New("[ERROR] Url is exist")
	ErrConnectTimeout = errors.New("connect timeout")
)
