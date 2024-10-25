package store

import "encoding/binary"

// Used for testing this package, should not be used in external tests.
type Mock struct {
	Data []byte
	Err  error
}

func (t *Mock) Size() int {
	return len(t.Data) + binary.MaxVarintLen64
}

func (t *Mock) Encode(e *Encoder) (int, error) {
	if t.Err != nil {
		return 0, t.Err
	}
	return e.Encode(t.Data)
}

func (t *Mock) Decode(d *Decoder) (err error) {
	if t.Err != nil {
		return t.Err
	}

	if t.Data, err = d.Decode(); err != nil {
		return err
	}

	return nil
}
