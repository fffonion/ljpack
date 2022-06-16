package ljpack

import (
	"fmt"
	"math"
	"reflect"

	"github.com/fffonion/ljpack/ljpcode"
)

const (
	minInternedStringLen = 3
	maxDictLen           = math.MaxUint16
)

//------------------------------------------------------------------------------

func encodeInternedInterfaceValue(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.EncodeNil()
	}

	v = v.Elem()
	if v.Kind() == reflect.String {
		return e.encodeInternedString(v.String(), true)
	}
	return e.EncodeValue(v)
}

func encodeInternedStringValue(e *Encoder, v reflect.Value) error {
	return e.encodeInternedString(v.String(), true)
}

func (e *Encoder) encodeInternedString(s string, intern bool) error {
	// Interned string takes at least 3 bytes. Plain string 1 byte + string len.
	if len(s) >= minInternedStringLen {
		if idx, ok := e.dict[s]; ok {
			return e.encodeInternedStringIndex(idx)
		}

		if intern && len(e.dict) < maxDictLen {
			if e.dict == nil {
				e.dict = make(map[string]int)
			}
			idx := len(e.dict)
			e.dict[s] = idx
		}
	}

	return e.encodeNormalString(s)
}

func (e *Encoder) encodeInternedStringIndex(idx int) error {
	e.buf = e.buf[:1]
	e.buf[0] = ljpcode.StringInterned
	err := e.write(e.buf)
	if err != nil {
		return err
	}

	if uint64(idx) <= math.MaxUint32 {
		return e.u124(uint32(idx))
	}

	return fmt.Errorf("ljpack: interned string index=%d is too large", idx)
}

//------------------------------------------------------------------------------

func decodeInternedInterfaceValue(d *Decoder, v reflect.Value) error {
	s, err := d.decodeInternedString(true)
	if err == nil {
		v.Set(reflect.ValueOf(s))
		return nil
	}
	if err != nil {
		if _, ok := err.(unexpectedCodeError); !ok {
			return err
		}
	}

	// if err := d.s.UnreadByte(); err != nil {
	// 	return err
	// }
	return decodeInterfaceValue(d, v)
}

func decodeInternedStringValue(d *Decoder, v reflect.Value) error {
	s, err := d.decodeInternedString(true)
	if err != nil {
		return err
	}

	v.SetString(s)
	return nil
}

func (d *Decoder) decodeInternedString(intern bool) (string, error) {
	c, err := d.readCode()
	if err != nil {
		return "", err
	}

	if c == ljpcode.StringInterned {
		n, err := d.u124()
		if err != nil {
			return "", err
		}

		return d.internedStringAtIndex(int(n))
	}

	if ljpcode.IsString(c) {
		err = d.s.UnreadByte()
		if err != nil {
			return "", err
		}
		return d.string()
	}

	return "", unexpectedCodeError{
		code: c,
		hint: "interned string",
	}

}

func (d *Decoder) internedStringAtIndex(idx int) (string, error) {
	if idx >= len(d.dict) {
		err := fmt.Errorf("ljpack: interned string at index=%d does not exist", idx)
		return "", err
	}
	return d.dict[idx], nil
}

func (d *Decoder) decodeInternedStringWithLen(n int, intern bool) (string, error) {
	if n <= 0 {
		return "", nil
	}

	s, err := d.stringWithLen(n)
	if err != nil {
		return "", err
	}

	if intern && len(s) >= minInternedStringLen && len(d.dict) < maxDictLen {
		d.dict = append(d.dict, s)
	}

	return s, nil
}
