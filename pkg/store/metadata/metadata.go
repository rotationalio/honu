package metadata

import (
	"encoding/binary"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/honu/pkg/store/lani"
)

//===========================================================================
// Metadata Storage
//===========================================================================

type Metadata struct {
	Version      *Version
	Schema       *SchemaVersion
	MIME         string
	Owner        ulid.ULID
	Group        ulid.ULID
	Permissions  uint8
	ACL          []*AccessControl
	WriteRegions []string
	Publisher    *Publisher
	Encryption   *Encryption
	Compression  *Compression
	Flags        uint8
	Created      time.Time
	Modified     time.Time
}

var _ lani.Encodable = &Metadata{}
var _ lani.Decodable = &Metadata{}

func (o *Metadata) Size() (s int) {
	// Version size + not nil bool
	s += 1
	if o.Version != nil {
		s += o.Version.Size()
	}

	// SchemaVersion size + not nil bool
	s += 1
	if o.Schema != nil {
		s += o.Schema.Size()
	}

	s += len([]byte(o.MIME)) + binary.MaxVarintLen64
	s += 16 + 16 + 1 // Owner, Group, Permissions

	// ACL
	s += (len(o.ACL) + 1) * binary.MaxVarintLen64
	for _, ac := range o.ACL {
		s += 1
		if ac != nil {
			s += ac.Size()
		}
	}

	// Write Regions
	s += (len(o.WriteRegions) + 1) * binary.MaxVarintLen64
	for _, wr := range o.WriteRegions {
		s += len([]byte(wr))
	}

	// Publisher size + not nil bool
	s += 1
	if o.Publisher != nil {
		s += o.Publisher.Size()
	}

	// Encryption size + not nil bool
	s += 1
	if o.Encryption != nil {
		s += o.Encryption.Size()
	}

	// Compression size + not nil bool
	s += 1
	if o.Compression != nil {
		s += o.Compression.Size()
	}

	s += 1                         // Flags
	s += 2 * binary.MaxVarintLen64 // Created, Modified
	return
}

func (o *Metadata) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeStruct(o.Version); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Schema); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeString(o.MIME); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeULID(o.Owner); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeULID(o.Group); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(o.Permissions); err != nil {
		return n + m, err
	}
	n += m

	// Encode ACL length
	if m, err = e.EncodeUint64(uint64(len(o.ACL))); err != nil {
		return n + m, err
	}
	n += m

	// Encode each ACL
	for _, ac := range o.ACL {
		if m, err = e.EncodeStruct(ac); err != nil {
			return n + m, err
		}
		n += m
	}

	if m, err = e.EncodeStringSlice(o.WriteRegions); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Publisher); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Encryption); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Compression); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(o.Flags); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeTime(o.Created); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeTime(o.Modified); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Metadata) Decode(d *lani.Decoder) (err error) {
	// Setup nested structs
	o.Version = &Version{}
	o.Schema = &SchemaVersion{}
	o.Publisher = &Publisher{}
	o.Encryption = &Encryption{}
	o.Compression = &Compression{}

	var isNil bool
	if isNil, err = d.DecodeStruct(o.Version); err != nil {
		return err
	} else if isNil {
		o.Version = nil
	}

	if isNil, err = d.DecodeStruct(o.Schema); err != nil {
		return err
	} else if isNil {
		o.Schema = nil
	}

	if o.MIME, err = d.DecodeString(); err != nil {
		return err
	}

	if o.Owner, err = d.DecodeULID(); err != nil {
		return err
	}

	if o.Group, err = d.DecodeULID(); err != nil {
		return err
	}

	if o.Permissions, err = d.DecodeUint8(); err != nil {
		return err
	}

	// Read the number of ACLs stored
	var nACLs uint64
	if nACLs, err = d.DecodeUint64(); err != nil {
		return err
	}

	// Decode all ACLs
	if nACLs > 0 {
		o.ACL = make([]*AccessControl, nACLs)
		for i := uint(0); i < uint(nACLs); i++ {
			o.ACL[i] = &AccessControl{}
			if isNil, err = d.DecodeStruct(o.ACL[i]); err != nil {
				return err
			}
			if isNil {
				o.ACL[i] = nil
			}
		}
	}

	if o.WriteRegions, err = d.DecodeStringSlice(); err != nil {
		return err
	}

	if isNil, err = d.DecodeStruct(o.Publisher); err != nil {
		return err
	} else if isNil {
		o.Publisher = nil
	}

	if isNil, err = d.DecodeStruct(o.Encryption); err != nil {
		return err
	} else if isNil {
		o.Encryption = nil
	}

	if isNil, err = d.DecodeStruct(o.Compression); err != nil {
		return err
	} else if isNil {
		o.Compression = nil
	}

	if o.Flags, err = d.DecodeUint8(); err != nil {
		return err
	}

	if o.Created, err = d.DecodeTime(); err != nil {
		return err
	}

	if o.Modified, err = d.DecodeTime(); err != nil {
		return err
	}

	return nil
}

//===========================================================================
// Encryption
//===========================================================================

type EncryptionAlgorithm uint8

const (
	Plaintext EncryptionAlgorithm = iota
	AES256_GCM
	AES192_GCM
	AES128_GCM
	HMAC_SHA256
	RSA_OEAP_SHA512
)

type Encryption struct {
	PublicKeyID         string
	EncryptionKey       []byte
	HMACSecret          []byte
	Signature           []byte
	SealingAlgorithm    EncryptionAlgorithm
	EncryptionAlgorithm EncryptionAlgorithm
	SignatureAlgorithm  EncryptionAlgorithm
}

var _ lani.Encodable = &Encryption{}
var _ lani.Decodable = &Encryption{}

func (o *Encryption) Size() (s int) {
	s += len([]byte(o.PublicKeyID)) + binary.MaxVarintLen64
	s += len(o.EncryptionKey) + binary.MaxVarintLen64
	s += len(o.HMACSecret) + binary.MaxVarintLen64
	s += len(o.Signature) + binary.MaxVarintLen64
	s += 3 // the three encryption algorithm bytes

	return
}

func (o *Encryption) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeString(o.PublicKeyID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.Encode(o.EncryptionKey); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.Encode(o.HMACSecret); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.Encode(o.Signature); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(uint8(o.SealingAlgorithm)); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(uint8(o.EncryptionAlgorithm)); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(uint8(o.SignatureAlgorithm)); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Encryption) Decode(d *lani.Decoder) (err error) {
	if o.PublicKeyID, err = d.DecodeString(); err != nil {
		return err
	}

	if o.EncryptionKey, err = d.Decode(); err != nil {
		return err
	}

	if o.HMACSecret, err = d.Decode(); err != nil {
		return err
	}

	if o.Signature, err = d.Decode(); err != nil {
		return err
	}

	var a uint8
	if a, err = d.DecodeUint8(); err != nil {
		return err
	}
	o.SealingAlgorithm = EncryptionAlgorithm(a)

	if a, err = d.DecodeUint8(); err != nil {
		return err
	}
	o.EncryptionAlgorithm = EncryptionAlgorithm(a)

	if a, err = d.DecodeUint8(); err != nil {
		return err
	}
	o.SignatureAlgorithm = EncryptionAlgorithm(a)

	return nil
}

//===========================================================================
// Compression
//===========================================================================

type CompressionAlgorithm uint8

const (
	None CompressionAlgorithm = iota
	GZIP
	COMPRESS
	DEFLATE
	BROTLI
)

type Compression struct {
	Algorithm CompressionAlgorithm
	Level     int64
}

var _ lani.Encodable = &Compression{}
var _ lani.Decodable = &Compression{}

func (o *Compression) Size() int {
	return 1 + binary.MaxVarintLen64
}

func (o *Compression) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeUint8(uint8(o.Algorithm)); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeInt64(o.Level); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Compression) Decode(d *lani.Decoder) (err error) {
	var a uint8
	if a, err = d.DecodeUint8(); err != nil {
		return err
	}
	o.Algorithm = CompressionAlgorithm(a)

	if o.Level, err = d.DecodeInt64(); err != nil {
		return err
	}

	return nil
}
