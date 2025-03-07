package main

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternalServer = errors.New("internal server error")
)