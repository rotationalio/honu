package pagination

import (
	"encoding/base64"
	"errors"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	DefaultPageSize int32 = 100
	MaximumPageSize int32 = 5000
	CursorDuration        = 24 * time.Hour
)

var (
	ErrMissingExpiration  = errors.New("cursor does not have an expires timestamp")
	ErrCursorExpired      = errors.New("cursor has expired and is no longer useable")
	ErrUnparsableToken    = errors.New("could not parse the next page token")
	ErrTokenQueryMismatch = errors.New("cannot change query parameters during pagination")
)

func New(nextKey []byte, namespace string, pageSize int32) *PageCursor {
	if pageSize == 0 {
		pageSize = DefaultPageSize
	}

	return &PageCursor{
		PageSize:  pageSize,
		NextKey:   nextKey,
		Namespace: namespace,
		Expires:   timestamppb.New(time.Now().Add(CursorDuration)),
	}
}

func Parse(token string) (cursor *PageCursor, err error) {
	var data []byte
	if data, err = base64.RawURLEncoding.DecodeString(token); err != nil {
		return nil, ErrUnparsableToken
	}

	cursor = &PageCursor{}
	if err = proto.Unmarshal(data, cursor); err != nil {
		return nil, ErrUnparsableToken
	}

	var expired bool
	if expired, err = cursor.HasExpired(); err != nil {
		return nil, err
	}

	if expired {
		return nil, ErrCursorExpired
	}

	return cursor, nil
}

func (c *PageCursor) NextPageToken() (_ string, err error) {
	var expired bool
	if expired, err = c.HasExpired(); err != nil {
		return "", err
	}

	if expired {
		return "", ErrCursorExpired
	}

	var data []byte
	if data, err = proto.Marshal(c); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func (c *PageCursor) HasExpired() (bool, error) {
	if c.Expires == nil {
		return false, ErrMissingExpiration
	}
	return time.Now().After(c.Expires.AsTime()), nil
}

func (c *PageCursor) IsZero() bool {
	return c.PageSize == 0 && len(c.NextKey) == 0 && c.Namespace == "" && (c.Expires == nil || c.Expires.AsTime().IsZero())
}
