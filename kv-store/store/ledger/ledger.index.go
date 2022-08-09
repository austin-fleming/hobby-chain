package store_ledger

import (
	"os"
	"sync"

	"github.com/austin-fleming/hobby-chain/record"
	"github.com/austin-fleming/hobby-chain/store"
)

type index struct {
	mu     sync.RWMutex
	table  map[string]int64
	cursor int64
}

func NewIndex() *index {
	return &index{}
}

func (i *index) Search(key string) (int64, bool) {
	i.mu.RLock()
	val, ok := i.table[key]
	i.mu.RUnlock()
	return val, ok
}

func (i *index) Insert(key string, bytesWritten int64) {
	i.mu.Lock()
	defer i.mu.Unlock()
	// set insertion point to current cursor
	i.table[key] = i.cursor
	// TODO: currently this is getting retroactively used, which is more fragile.
	// update cursor for next call.
	i.cursor += bytesWritten
}

func (i *index) LoadIndex(filePath string, maxRecordSize int) (*index, error) {
	idx := new(index)

	file, openErr := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0600)
	defer file.Close()
	if openErr != nil {
		return nil, store.ErrOpenFile(filePath, openErr)
	}

	scanner, scannerErr := record.NewScanner(file, maxRecordSize)
	if scannerErr != nil {
		return nil, scannerErr
	}

	for scanner.Scan() {
		parsedRecord, scanErr := scanner.Record()
		if scanErr != nil {
			return nil, store.ErrDeserialize("Failed to parse record when building index")
		}

		idx.Insert(parsedRecord.GetKey(), int64(parsedRecord.Size()))
	}

	if scanner.Err() != nil {
		return nil, store.ErrScan(scanner.Err())
	}

	return idx, nil
}
