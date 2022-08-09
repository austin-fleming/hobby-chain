package store_ledger

import (
	"log"
	"os"
	"sync"

	"github.com/austin-fleming/hobby-chain/record"
	"github.com/austin-fleming/hobby-chain/store"
)

type ledgerStore struct {
	storagePath   string
	maxRecordSize int
	logger        *log.Logger
	index         *index
	writeMutex    sync.Mutex
}

func NewLedgerStore(storagePath string, maxRecordSize int, logger *log.Logger) ledgerStore {
	return ledgerStore{
		storagePath:   storagePath,
		maxRecordSize: maxRecordSize,
		logger:        logger,
	}
}

func (s *ledgerStore) Get(key string) ([]byte, error) {
	file, openErr := os.Open(s.storagePath)
	defer file.Close()
	if openErr != nil {
		return nil, store.ErrOpenFile(s.storagePath, openErr)
	}

	keyOffset, ok := s.index.Search(key)
	if !ok {
		return nil, store.ErrIndexNotFound(key, nil)
	}

	// jump to index
	_, seekErr := file.Seek(keyOffset, 0)
	if seekErr != nil {
		return nil, store.ErrIndexSeek(key, seekErr)
	}

	scanner, scanErr := record.NewScanner(file, s.maxRecordSize)
	if scanErr != nil {
		return nil, store.ErrScannerCreationError(scanErr)
	}

	var foundRecord *record.Record
	for scanner.Scan() {
		// TODO: dodged err
		record, _ := scanner.Record()
		if record.GetKey() == key {
			foundRecord = record
		}
	}

	if scanner.Err() != nil {
		return nil, store.NewStoreError("ErrScannerError", "Scanner Error", "", scanner.Err())
	}

	if foundRecord == nil || foundRecord.IsTombstone() {
		return nil, store.ErrNotFound(key)
	}

	return foundRecord.GetValue(), nil
}

func (s *ledgerStore) Write(key string, value []byte) error {
	record := record.NewValue(key, value)
	return s.append(record)
}

func (s *ledgerStore) Delete(key string) error {
	record := record.NewTombstone(key)
	return s.append(record)
}

func (s *ledgerStore) Close() error {
	s.logger.Print("Database close -- this doesn't actually do anything right now.")
	return nil
}

func (s *ledgerStore) append(rec *record.Record) error {
	if rec.Size() > s.maxRecordSize {
		return store.ErrRecordTooLarge(s.maxRecordSize, rec.Size())
	}

	s.writeMutex.Lock()
	defer s.writeMutex.Unlock()

	file, openErr := os.OpenFile(s.storagePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	defer file.Close()
	if openErr != nil {
		return store.ErrOpenFile(s.storagePath, openErr)
	}

	// TODO: Now that I'm calling it, this usage kind of usage is kind of inverted (in a bad way)
	bytesWritten, writeErr := rec.Write(file)
	if writeErr != nil {
		return store.ErrWriteFile(s.storagePath, writeErr)
	}

	if closeErr := file.Close(); closeErr != nil {
		return store.ErrFileClose(closeErr)
	}
	if syncErr := file.Sync(); syncErr != nil {
		return store.ErrFileSync(syncErr)
	}

	s.index.Insert(rec.GetKey(), int64(bytesWritten))

	return nil
}
