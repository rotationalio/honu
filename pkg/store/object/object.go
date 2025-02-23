package object

import (
	"encoding/binary"

	"github.com/rotationalio/honu/pkg/store/key"
	"github.com/rotationalio/honu/pkg/store/lani"
	"github.com/rotationalio/honu/pkg/store/metadata"
)

// A compatibility indicator, increment this number any time the underlying storage is
// no longer compatible with the previous storage version and needs a different type of
// deserialization mechanism.
const StorageVersion uint8 = 1

// An object is the serial data that's written to the underlying storage and is composed
// of a one byte version indicator, the length of the document data, the document data,
// and the metadata serialized in a format that can be easily unmarshaled and marshed
// without requiring copying of data into multiple byte slices.
type Object []byte

// Create an object for storage by writing the metadata and the data into a new byte
// slice ready for storage on disk.
func Marshal(meta *metadata.Metadata, data []byte) (_ Object, err error) {
	// Create an encoder of the correct size that requires no allocations.
	encoder := lani.Encoder{}
	encoder.Grow(1 + binary.MaxVarintLen64 + meta.Size() + len(data))

	// Write the storage version to the encoder
	if _, err = encoder.EncodeUint8(StorageVersion); err != nil {
		return nil, err
	}

	// Write the bytes to the encoder
	if _, err = encoder.Encode(data); err != nil {
		return nil, err
	}

	// Append the metadata to the encoder
	if _, err = encoder.EncodeStruct(meta); err != nil {
		return nil, err
	}

	return Object(encoder.Bytes()), nil
}

func (o Object) StorageVersion() uint8 {
	if len(o) == 0 {
		return 0
	}
	return uint8(o[0])
}

// Shortcut for parsing the metadata then getting the key from it. However, it is not
// recommended to do this if you need access to the metadata since that will require
// parsing the metadata twice.
func (o Object) Key() (_ key.Key, err error) {
	// TODO: parse the key from the metadata without parsing the entire struct.
	var meta *metadata.Metadata
	if meta, err = o.Metadata(); err != nil {
		return nil, err
	}
	return meta.Key(), nil
}

func (o Object) Metadata() (*metadata.Metadata, error) {
	if o.StorageVersion() != StorageVersion {
		return nil, ErrBadVersion
	}

	d, b := o.dataLength()
	if d < 0 {
		return nil, ErrMalformed
	}

	meta := &metadata.Metadata{}
	decoder := lani.NewDecoder(o[1+d+b:])
	if _, err := decoder.DecodeStruct(meta); err != nil {
		return nil, err
	}

	return meta, nil
}

func (o Object) Data() ([]byte, error) {
	if o.StorageVersion() != StorageVersion {
		return nil, ErrBadVersion
	}

	d, b := o.dataLength()
	switch {
	case d < 0:
		return nil, ErrMalformed
	case d == 0:
		return nil, nil
	default:
		return o[1+b : 1+b+d], nil
	}
}

// If true the object is a tombstone meaning that it only contains metadata and has no
// associated data. Tombstones are used to indicate that a key has been deleted.
func (o Object) Tombstone() bool {
	if o.StorageVersion() != StorageVersion {
		return false
	}

	d, _ := o.dataLength()
	return d == 0
}

func (o Object) dataLength() (int, int) {
	if len(o) == 0 {
		return -1, -1
	}

	j := 1 + binary.MaxVarintLen64
	if j > len(o)-1 {
		j = len(o) - 1
	}

	if j < 1 {
		return -1, -1
	}

	rl, k := binary.Uvarint(o[1:j])
	if k <= 0 {
		return -1, -1
	}

	return int(rl), k
}
