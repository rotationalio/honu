package store

import (
	"net"
	"time"

	"github.com/oklog/ulid"
)

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

func (o Object) Write(w *Writer) (n int, err error) {
	var m int
	if m, err = w.WriteBool(o.Version != nil); err != nil {
		return n + m, err
	}
	n += m

	if o.Version != nil {
		if m, err = o.Version.Write(w); err != nil {
			return n + m, err
		}
		n += m
	}

	if m, err = w.WriteBool(o.Schema != nil); err != nil {
		return n + m, err
	}
	n += m

	if o.Schema != nil {
		if m, err = o.Schema.Write(w); err != nil {
			return n + m, err
		}
		n += m
	}

	if m, err = w.WriteString(o.MIME); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteULID(o.Owner); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteULID(o.Group); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteUint8(o.Permissions); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteUint64(uint64(len(o.ACL))); err != nil {
		return n + m, err
	}
	n += m

	for _, acl := range o.ACL {
		// TODO: nil checks?
		if m, err = acl.Write(w); err != nil {
			return n + m, err
		}
		n += m
	}

	if m, err = w.WriteStrings(o.WriteRegions); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteBool(o.Publisher != nil); err != nil {
		return n + m, err
	}
	n += m

	if o.Publisher != nil {
		if m, err = o.Publisher.Write(w); err != nil {
			return n + m, err
		}
		n += m
	}

	if m, err = w.WriteBool(o.Encryption != nil); err != nil {
		return n + m, err
	}
	n += m

	if o.Encryption != nil {
		if m, err = o.Encryption.Write(w); err != nil {
			return n + m, err
		}
		n += m
	}

	if m, err = w.WriteBool(o.Compression != nil); err != nil {
		return n + m, err
	}
	n += m

	if o.Compression != nil {
		if m, err = o.Compression.Write(w); err != nil {
			return n + m, err
		}
		n += m
	}

	if m, err = w.WriteUint8(o.Flags); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteTime(o.Created); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteTime(o.Modified); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteBytes(o.Data); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Object) Read(r *Reader) (err error) {
	var hasStruct bool
	if hasStruct, err = r.ReadBool(); err != nil {
		return err
	}

	if hasStruct {
		o.Version = &Version{}
		if err = o.Version.Read(r); err != nil {
			return err
		}
	} else {
		o.Version = nil
	}

	if hasStruct, err = r.ReadBool(); err != nil {
		return err
	}

	if hasStruct {
		o.Schema = &SchemaVersion{}
		if err = o.Schema.Read(r); err != nil {
			return err
		}
	} else {
		o.Schema = nil
	}

	if o.MIME, err = r.ReadString(); err != nil {
		return err
	}

	if o.Owner, err = r.ReadULID(); err != nil {
		return err
	}

	if o.Group, err = r.ReadULID(); err != nil {
		return err
	}

	if o.Permissions, err = r.ReadUint8(); err != nil {
		return err
	}

	var nACLs uint64
	if nACLs, err = r.ReadUint64(); err != nil {
		return err
	}

	if nACLs > 0 {
		o.ACL = make([]*AccessControl, nACLs)
		for i := uint64(0); i < nACLs; i++ {
			// TODO: check for nil?
			o.ACL[i] = &AccessControl{}
			if err = o.ACL[i].Read(r); err != nil {
				return err
			}
		}
	} else {
		o.ACL = nil
	}

	if o.WriteRegions, err = r.ReadStrings(); err != nil {
		return err
	}

	if hasStruct, err = r.ReadBool(); err != nil {
		return err
	}

	if hasStruct {
		o.Publisher = &Publisher{}
		if err = o.Publisher.Read(r); err != nil {
			return err
		}
	} else {
		o.Publisher = nil
	}

	if hasStruct, err = r.ReadBool(); err != nil {
		return err
	}

	if hasStruct {
		o.Encryption = &Encryption{}
		if err = o.Encryption.Read(r); err != nil {
			return err
		}
	} else {
		o.Encryption = nil
	}

	if hasStruct, err = r.ReadBool(); err != nil {
		return err
	}

	if hasStruct {
		o.Compression = &Compression{}
		if err = o.Compression.Read(r); err != nil {
			return err
		}
	} else {
		o.Compression = nil
	}

	if o.Flags, err = r.ReadUint8(); err != nil {
		return err
	}

	if o.Created, err = r.ReadTime(); err != nil {
		return err
	}

	if o.Modified, err = r.ReadTime(); err != nil {
		return err
	}

	if o.Data, err = r.ReadBytes(); err != nil {
		return err
	}

	return
}

type Version struct {
	PID       uint64
	Version   uint64
	Region    string
	Parent    *Version
	Tombstone bool
	Created   time.Time
}

func (v Version) Write(w *Writer) (n int, err error) {
	var m int
	if m, err = w.WriteUint64(v.PID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteUint64(v.Version); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteString(v.Region); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteBool(v.Parent != nil); err != nil {
		return n + m, err
	}
	n += m

	if v.Parent != nil {
		if m, err = v.Parent.Write(w); err != nil {
			return n + m, err
		}
		n += m
	}

	if m, err = w.WriteBool(v.Tombstone); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteTime(v.Created); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (v *Version) Read(r *Reader) (err error) {
	if v.PID, err = r.ReadUint64(); err != nil {
		return err
	}

	if v.Version, err = r.ReadUint64(); err != nil {
		return err
	}

	if v.Region, err = r.ReadString(); err != nil {
		return err
	}

	var hasParent bool
	if hasParent, err = r.ReadBool(); err != nil {
		return err
	}

	if hasParent {
		v.Parent = &Version{}
		if err = v.Parent.Read(r); err != nil {
			return err
		}
	} else {
		v.Parent = nil
	}

	if v.Tombstone, err = r.ReadBool(); err != nil {
		return err
	}

	if v.Created, err = r.ReadTime(); err != nil {
		return err
	}

	return
}

type SchemaVersion struct {
	Name  string
	Major uint32
	Minor uint32
	Patch uint32
}

func (v SchemaVersion) Write(w *Writer) (n int, err error) {
	var m int
	if m, err = w.WriteString(v.Name); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteUint32(v.Major); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteUint32(v.Minor); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteUint32(v.Patch); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (v *SchemaVersion) Read(r *Reader) (err error) {
	if v.Name, err = r.ReadString(); err != nil {
		return err
	}

	if v.Major, err = r.ReadUint32(); err != nil {
		return err
	}

	if v.Minor, err = r.ReadUint32(); err != nil {
		return err
	}

	if v.Patch, err = r.ReadUint32(); err != nil {
		return err
	}

	return
}

type AccessControl struct {
	ClientID    ulid.ULID
	Permissions uint8
}

func (c AccessControl) Write(w *Writer) (n int, err error) {
	var m int
	if m, err = w.WriteULID(c.ClientID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteUint8(c.Permissions); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (c *AccessControl) Read(r *Reader) (err error) {
	if c.ClientID, err = r.ReadULID(); err != nil {
		return err
	}

	if c.Permissions, err = r.ReadUint8(); err != nil {
		return err
	}

	return
}

type Publisher struct {
	PublisherID ulid.ULID
	ClientID    ulid.ULID
	IPAddress   net.IP
	UserAgent   string
}

func (p Publisher) Write(w *Writer) (n int, err error) {
	var m int
	if m, err = w.WriteULID(p.PublisherID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteULID(p.ClientID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteBytes(p.IPAddress); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteString(p.UserAgent); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (p *Publisher) Read(r *Reader) (err error) {
	if p.PublisherID, err = r.ReadULID(); err != nil {
		return err
	}

	if p.ClientID, err = r.ReadULID(); err != nil {
		return err
	}

	var ip []byte
	if ip, err = r.ReadBytes(); err != nil {
		return err
	}
	p.IPAddress = net.IP(ip)

	if p.UserAgent, err = r.ReadString(); err != nil {
		return err
	}

	return
}

type Encryption struct {
	PublicKeyID         string
	EncryptionKey       []byte
	HMACSecret          []byte
	Signature           []byte
	SealingAlgorithm    EncryptionAlgorithm
	EncryptionAlgorithm EncryptionAlgorithm
	SignatureAlgorithm  EncryptionAlgorithm
}

func (e Encryption) Write(w *Writer) (n int, err error) {
	var m int
	if m, err = w.WriteString(e.PublicKeyID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteBytes(e.EncryptionKey); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteBytes(e.HMACSecret); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteBytes(e.Signature); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteUint8(uint8(e.SealingAlgorithm)); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteUint8(uint8(e.EncryptionAlgorithm)); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteUint8(uint8(e.SignatureAlgorithm)); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (e *Encryption) Read(r *Reader) (err error) {
	if e.PublicKeyID, err = r.ReadString(); err != nil {
		return err
	}

	if e.EncryptionKey, err = r.ReadBytes(); err != nil {
		return err
	}

	if e.HMACSecret, err = r.ReadBytes(); err != nil {
		return err
	}

	if e.Signature, err = r.ReadBytes(); err != nil {
		return err
	}

	var a uint8
	if a, err = r.ReadUint8(); err != nil {
		return err
	}
	e.SealingAlgorithm = EncryptionAlgorithm(a)

	if a, err = r.ReadUint8(); err != nil {
		return err
	}
	e.EncryptionAlgorithm = EncryptionAlgorithm(a)

	if a, err = r.ReadUint8(); err != nil {
		return err
	}
	e.SignatureAlgorithm = EncryptionAlgorithm(a)

	return
}

type Compression struct {
	Algorithm CompressionAlgorithm
	Level     int64
}

func (c Compression) Write(w *Writer) (n int, err error) {
	var m int
	if m, err = w.WriteUint8(uint8(c.Algorithm)); err != nil {
		return n + m, err
	}
	n += m

	if m, err = w.WriteInt64(c.Level); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (c *Compression) Read(r *Reader) (err error) {
	var a uint8
	if a, err = r.ReadUint8(); err != nil {
		return err
	}
	c.Algorithm = CompressionAlgorithm(a)

	if c.Level, err = r.ReadInt64(); err != nil {
		return err
	}

	return nil
}

type EncryptionAlgorithm uint8

const (
	Plaintext EncryptionAlgorithm = iota
	AES256_GCM
	AES192_GCM
	AES128_GCM
	HMAC_SHA256
	RSA_OEAP_SHA512
)

type CompressionAlgorithm uint8

const (
	None CompressionAlgorithm = iota
	GZIP
	COMPRESS
	DEFLATE
	BROTLI
)
