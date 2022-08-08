package store

import "fmt"

type StoreError struct {
	message string
	reason  string
	inner   error
}

func NewStoreError(message string, reason string, inner error) StoreError {
	return StoreError{message, reason, inner}
}

func (se StoreError) Error() string {
	return fmt.Sprintf("Store Error: %s | %s", se.message, se.reason)
}

func (se StoreError) InnerError() string {
	return se.inner.Error()
}

func (se *StoreError) SetInner(err error) {
	se.inner = err
}

// -----------------
// ERRORS
// -----------------

func ErrNotFound(key string) StoreError {
	return NewStoreError("Not Found", fmt.Sprintf("no entry for key '%s'", key), nil)
}
func ErrBadRequest(reason string) StoreError {
	return NewStoreError("Bad Request", reason, nil)
}
