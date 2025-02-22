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

func New(cid, oid ulid.ULID, vers lamport.Scalar) Key {
	key := make([]byte, keySize)
	copy(key[0:16], cid[:])
	copy(key[16:32], oid[:])
	binary.BigEndian.PutUint64(key[32:40], vers.VID)
	binary.BigEndian.PutUint32(key[40:44], vers.PID)
	key[44] = keyVersion
	return Key(key)
}

func (k Key) CollectionID() ulid.ULID {
	if err := k.Check(); err != nil {
		panic(err)
	}
	return ulid.ULID(k[0:16])
}

func (k Key) ObjectID() ulid.ULID {
	if err := k.Check(); err != nil {
		panic(err)
	}
	return ulid.ULID(k[16:32])
}

func (k Key) Version() lamport.Scalar {
	if err := k.Check(); err != nil {
		panic(err)
	}
	return lamport.Scalar{
		VID: binary.BigEndian.Uint64(k[32:40]),
		PID: binary.BigEndian.Uint32(k[40:44]),
	}
}

func (k Key) Check() error {
	if len(k) != keySize {
		return ErrBadSize
	}

	if k[44] != keyVersion {
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
