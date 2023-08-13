package storage

import "errors"

var (
	ErrorUrlNotFound = errors.New("Url not found")
	ErrorUrlExists   = errors.New("Url exists")
)
