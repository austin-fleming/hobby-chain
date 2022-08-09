package store

import "fmt"

// TODO: add op
type StoreError struct {
	kind    string
	message string
	reason  string
	inner   error
}

func NewStoreError(kind string, message string, reason string, inner error) StoreError {
	return StoreError{kind, message, reason, inner}
}

func (se StoreError) Error() string {
	return fmt.Sprintf("Store Error: %s | %s", se.message, se.reason)
}

func (se StoreError) Kind() string {
	return se.kind
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
	return NewStoreError("ErrNotFound", "Not Found", fmt.Sprintf("no entry for key '%s'", key), nil)
}
func ErrBadRequest(reason string) StoreError {
	return NewStoreError("ErrBadRequest", "Bad Request", reason, nil)
}
func ErrDeserialize(reason string) StoreError {
	return NewStoreError("ErrDeserialize", "Failed to Deserialize", reason, nil)
}
func ErrInsufficientData() StoreError {
	return NewStoreError("ErrInsufficientData", "Insufficient Data", "data missing from record", nil)
}
func ErrDataCorruption() StoreError {
	return NewStoreError("ErrDataCorruption", "Failed to Deserialize", "Checksums do not match", nil)
}
func ErrRecordTooLarge(maxSize int, recordSize int) StoreError {
	return NewStoreError("ErrRecordTooLarge", "Record too Large", fmt.Sprintf("max %d, record %d", maxSize, recordSize), nil)
}
func ErrOpenFile(filePath string, err error) StoreError {
	return NewStoreError("ErrOpenFile", "Failed to Open File", fmt.Sprintf("file at %s", filePath), err)
}
func ErrWriteFile(filePath string, err error) StoreError {
	return NewStoreError("ErrStoreWrite", "Failed to Write Record", fmt.Sprintf("file at %s", filePath), err)
}
func ErrScannerCreationError(err error) StoreError {
	return NewStoreError("ErrScannerCreationError", "Failed to Create Scanner", "", err)
}
func ErrScan(err error) StoreError {
	return NewStoreError("ErrScan", "Scanner Error", "An issue occurred when scanning value", err)
}
func ErrFileClose(err error) StoreError {
	return NewStoreError("ErrFileClose", "File Close Error", "", err)
}
func ErrFileSync(err error) StoreError {
	return NewStoreError("ErrFileSync", "File Sync Error", "", err)
}
func ErrIndexNotFound(key string, err error) StoreError {
	return NewStoreError("ErrIndexNotFound", "Index Not Found", fmt.Sprintf("no index for key: %s", key), err)
}
func ErrIndexSeek(key string, err error) StoreError {
	return NewStoreError("ErrIndexSeek", "Index Seek Failed", fmt.Sprintf("could not seek for item with the key: %s", key), err)
}
