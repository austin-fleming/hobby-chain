package store

type Store interface {
	Get(key string) ([]byte, error)
	Write(key string, value []byte) error
	Delete(key string) error
	Close() error
}

type ErrNotFound struct {
	key string
}

type ErrBadRequest struct {
	reason string
}
