package pkg

import "errors"

var (
	ErrDbInternal          = errors.New("db internal error")
	ErrEntityAlreadyExists = errors.New("entity already exists")
	ErrInvalidPayload      = errors.New("invalid payload")
)
