package store

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/oklog/ulid"
)

//===========================================================================
// Object Storage
//===========================================================================

type Object struct {
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
	Data         []byte
}

func (o *Object) Size() (s int) {
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
	s += len(o.Data)
	return
}

func (o *Object) Encode(e *Encoder) (n int, err error) {
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

	if m, err = e.Encode(o.Data); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Object) Decode(d *Decoder) (err error) {
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

	if o.Data, err = d.Decode(); err != nil {
		return err
	}

	return nil
}

//===========================================================================
// Object Version
//===========================================================================

type Version struct {
	PID       uint64
	Version   uint64
	Region    string
	Parent    *Version
	Tombstone bool
	Created   time.Time
}

func (o *Version) Size() (s int) {
	s += 2 * binary.MaxVarintLen64
	s += len([]byte(o.Region)) + binary.MaxVarintLen64

	if o.Parent != nil {
		s += o.Parent.Size() + 1 // Add 1 for the not nil bool
	} else {
		s += 1 // Add 1 for the nil bool
	}

	s += 1                     // Tombstone bool
	s += binary.MaxVarintLen64 // Timestamp int64

	return
}

func (o *Version) Encode(e *Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeUint64(o.PID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint64(o.Version); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeString(o.Region); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Parent); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeBool(o.Tombstone); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeTime(o.Created); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Version) Decode(d *Decoder) (err error) {
	if o.PID, err = d.DecodeUint64(); err != nil {
		return err
	}

	if o.Version, err = d.DecodeUint64(); err != nil {
		return err
	}

	if o.Region, err = d.DecodeString(); err != nil {
		return err
	}

	var isNil bool
	o.Parent = &Version{}

	if isNil, err = d.DecodeStruct(o.Parent); err != nil {
		return err
	}

	if isNil {
		o.Parent = nil
	}

	if o.Tombstone, err = d.DecodeBool(); err != nil {
		return err
	}

	if o.Created, err = d.DecodeTime(); err != nil {
		return err
	}

	return nil
}

//===========================================================================
// Schema Version
//===========================================================================

type SchemaVersion struct {
	Name  string
	Major uint32
	Minor uint32
	Patch uint32
}

func (o *SchemaVersion) Size() (s int) {
	s += len([]byte(o.Name)) + binary.MaxVarintLen64
	s += 3 * binary.MaxVarintLen32
	return
}

func (o *SchemaVersion) Encode(e *Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeString(o.Name); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint32(o.Major); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint32(o.Minor); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint32(o.Patch); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *SchemaVersion) Decode(d *Decoder) (err error) {
	if o.Name, err = d.DecodeString(); err != nil {
		return err
	}

	if o.Major, err = d.DecodeUint32(); err != nil {
		return err
	}

	if o.Minor, err = d.DecodeUint32(); err != nil {
		return err
	}

	if o.Patch, err = d.DecodeUint32(); err != nil {
		return err
	}

	return nil
}

//===========================================================================
// ACL
//===========================================================================

type AccessControl struct {
	ClientID    ulid.ULID
	Permissions uint8
}

func (o *AccessControl) Size() int {
	// ULID + 1 byte
	return 17
}

func (o *AccessControl) Encode(e *Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeULID(o.ClientID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(o.Permissions); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *AccessControl) Decode(d *Decoder) (err error) {
	if o.ClientID, err = d.DecodeULID(); err != nil {
		return err
	}

	if o.Permissions, err = d.DecodeUint8(); err != nil {
		return err
	}

	return nil
}

//===========================================================================
// Provenance
//===========================================================================

type Publisher struct {
	PublisherID ulid.ULID
	ClientID    ulid.ULID
	IPAddress   net.IP
	UserAgent   string
}

func (o *Publisher) Size() (s int) {
	// 2 ULIDs and 2 variable byte arrays
	s += 16 + 16
	s += len(o.IPAddress) + binary.MaxVarintLen64
	s += len([]byte(o.UserAgent)) + binary.MaxVarintLen64
	return
}

func (o *Publisher) Encode(e *Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeULID(o.PublisherID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeULID(o.ClientID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.Encode(o.IPAddress); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeString(o.UserAgent); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Publisher) Decode(d *Decoder) (err error) {
	if o.PublisherID, err = d.DecodeULID(); err != nil {
		return err
	}

	if o.ClientID, err = d.DecodeULID(); err != nil {
		return err
	}

	var ip []byte
	if ip, err = d.Decode(); err != nil {
		return err
	}
	o.IPAddress = net.IP(ip)

	if o.UserAgent, err = d.DecodeString(); err != nil {
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

func (o *Encryption) Size() (s int) {
	s += len([]byte(o.PublicKeyID)) + binary.MaxVarintLen64
	s += len(o.EncryptionKey) + binary.MaxVarintLen64
	s += len(o.HMACSecret) + binary.MaxVarintLen64
	s += len(o.Signature) + binary.MaxVarintLen64
	s += 3 // the three encryption algorithm bytes

	return
}

func (o *Encryption) Encode(e *Encoder) (n int, err error) {
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

func (o *Encryption) Decode(d *Decoder) (err error) {
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

func (o *Compression) Size() int {
	return 1 + binary.MaxVarintLen64
}

func (o *Compression) Encode(e *Encoder) (n int, err error) {
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

func (o *Compression) Decode(d *Decoder) (err error) {
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
