package ljpack

import (
	"reflect"

	"github.com/fffonion/ljpack/ljpcode"
)

var stringSliceType = reflect.TypeOf(([]string)(nil))

func encodeStringValue(e *Encoder, v reflect.Value) error {
	return e.EncodeString(v.String())
}

func encodeByteSliceValue(e *Encoder, v reflect.Value) error {
	return e.EncodeBytes(v.Bytes())
}

func encodeByteArrayValue(e *Encoder, v reflect.Value) error {
	if err := e.EncodeBytesLen(v.Len()); err != nil {
		return err
	}

	if v.CanAddr() {
		b := v.Slice(0, v.Len()).Bytes()
		return e.write(b)
	}

	e.buf = grow(e.buf, v.Len())
	reflect.Copy(reflect.ValueOf(e.buf), v)
	return e.write(e.buf)
}

func grow(b []byte, n int) []byte {
	if cap(b) >= n {
		return b[:n]
	}
	b = b[:cap(b)]
	b = append(b, make([]byte, n-len(b))...)
	return b
}

func (e *Encoder) EncodeBytesLen(l int) error {
	// TODO: overflow?
	return e.u124(uint32(l) + 0x20)
}

func (e *Encoder) encodeStringLen(l int) error {
	// TODO: overflow?
	return e.u124(uint32(l) + 0x20)
}

func (e *Encoder) EncodeString(v string) error {
	if intern := e.flags&useInternedStringsFlag != 0; intern || len(e.dict) > 0 {
		return e.encodeInternedString(v, intern)
	}
	return e.encodeNormalString(v)
}

func (e *Encoder) encodeNormalString(v string) error {
	if err := e.encodeStringLen(len(v)); err != nil {
		return err
	}
	return e.writeString(v)
}

func (e *Encoder) EncodeBytes(v []byte) error {
	if v == nil {
		return e.EncodeNull()
	}

	if err := e.EncodeBytesLen(len(v)); err != nil {
		return err
	}
	return e.write(v)
}

func (e *Encoder) EncodeArrayLen(l int) error {
	// TODO: overflow?
	err := e.writeCode(ljpcode.OneBasedArray)
	if err != nil {
		return err
	}
	return e.u124(uint32(l) + 1)
}

func encodeStringSliceValue(e *Encoder, v reflect.Value) error {
	ss := v.Convert(stringSliceType).Interface().([]string)
	return e.encodeStringSlice(ss)
}

func (e *Encoder) encodeStringSlice(s []string) error {
	if s == nil {
		return e.EncodeNull()
	}
	if err := e.EncodeArrayLen(len(s)); err != nil {
		return err
	}
	for _, v := range s {
		if err := e.encodeNormalString(v); err != nil {
			return err
		}
	}
	return nil
}

func encodeSliceValue(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.EncodeNull()
	}
	return encodeArrayValue(e, v)
}

func encodeArrayValue(e *Encoder, v reflect.Value) error {
	l := v.Len()

	if l == 0 {
		return e.EncodeEmpty()
	}

	if err := e.EncodeArrayLen(l); err != nil {
		return err
	}
	for i := 0; i < l; i++ {
		if err := e.EncodeValue(v.Index(i)); err != nil {
			return err
		}
	}
	return nil
}
