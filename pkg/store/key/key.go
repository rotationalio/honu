package key

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/rotationalio/honu/pkg/store/lamport"
	"go.rtnl.ai/ulid"
)

const (
	// The default size of a v1 object storage key.
	keySize int = 45

	// The version of the key for compatibility indication; increment this number any time
	// the underlying key data is no longer compatible with the previous version.
	keyVersion byte = 0x1
)

var (
	ErrBadVersion = errors.New("key is malformed: cannot decode specified version")
	ErrBadSize    = errors.New("key is malformed: incorrect size")
	ErrMalformed  = errors.New("key is malformed: cannot parse version components")
)

// Keys are used to store objects in the underlying key/value store. It is a 45 byte key
// that is composed of 16 byte object and collection IDs and a 4 byte uint32 and 8 byte
// uint64 representing the lamport scalar version number. The last byte indicates the
// key version and marshaling compatibility. There are no separator characters
// between the components of the key since all components are a fixed length.
//
// A key is structured as: collection::oid::vid::pid::keyVersion
//
// Note that the version is serialized differently than the lamport scalar in order to
// maintain lexicographic sorting of the the data.
type Key []byte

// Create a new key for the specified collection and object ID with the given version.
// If the version is nil, the key is treated as an object prefix and either the latest
// version of the object is returned or all versions related to the object.
func New(cid, oid ulid.ULID, vers *lamport.Scalar) Key {
	key := make([]byte, keySize)
	key[0] = keyVersion
	copy(key[1:17], cid[:])
	copy(key[17:33], oid[:])
	if vers != nil {
		binary.BigEndian.PutUint64(key[33:41], vers.VID)
		binary.BigEndian.PutUint32(key[41:45], vers.PID)
	}
	return Key(key)
}

// Returns the collection ID encoded in the key as a ulid.
func (k Key) CollectionID() ulid.ULID {
	if err := k.Check(); err != nil {
		panic(err)
	}
	return ulid.ULID(k[1:17])
}

// Returns the object ID encoded in the key as a ulid.
func (k Key) ObjectID() ulid.ULID {
	if err := k.Check(); err != nil {
		panic(err)
	}
	return ulid.ULID(k[17:33])
}

// Returns the version specified by the key if any (if no version is specified then
// returns a zero valued version rather than nil).
func (k Key) Version() lamport.Scalar {
	if err := k.Check(); err != nil {
		panic(err)
	}
	return lamport.Scalar{
		VID: binary.BigEndian.Uint64(k[33:41]),
		PID: binary.BigEndian.Uint32(k[41:45]),
	}
}

// ObjectPrefix returns the collection and object IDs without any version information.
func (k Key) ObjectPrefix() []byte {
	if err := k.Check(); err != nil {
		panic(err)
	}
	return k[0:33]
}

// ObjectLimit returns a byte slice with the collection and object ID with the last
// byte incremented by one. This can be used to create a range query for all versions
// of an object where the start is the ObjectPrefix.
func (k Key) ObjectLimit() []byte {
	if err := k.Check(); err != nil {
		panic(err)
	}
	limit := make([]byte, 33)
	copy(limit, k[0:33])
	limit[32]++
	return limit
}

// HasVersion checks if there is any version information or if only the object prefix
// is specified by the key. If false, then the Version() method is guaranteed to return
// a zero valued version. If true, then there is a specific version described by the key.
func (k Key) HasVersion() bool {
	if err := k.Check(); err != nil {
		panic(err)
	}

	for i := keySize - 1; i > 32; i-- {
		if k[i] != 0 {
			return true
		}
	}
	return false
}

func (k Key) Check() error {
	if len(k) != keySize {
		return ErrBadSize
	}

	if k[0] != keyVersion {
		return ErrBadVersion
	}

	return nil
}

//===========================================================================
// Sort Interface
//===========================================================================

type Keys []Key

func (k Keys) Len() int           { return len(k) }
func (k Keys) Less(i, j int) bool { return bytes.Compare(k[i][:], k[j][:]) < 0 }
func (k Keys) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
