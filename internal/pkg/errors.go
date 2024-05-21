package pkg

import "errors"

var (
	ErrDbInternal           = errors.New("db internal error")
	ErrEntityAlreadyExists  = errors.New("entity already exists")
	ErrInvalidPayload       = errors.New("invalid payload")
	ErrNotFound             = errors.New("entity not found")
	ErrInvalidUriParameters = errors.New("invalid uri paramteres")
	ErrEntityAlreadyDeleted = errors.New("entity already delted")
)
