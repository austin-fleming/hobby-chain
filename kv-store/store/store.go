package store

// TODO: would it better to pass Record instead of []byte?
type Store interface {
	Get(key string) ([]byte, error)
	Write(key string, value []byte) error
	Delete(key string) error
	Close() error
}
