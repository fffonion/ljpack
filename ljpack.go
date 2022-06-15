package ljpack

import "fmt"

type Marshaler interface {
	MarshalLJpack() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalLJpack([]byte) error
}

type CustomEncoder interface {
	EncodeLJpack(*Encoder) error
}

type CustomDecoder interface {
	DecodeLJpack(*Decoder) error
}

//------------------------------------------------------------------------------

type RawMessage []byte

var (
	_ CustomEncoder = (RawMessage)(nil)
	_ CustomDecoder = (*RawMessage)(nil)
)

func (m RawMessage) EncodeLJpack(enc *Encoder) error {
	return enc.write(m)
}

func (m *RawMessage) DecodeLJpack(dec *Decoder) error {
	msg, err := dec.DecodeRaw()
	if err != nil {
		return err
	}
	*m = msg
	return nil
}

//------------------------------------------------------------------------------

type unexpectedCodeError struct {
	code byte
	hint string
}

func (err unexpectedCodeError) Error() string {
	return fmt.Sprintf("ljpack: unexpected code=%x decoding %s", err.code, err.hint)
}
