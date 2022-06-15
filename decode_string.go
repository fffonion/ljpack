package ljpack

import (
	"fmt"
	"reflect"

	"github.com/fffonion/ljpack/ljpcode"
)

func (d *Decoder) bytesLen() (int, error) {
	n, err := d.u124()
	return int(n) - int(ljpcode.String), err
}

func (d *Decoder) DecodeString() (string, error) {
	if intern := d.flags&useInternedStringsFlag != 0; intern || len(d.dict) > 0 {
		return d.decodeInternedString(intern)
	}

	return d.string()
}

func (d *Decoder) string() (string, error) {
	n, err := d.bytesLen()
	if err != nil {
		return "", err
	}
	return d.stringWithLen(n)
}

func (d *Decoder) stringWithLen(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	b, err := d.readN(n)
	return string(b), err
}

func decodeStringValue(d *Decoder, v reflect.Value) error {
	s, err := d.DecodeString()
	if err != nil {
		return err
	}
	v.SetString(s)
	return nil
}

func (d *Decoder) DecodeBytesLen() (int, error) {
	return d.bytesLen()
}

func (d *Decoder) DecodeBytes() ([]byte, error) {
	return d.bytes(nil)
}

func (d *Decoder) bytes(b []byte) ([]byte, error) {
	n, err := d.bytesLen()
	if err != nil {
		return nil, err
	}
	if n == -1 {
		return nil, nil
	}
	return readN(d.r, b, n)
}

func (d *Decoder) decodeStringTemp() (string, error) {
	if intern := d.flags&useInternedStringsFlag != 0; intern || len(d.dict) > 0 {
		return d.decodeInternedString(intern)
	}

	n, err := d.bytesLen()
	if err != nil {
		return "", err
	}
	if n == -1 {
		return "", nil
	}

	b, err := d.readN(n)
	if err != nil {
		return "", err
	}

	return bytesToString(b), nil
}

func (d *Decoder) decodeBytesPtr(ptr *[]byte) error {
	return d.bytesPtr(ptr)
}

func (d *Decoder) bytesPtr(ptr *[]byte) error {
	n, err := d.bytesLen()
	if err != nil {
		return err
	}
	if n == -1 {
		*ptr = nil
		return nil
	}

	*ptr, err = readN(d.r, *ptr, n)
	return err
}

func (d *Decoder) skipBytes() error {
	n, err := d.bytesLen()
	if err != nil {
		return err
	}
	if n <= 0 {
		return nil
	}
	return d.skipN(n)
}

func decodeBytesValue(d *Decoder, v reflect.Value) error {
	b, err := d.bytes(v.Bytes())
	if err != nil {
		return err
	}

	v.SetBytes(b)

	return nil
}

func decodeByteArrayValue(d *Decoder, v reflect.Value) error {
	n, err := d.bytesLen()
	if err != nil {
		return err
	}
	if n == -1 {
		return nil
	}
	if n > v.Len() {
		return fmt.Errorf("%s len is %d, but ljpack has %d elements", v.Type(), v.Len(), n)
	}

	b := v.Slice(0, n).Bytes()
	return d.readFull(b)
}
