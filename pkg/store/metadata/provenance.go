package metadata

import (
	"net"

	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/ulid"
)

//===========================================================================
// Provenance
//===========================================================================

type Publisher struct {
	PublisherID ulid.ULID `json:"publisher_id" msg:"publisher_id"`
	ClientID    ulid.ULID `json:"client_id" msg:"client_id"`
	IPAddress   net.IP    `json:"ipaddr" msg:"ipaddr"`
	UserAgent   string    `json:"user_agent,omitempty" msg:"user_agent,omitempty"`
}

var _ lani.Encodable = &Publisher{}
var _ lani.Decodable = &Publisher{}

// The static size of a zero valued Publisher object; see TestPublisherSize for details.
const publisherStaticSize = 52

func (o *Publisher) Size() (s int) {
	s = publisherStaticSize
	s += len(o.IPAddress)
	s += len([]byte(o.UserAgent))
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
