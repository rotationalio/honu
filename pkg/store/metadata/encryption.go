package metadata

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.rtnl.ai/honu/pkg/store/lani"
)

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

var _ lani.Encodable = (*Encryption)(nil)
var _ lani.Decodable = (*Encryption)(nil)

// The static size of a zero valued Encryption object; see TestEncryptionSize for details.
const encryptionStaticSize = 43

func (o *Encryption) Size() (s int) {
	s = encryptionStaticSize
	s += len([]byte(o.PublicKeyID))
	s += len(o.EncryptionKey)
	s += len(o.HMACSecret)
	s += len(o.Signature)
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

func (ea EncryptionAlgorithm) String() string {
	switch ea {
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

func (ea *EncryptionAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(ea.String())
}

func (ea *EncryptionAlgorithm) UnmarshalJSON(data []byte) (err error) {
	var alg string
	if err := json.Unmarshal(data, &alg); err != nil {
		return err
	}
	if *ea, err = ParseEncryptionAlgorithm(alg); err != nil {
		return err
	}
	return nil
}

func (ea EncryptionAlgorithm) Value() uint8 {
	return uint8(ea)
}
