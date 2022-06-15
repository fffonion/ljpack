package ljpack

import (
	"math"
	"reflect"

	"github.com/fffonion/ljpack/ljpcode"
)

func (e *Encoder) u124(n uint32) error {
	if n < 0xe0 {
		e.buf = e.buf[:1]
		e.buf[0] = byte(n & 0xFF)
		return e.write(e.buf)
	}

	if n < 0x1fe0 {
		n -= 0xe0
		e.buf = e.buf[:2]
		e.buf[0] = byte(0xe0 | (n >> 8))
		e.buf[1] = byte(n & 0xFF)
		return e.write(e.buf)
	}

	return e.write4(0xFF, n)
}

// EncodeFFIUint64 encodes an uint16 in 9 bytes preserving type of the number.
func (e *Encoder) EncodeFFIUint64(n uint64) error {
	return e.write8(ljpcode.FFIUint64, n)
}

func (e *Encoder) encodeFFIUint64Cond(n uint64) error {
	if e.flags&useCompactIntsFlag != 0 {
		return e.EncodeFFIUint(n)
	}
	return e.EncodeFFIUint64(n)
}

// EncodeInt32 encodes an int32 in 5 bytes preserving type of the number.
func (e *Encoder) EncodeInt32(n int32) error {
	return e.write4(ljpcode.Int, uint32(n))
}

func (e *Encoder) encodeInt32Cond(n int32) error {
	return e.EncodeInt32(n)
}

// EncodeFFIInt64 encodes an int64 in 9 bytes preserving type of the number.
func (e *Encoder) EncodeFFIInt64(n int64) error {
	return e.write8(ljpcode.FFIInt64, uint64(n))
}

func (e *Encoder) encodeFFIInt64Cond(n int64) error {
	if e.flags&useCompactIntsFlag != 0 {
		return e.EncodeInt(n)
	}
	return e.EncodeFFIInt64(n)
}

// EncodeUnsignedNumber encodes an uint64 in 1, 2, 3, 5, or 9 bytes.
// Type of the number is lost during encoding.
func (e *Encoder) EncodeFFIUint(n uint64) error {
	return e.EncodeFFIUint64(n)
}

// EncodeNumber encodes an int64 in 1, 2, 3, 5, or 9 bytes.
// Type of the number is lost during encoding.
func (e *Encoder) EncodeInt(n int64) error {
	if n >= 0 {
		return e.EncodeFFIUint(uint64(n))
	}
	if n >= math.MinInt32 {
		return e.EncodeInt32(int32(n))
	}
	return e.EncodeFFIInt64(n)
}

func (e *Encoder) EncodeDouble(n float64) error {
	if e.flags&useCompactFloatsFlag != 0 {
		// Both NaN and Inf convert to int64(-0x8000000000000000)
		// If n is NaN then it never compares true with any other value
		// If n is Inf then it doesn't convert from int64 back to +/-Inf
		// In both cases the comparison works.
		if float64(int64(n)) == n {
			return e.EncodeInt(int64(n))
		}
	}
	return e.write8(ljpcode.Double, math.Float64bits(n))
}

func (e *Encoder) encodeFFIComplex(c complex128) error {
	re := real(c)
	err := e.write8(ljpcode.FFIComplex, math.Float64bits(re))
	if err != nil {
		return err
	}
	im := imag(c)
	n := math.Float64bits(im)
	e.buf = e.buf[:8]
	e.buf[7] = byte(n >> 56)
	e.buf[6] = byte(n >> 48)
	e.buf[5] = byte(n >> 40)
	e.buf[4] = byte(n >> 32)
	e.buf[3] = byte(n >> 24)
	e.buf[2] = byte(n >> 16)
	e.buf[1] = byte(n >> 8)
	e.buf[0] = byte(n)
	return e.write(e.buf)
}

func (e *Encoder) write1(code byte, n uint8) error {
	e.buf = e.buf[:2]
	e.buf[0] = code
	e.buf[1] = n
	return e.write(e.buf)
}

func (e *Encoder) write2(code byte, n uint16) error {
	e.buf = e.buf[:3]
	e.buf[0] = code
	e.buf[2] = byte(n >> 8)
	e.buf[1] = byte(n)
	return e.write(e.buf)
}

func (e *Encoder) write4(code byte, n uint32) error {
	e.buf = e.buf[:5]
	e.buf[0] = code
	e.buf[4] = byte(n >> 24)
	e.buf[3] = byte(n >> 16)
	e.buf[2] = byte(n >> 8)
	e.buf[1] = byte(n)
	return e.write(e.buf)
}

func (e *Encoder) write8(code byte, n uint64) error {
	e.buf = e.buf[:9]
	e.buf[0] = code
	e.buf[8] = byte(n >> 56)
	e.buf[7] = byte(n >> 48)
	e.buf[6] = byte(n >> 40)
	e.buf[5] = byte(n >> 32)
	e.buf[4] = byte(n >> 24)
	e.buf[3] = byte(n >> 16)
	e.buf[2] = byte(n >> 8)
	e.buf[1] = byte(n)
	return e.write(e.buf)
}

func encodeFFIUintValue(e *Encoder, v reflect.Value) error {
	return e.EncodeFFIUint(v.Uint())
}

func encodeIntValue(e *Encoder, v reflect.Value) error {
	return e.EncodeInt(v.Int())
}

func encodeFFIUint64CondValue(e *Encoder, v reflect.Value) error {
	return e.encodeFFIUint64Cond(v.Uint())
}

func encodeInt8CondValue(e *Encoder, v reflect.Value) error {
	return e.encodeInt32Cond(int32(v.Int()))
}

func encodeInt16CondValue(e *Encoder, v reflect.Value) error {
	return e.encodeInt32Cond(int32(v.Int()))
}

func encodeInt32CondValue(e *Encoder, v reflect.Value) error {
	return e.encodeInt32Cond(int32(v.Int()))
}

func encodeFFIInt64CondValue(e *Encoder, v reflect.Value) error {
	return e.encodeFFIInt64Cond(v.Int())
}

func encodeFFIComplex(e *Encoder, v reflect.Value) error {
	return e.encodeFFIComplex(v.Complex())
}

func encodeDoubleValue(e *Encoder, v reflect.Value) error {
	return e.EncodeDouble(v.Float())
}
