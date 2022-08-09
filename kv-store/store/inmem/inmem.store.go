package store_inmem

import (
	"fmt"
	"log"
	"sync"

	"github.com/austin-fleming/hobby-chain/store"
)

type Store struct {
	maxRecordSize int
	logger        *log.Logger
	sync.RWMutex
	table map[string][]byte
}

type StoreConfig struct {
	MaxRecordSize int
	Logger        *log.Logger
}

func NewStore(config StoreConfig) *Store {
	return &Store{
		maxRecordSize: config.MaxRecordSize,
		logger:        config.Logger,
		table:         map[string][]byte{},
	}
}

// implement methods for ../store interface
func (s *Store) Get(key string) ([]byte, error) {
	s.RLock()
	value, ok := s.table[key]
	s.RUnlock()
	if !ok {
		return nil, store.ErrNotFound(key)
	}

	return value, nil
}

// TODO: Should overwrites be allowed?
func (s *Store) Write(key string, value []byte) error {

	if kvSize := len([]byte(key)) + len(value); kvSize > s.maxRecordSize {
		return store.ErrBadRequest(fmt.Sprintf("key + value is %d bytes, exceeding the maximum size of %d bytes", kvSize, s.maxRecordSize))
	}

	s.Lock()
	s.table[key] = value
	s.Unlock()
	return nil
}

// TODO: Should this be disallowed?
// TODO: Should trying to delete a non-existent key error?
func (s *Store) Delete(key string) error {
	s.Lock()
	delete(s.table, key)
	s.Unlock()
	return nil
}

func (s *Store) Close() error {
	s.logger.Print("Db closed: in-memory database being used, therefore this is just a placeholder message.")
	return nil
}
