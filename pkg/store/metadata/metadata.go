package metadata

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/ulid"
)

//===========================================================================
// Metadata Storage
//===========================================================================

type Metadata struct {
	ObjectID     ulid.ULID        `json:"oid" msg:"oid"`
	CollectionID ulid.ULID        `json:"collection" msg:"collection"`
	Version      *Version         `json:"version" msg:"version"`
	Schema       *SchemaVersion   `json:"schema,omitempty" msg:"schema,omitempty"`
	MIME         string           `json:"mime" msg:"mime"`
	Owner        ulid.ULID        `json:"owner" msg:"owner"`
	Group        ulid.ULID        `json:"group" msg:"group"`
	Permissions  uint8            `json:"permissions" msg:"permissions"`
	ACL          []*AccessControl `json:"acl,omitempty" msg:"acl,omitempty"`
	WriteRegions []string         `json:"write_regions,omitempty" msg:"write_regions,omitempty"`
	Publisher    *Publisher       `json:"publisher,omitempty" msg:"publisher,omitempty"`
	Encryption   *Encryption      `json:"encryption,omitempty" msg:"encryption,omitempty"`
	Compression  *Compression     `json:"compression,omitempty" msg:"compression,omitempty"`
	Flags        uint8            `json:"flags" msg:"flags"`
	Created      time.Time        `json:"created" msg:"created"`
	Modified     time.Time        `json:"modified" msg:"modified"`
}

var _ lani.Encodable = &Metadata{}
var _ lani.Decodable = &Metadata{}

func (o *Metadata) Key() key.Key {
	// TODO: should the key be cached on the metadata to prevent multiple allocations?
	return key.New(o.CollectionID, o.ObjectID, &o.Version.Scalar)
}

func (o *Metadata) Size() (s int) {
	// ObjectID, CollectionID
	s += 16 + 16

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
	if m, err = e.EncodeULID(o.ObjectID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeULID(o.CollectionID); err != nil {
		return n + m, err
	}
	n += m

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

	if o.ObjectID, err = d.DecodeULID(); err != nil {
		return err
	}

	if o.CollectionID, err = d.DecodeULID(); err != nil {
		return err
	}

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
	PublicKeyID         string              `json:"public_key_id,omitempty" msg:"public_key_id,omitempty"`
	EncryptionKey       []byte              `json:"encryption_key,omitempty" msg:"encryption_key,omitempty"`
	HMACSecret          []byte              `json:"hmac_secret,omitempty" msg:"hmac_secret,omitempty"`
	Signature           []byte              `json:"signature,omitempty" msg:"signature,omitempty"`
	SealingAlgorithm    EncryptionAlgorithm `json:"sealing_algorithm,omitempty" msg:"sealing_algorithm,omitempty"`
	EncryptionAlgorithm EncryptionAlgorithm `json:"encryption_algoirthm" msg:"encryption_algorithm"`
	SignatureAlgorithm  EncryptionAlgorithm `json:"signature_algorithm,omitempty" msg:"signature_algorithm,omitempty"`
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

func ParseEncryptionAlgorithm(s string) (EncryptionAlgorithm, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case "PLAINTEXT":
		return Plaintext, nil
	case "AES256_GCM":
		return AES256_GCM, nil
	case "AES192_GCM":
		return AES192_GCM, nil
	case "AES128_GCM":
		return AES128_GCM, nil
	case "HMAC_SHA256":
		return HMAC_SHA256, nil
	case "RSA_OEAP_SHA512":
		return RSA_OEAP_SHA512, nil
	default:
		return 0, fmt.Errorf("%q is not a valid compression algorithm", s)
	}
}

func (o EncryptionAlgorithm) String() string {
	switch o {
	case Plaintext:
		return "PLAINTEXT"
	case AES256_GCM:
		return "AES256_GCM"
	case AES192_GCM:
		return "AES192_GCM"
	case AES128_GCM:
		return "AES128_GCM"
	case HMAC_SHA256:
		return "HMAC_SHA256"
	case RSA_OEAP_SHA512:
		return "RSA_OEAP_SHA512"
	default:
		return "UNKNOWN"
	}
}

func (o *EncryptionAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *EncryptionAlgorithm) UnmarshalJSON(data []byte) (err error) {
	var alg string
	if err := json.Unmarshal(data, &alg); err != nil {
		return err
	}
	if *o, err = ParseEncryptionAlgorithm(alg); err != nil {
		return err
	}
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
	Algorithm CompressionAlgorithm `json:"algorithm" msg:"algorithm"`
	Level     int64                `json:"level,omitempty" msg:"level,omitempty"`
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

func ParseCompressionAlgorithm(s string) (CompressionAlgorithm, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case "NONE":
		return None, nil
	case "GZIP":
		return GZIP, nil
	case "COMPRESS":
		return COMPRESS, nil
	case "DEFLATE":
		return DEFLATE, nil
	case "BROTLI":
		return BROTLI, nil
	default:
		return 0, fmt.Errorf("%q is not a valid compression algorithm", s)
	}
}

func (o CompressionAlgorithm) String() string {
	switch o {
	case None:
		return "NONE"
	case GZIP:
		return "GZIP"
	case COMPRESS:
		return "COMPRESS"
	case DEFLATE:
		return "DEFLATE"
	case BROTLI:
		return "BROTLI"
	default:
		return "UNKNOWN"
	}
}

func (o *CompressionAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *CompressionAlgorithm) UnmarshalJSON(data []byte) (err error) {
	var alg string
	if err := json.Unmarshal(data, &alg); err != nil {
		return err
	}
	if *o, err = ParseCompressionAlgorithm(alg); err != nil {
		return err
	}
	return nil
}
