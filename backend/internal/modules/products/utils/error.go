package utils

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request body")
	ErrInvalidId      = errors.New("invalid id")

	ErrInvalidCode = errors.New("invalid code39 value")
	ErrDuplicate   = errors.New("product code already exists")
	ErrNotFound    = errors.New("barcode not found")

	ErrList   = errors.New("failed to list barcodes")
	ErrCreate = errors.New("failed to create barcode")
	ErrDelete = errors.New("failed to delete barcode")
)
