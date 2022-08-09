package record

/*
FORMAT:

| length	| name 	| description 			|
|-----------|-------|-----------------------|
| 4 bytes 	| sum 	| checksum 				|
| 1 bytes 	| kind 	| Active or Tombstone 	|
| 4 bytes 	| kL 	| key length 			|
| 4 bytes 	| vL 	| value length 			|
| 1..n 		| key   | key 					|
| 0..n 		| val 	| value 				|

*/
// BYTE LENGTHS
const (
	BL_KIND     = 1
	BL_KEY      = 4
	BL_VALUE    = 4
	BL_CHECKSUM = 4
	BL_TOTAL    = BL_KIND + BL_CHECKSUM + BL_KEY + BL_VALUE
)

// OFFSETS
const (
	OFFSET_CHECKSUM = 0
	OFFSET_KIND     = OFFSET_CHECKSUM + BL_CHECKSUM
	OFFSET_KEY      = OFFSET_KIND + BL_KIND
	OFFSET_VALUE    = OFFSET_KEY + BL_KEY
)

const (
	KIND_VALUE = iota
	KIND_TOMBSTONE
)

type Record struct {
	kind  byte
	key   string
	value []byte
}

func NewValue(key string, value []byte) *Record {
	return &Record{
		kind:  KIND_VALUE,
		key:   key,
		value: value,
	}
}

func NewTombstone(key string) *Record {
	return &Record{
		kind:  KIND_TOMBSTONE,
		key:   key,
		value: []byte{},
	}
}

func (r *Record) GetKey() string {
	return r.key
}
func (r *Record) GetValue() []byte {
	return r.value
}
func (r *Record) IsTombstone() bool {
	return r.kind == KIND_TOMBSTONE
}
func (r *Record) Size() int {
	return BL_TOTAL + len(r.key) + len(r.value)
}
