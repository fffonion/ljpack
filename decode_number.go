package ljpack

import (
	"fmt"
	"math"
	"reflect"

	"github.com/fffonion/ljpack/ljpcode"
)

func (d *Decoder) skipN(n int) error {
	_, err := d.readN(n)
	return err
}

func (d *Decoder) u124() (uint32, error) {
	c1, err := d.readCode()
	if err != nil {
		return 0, err
	}
	if c1 < 0xe0 {
		return uint32(c1), nil
	}
	c2, err := d.readCode()
	if err != nil {
		return 0, err
	}

	if c1 != 0xff {
		return uint32(c1&0x1f)<<8 + uint32(c2) + 0xe0, nil
	}
	b, err := d.readN(3) // read three remaining bytes, c1 is the marker 0xff
	if err != nil {
		return 0, err
	}

	n := (uint32(b[2]) << 24) |
		(uint32(b[1]) << 16) |
		(uint32(b[0]) << 8) |
		uint32(c2)
	return n, nil
}

func (d *Decoder) int32() (int32, error) {
	b, err := d.readN(4)
	// Little-Endian
	n := (int32(b[3]) << 24) |
		(int32(b[2]) << 16) |
		(int32(b[1]) << 8) |
		int32(b[0])
	return int32(n), err
}

func (d *Decoder) ffiUint64() (uint64, error) {
	b, err := d.readN(8)
	if err != nil {
		return 0, err
	}
	// Little-Endian
	n := (uint64(b[7]) << 56) |
		(uint64(b[6]) << 48) |
		(uint64(b[5]) << 40) |
		(uint64(b[4]) << 32) |
		(uint64(b[3]) << 24) |
		(uint64(b[2]) << 16) |
		(uint64(b[1]) << 8) |
		uint64(b[0])
	return n, nil
}

func (d *Decoder) ffiInt64() (int64, error) {
	n, err := d.ffiUint64()
	return int64(n), err
}

func (d *Decoder) ffiComplex() (complex128, error) {
	re, err := d.ffiUint64()
	if err != nil {
		return 0, err
	}
	ref := math.Float64frombits(re)

	im, err := d.ffiUint64()
	if err != nil {
		return 0, err
	}
	imf := math.Float64frombits(im)

	return complex(ref, imf), nil
}

// DecodeFFIUint64 decodes ljpack int8/16/32/64 and uint8/16/32/64
// into Go uint64.
func (d *Decoder) DecodeFFIUint64() (uint64, error) {
	c, err := d.readCode()
	if err != nil {
		return 0, err
	}
	return d.uint(c)
}

func (d *Decoder) uint(c byte) (uint64, error) {
	if c == ljpcode.Nil {
		return 0, nil
	}
	switch c {
	case ljpcode.Int:
		n, err := d.int32()
		return uint64(n), err
	case ljpcode.FFIUint64, ljpcode.FFIInt64:
		return d.ffiUint64()
	}
	return 0, fmt.Errorf("ljpack: invalid code=%x decoding uint64", c)
}

// DecodeFFIInt64 decodes ljpack int8/16/32/64 and uint8/16/32/64
// into Go int64.
func (d *Decoder) DecodeFFIInt64() (int64, error) {
	c, err := d.readCode()
	if err != nil {
		return 0, err
	}
	return d.int(c)
}

func (d *Decoder) int(c byte) (int64, error) {
	if c == ljpcode.Nil {
		return 0, nil
	}
	switch c {
	case ljpcode.Int:
		n, err := d.int32()
		return int64(int32(n)), err
	case ljpcode.FFIUint64, ljpcode.FFIInt64:
		n, err := d.ffiUint64()
		return int64(n), err
	}
	return 0, fmt.Errorf("ljpack: invalid code=%x decoding int64", c)
}

// DecodeDouble decodes ljpack float32/64 into Go float64.
func (d *Decoder) DecodeDouble() (float64, error) {
	c, err := d.readCode()
	if err != nil {
		return 0, err
	}
	return d.double(c)
}

func (d *Decoder) double(c byte) (float64, error) {
	switch c {
	case ljpcode.Double:
		n, err := d.ffiUint64()
		if err != nil {
			return 0, err
		}
		return math.Float64frombits(n), nil
	}

	n, err := d.int(c)
	if err != nil {
		return 0, fmt.Errorf("ljpack: invalid code=%x decoding float64", c)
	}
	return float64(n), nil
}

func (d *Decoder) DecodeUint() (uint, error) {
	n, err := d.DecodeFFIUint64()
	return uint(n), err
}

func (d *Decoder) DecodeUint32() (uint32, error) {
	n, err := d.DecodeFFIUint64()
	return uint32(n), err
}

func (d *Decoder) DecodeInt() (int, error) {
	n, err := d.DecodeFFIInt64()
	return int(n), err
}

func (d *Decoder) DecodeInt32() (int32, error) {
	n, err := d.DecodeFFIInt64()
	return int32(n), err
}

func (d *Decoder) DecodeFFIComplex() (complex128, error) {
	c, err := d.readCode()
	if err != nil {
		return 0, err
	}
	if c != ljpcode.FFIComplex {
		return 0, fmt.Errorf("ljpack: invalid code=%x decoding complex", c)
	}
	return d.ffiComplex()
}

func decodeDoubleValue(d *Decoder, v reflect.Value) error {
	f, err := d.DecodeDouble()
	if err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}

func decodeFFIComplexValue(d *Decoder, v reflect.Value) error {
	f, err := d.DecodeDouble()
	if err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}

func decodeFFIInt64Value(d *Decoder, v reflect.Value) error {
	n, err := d.DecodeFFIInt64()
	if err != nil {
		return err
	}
	v.SetInt(n)
	return nil
}

func decodeFFIUint64Value(d *Decoder, v reflect.Value) error {
	n, err := d.DecodeFFIUint64()
	if err != nil {
		return err
	}
	v.SetUint(n)
	return nil
}
