package metadata

import (
	"encoding/binary"
	"net"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/honu/pkg/store/lani"
)

//===========================================================================
// Provenance
//===========================================================================

type Publisher struct {
	PublisherID ulid.ULID
	ClientID    ulid.ULID
	IPAddress   net.IP
	UserAgent   string
}

var _ lani.Encodable = &Publisher{}
var _ lani.Decodable = &Publisher{}

func (o *Publisher) Size() (s int) {
	// 2 ULIDs and 2 variable byte arrays
	s += 16 + 16
	s += len(o.IPAddress) + binary.MaxVarintLen64
	s += len([]byte(o.UserAgent)) + binary.MaxVarintLen64
	return
}

func (o *Publisher) Encode(e *lani.Encoder) (n int, err error) {
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

func (o *Publisher) Decode(d *lani.Decoder) (err error) {
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
