package object

import (
	"encoding/binary"

	"github.com/rotationalio/honu/pkg/store/lani"
	"github.com/rotationalio/honu/pkg/store/metadata"
)

// A compatibility indicator, increment this number any time the underlying storage is
// no longer compatible with the previous storage version and needs a different type of
// deserialization mechanism.
const StorageVersion uint8 = 1

// An object is the serial data that's written to the underlying storage and is composed
// of a one byte version indicator, metadata, and the document data serialized in a
// format that can be easily unmarshaled and marshed without requiring copying of data
// into multiple byte slices.
type Object []byte

// Create an object for storage by writing the metadata and the data into a new byte
// slice ready for storage on disk.
func Marshal(meta *metadata.Metadata, data []byte) (Object, error) {
	// Create an encoder of the correct size that requires no allocations.
	encoder := lani.Encoder{}
	encoder.Grow(1 + meta.Size() + len(data))

	// Write the storage version to the encoder
	if _, err := encoder.EncodeUint8(StorageVersion); err != nil {
		return nil, err
	}

	// Write the metadata to the encoder
	if _, err := encoder.EncodeStruct(meta); err != nil {
		return nil, err
	}

	// Write the fixed bytes to the encoder
	if _, err := encoder.EncodeFixed(data); err != nil {
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

func (o Object) Metadata() (*metadata.Metadata, error) {
	if o.StorageVersion() != StorageVersion {
		return nil, ErrBadVersion
	}

	i := o.metadataLength()
	if i < 1 {
		return nil, ErrNoMetadata
	}

	meta := &metadata.Metadata{}
	decoder := lani.NewDecoder(o[1 : i+1])
	if _, err := decoder.DecodeStruct(meta); err != nil {
		return nil, err
	}

	return meta, nil
}

func (o Object) Data() ([]byte, error) {
	if o.StorageVersion() != StorageVersion {
		return nil, ErrBadVersion
	}

	i := o.metadataLength()
	if i < 1 {
		return nil, ErrNoMetadata
	}

	return o[i+1:], nil
}

func (o Object) metadataLength() int {
	if len(o) == 0 {
		return -1
	}

	j := 1 + binary.MaxVarintLen64
	if j > len(o)-1 {
		j = len(o) - 1
	}

	rl, k := binary.Uvarint(o[1:j])
	if k <= 0 {
		return -1
	}

	return int(rl) + k
}
