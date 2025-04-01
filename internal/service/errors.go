package service

import "errors"

var (
	ErrCannotShortenURI = errors.New("cannot shorten uri")
	ErrCannotExpandURI  = errors.New("cannot expand uri")
	ErrURINotFound      = errors.New("uri not found")
)
