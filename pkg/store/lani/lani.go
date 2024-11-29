/*
Lani is sky or heavens in Hawaiian and also sounds like line; we've chosen this word to
represent Honu's byte serialization format for data storage; encoding in a serial or
single skyward direction and decoding in a landward direction.
*/
package lani

// Specifies the interface for objects that can be encoded using lani.
type Encodable interface {
	Size() int
	Encode(*Encoder) (int, error)
}

// Specifies the interface for objects that can be decoded using lani.
type Decodable interface {
	Decode(*Decoder) error
}

// Marshal an encodable type into a byte slice for storage or serialization.
func Marshal(v Encodable) (_ []byte, err error) {
	encoder := &Encoder{}
	encoder.Grow(v.Size())
	if _, err = v.Encode(encoder); err != nil {
		return nil, err
	}
	return encoder.Bytes(), nil
}

// Unmarshal a decodable type from a byte slice for deserialization.
func Unmarshal(data []byte, v Decodable) (err error) {
	decoder := NewDecoder(data)
	return v.Decode(decoder)
}
