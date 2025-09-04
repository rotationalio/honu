package object

import (
	"go.rtnl.ai/honu/pkg/store/lani"
)

// Create an object to store system data, which do not need metadata appended to the
// object since these objects are themselves metadata. System objects are used entirely
// for the store's internal management and are not user-facing.
func MarshalSystem(obj lani.Encodable) (_ Object, err error) {
	// Create an encoder of the correct size that requires no allocations.
	encoder := lani.Encoder{}
	encoder.Grow(2 + obj.Size())

	// Write the storage version to the encoder
	if _, err = encoder.EncodeUint8(StorageVersion); err != nil {
		return nil, err
	}

	// Append the object to the encoder
	if _, err = encoder.EncodeStruct(obj); err != nil {
		return nil, err
	}

	// Append nil metadata to the encoder, since system objects do not have metadata.
	if _, err = encoder.EncodeStruct(nil); err != nil {
		return nil, err
	}

	return Object(encoder.Bytes()), nil
}

// Unmarshal system metadata from an object; e.g. to unpack a system object from the
// bytes stored on disk. System objects do not have metadata and they are decodable,
// unlike user data which cannot be unmarshaled without a schema or metadata.
func UnmarshalSystem(obj Object, v lani.Decodable) (err error) {
	// NOTE: cannot use obj.Data() because the structure of a system object is an
	// storage version byte, followed by an encoded struct, then nil metadata. The
	// Data() method expects a raw byte slice with a length prefix.
	decoder := lani.NewDecoder(obj[1 : len(obj)-1])
	if _, err = decoder.DecodeStruct(v); err != nil {
		return err
	}
	return nil
}
