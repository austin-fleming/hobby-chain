package record

import (
	"bufio"
	"encoding/binary"
	"hash/crc32"
	"io"

	"github.com/austin-fleming/hobby-chain/store"
)

func (r *Record) Serialize() []byte {

	keyLength := uint32(len([]byte(r.key)))
	keyBuffer := make([]byte, BL_KEY)
	binary.BigEndian.PutUint32(keyBuffer, keyLength)

	valueLength := uint32(len(r.value))
	valueBuffer := make([]byte, BL_VALUE)
	binary.BigEndian.PutUint32(valueBuffer, valueLength)

	// TODO: this currently doesn't have a fixed size or chunking
	recordData := [][]byte{{r.kind}, keyBuffer, valueBuffer, []byte(r.key), r.value}
	recordDataBuffer := []byte{}
	checksum := crc32.NewIEEE()
	for _, item := range recordData {
		recordDataBuffer = append(recordDataBuffer, item...)
		checksum.Write(item)
	}

	checksumBuffer := make([]byte, BL_CHECKSUM)
	binary.BigEndian.PutUint32(checksumBuffer, checksum.Sum32())

	return append(checksumBuffer, recordDataBuffer...)
}

func DeserializeRecord(serialData []byte) (*Record, error) {
	// Should absolutely not be the case
	if len(serialData) < BL_TOTAL {
		return nil, store.ErrInsufficientData()
	}

	// get data
	checksum := binary.BigEndian.Uint32(serialData[OFFSET_CHECKSUM : OFFSET_CHECKSUM+BL_CHECKSUM])
	kind := serialData[OFFSET_KIND]
	keyLength := int(binary.BigEndian.Uint32(serialData[OFFSET_KEY : OFFSET_KEY+BL_KEY]))
	valueLength := int(binary.BigEndian.Uint32(serialData[OFFSET_VALUE : OFFSET_VALUE+BL_VALUE]))

	if len(serialData) < BL_TOTAL+keyLength+valueLength {
		return nil, store.ErrInsufficientData()
	}

	// find kv cursors
	keyStart := BL_TOTAL
	keyEnd := keyStart + keyLength
	valueStart := keyEnd
	valueEnd := valueStart + int(valueLength)

	// verify data
	checksumCheck := crc32.NewIEEE()
	checksumCheck.Write(serialData[OFFSET_KIND:valueEnd])
	if checksumCheck.Sum32() != checksum {
		return nil, store.ErrDataCorruption()
	}

	// get kv pair
	key := make([]byte, keyLength)
	value := make([]byte, valueLength)
	copy(key, serialData[keyStart:keyEnd])
	copy(value, serialData[valueStart:valueEnd])

	return &Record{kind: kind, key: string(key), value: value}, nil
}

func (r *Record) Write(w io.Writer) (int, error) {
	data := r.Serialize()
	return w.Write(data)
}

// --------------
// READER
// --------------
type Scanner struct {
	*bufio.Scanner
}

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// TODO: this is nasty. Find a better way to signal a particular error type.
	record, deserializeErr := DeserializeRecord(data)
	if err, ok := deserializeErr.(store.StoreError); ok && err.Kind() == "ErrInsufficientData" {
		return 0, nil, nil
	}

	if deserializeErr != nil {
		return 0, nil, deserializeErr
	}

	advanceAmount := record.Size()

	return advanceAmount, data[:advance], nil
}

func NewScanner(r io.Reader, maxScanSize int) (*Scanner, error) {
	scanner := bufio.NewScanner(r)

	buffer := make([]byte, 4096)
	// Make sure buffer is always at least large enough to consume headers
	scanner.Buffer(buffer, maxScanSize+BL_TOTAL)
	scanner.Split(split)

	return &Scanner{scanner}, nil
}

func (scanner *Scanner) Record() (*Record, error) {
	data := scanner.Bytes()
	record, deserializeErr := DeserializeRecord(data)
	if deserializeErr != nil {
		return nil, deserializeErr
	}

	return record, nil
}
