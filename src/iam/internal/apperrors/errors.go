package apperrors

import "errors"

var (
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrUserNotFound              = errors.New("user not found")
	ErrPasswordNotFound          = errors.New("password not found")
	ErrUUIDCreationFailed        = errors.New("failed to create UUID")
	ErrInvalidID                 = errors.New("invalid ID")
	ErrPasswordTooLong           = errors.New("password should be less than 72 bytes")
	ErrUnexpectedInternalFailure = errors.New("Unexepected internal failure")
	ErrFailedToCreateResource    = errors.New("failed to create resource")
)
