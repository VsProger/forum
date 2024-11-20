package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")

type Err struct {
	Text_err string
	Code_err int
}

var (
	ErrInvalidComment  error = errors.New("invalid length of text")
	ErrEmptyComment    error = errors.New("empty comment")
	ErrNotAscii        error = errors.New("text is not in Ascii")
	ErrUserNotFound    error = errors.New("user not found")
	ErrInvalidPassword error = errors.New("Password does not match")
)
