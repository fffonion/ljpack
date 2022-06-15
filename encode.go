package ljpack

import (
	"bytes"
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/fffonion/ljpack/ljpcode"
)

const (
	sortMapKeysFlag uint32 = 1 << iota
	arrayEncodedStructsFlag
	useCompactIntsFlag
	useCompactFloatsFlag
	useInternedStringsFlag
	omitEmptyFlag
)

type writer interface {
	io.Writer
	WriteByte(byte) error
}

type byteWriter struct {
	io.Writer
}

func newByteWriter(w io.Writer) byteWriter {
	return byteWriter{
		Writer: w,
	}
}

func (bw byteWriter) WriteByte(c byte) error {
	_, err := bw.Write([]byte{c})
	return err
}

//------------------------------------------------------------------------------

var encPool = sync.Pool{
	New: func() interface{} {
		return NewEncoder(nil)
	},
}

func GetEncoder() *Encoder {
	return encPool.Get().(*Encoder)
}

func PutEncoder(enc *Encoder) {
	enc.w = nil
	encPool.Put(enc)
}

// Marshal returns the MessagePack encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	enc := GetEncoder()

	var buf bytes.Buffer
	enc.Reset(&buf)

	err := enc.Encode(v)
	b := buf.Bytes()

	PutEncoder(enc)

	if err != nil {
		return nil, err
	}
	return b, err
}

type Encoder struct {
	w writer

	buf []byte
	//timeBuf []byte

	dict map[string]int

	flags     uint32
	structTag string
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{
		buf: make([]byte, 9),
	}
	e.Reset(w)
	return e
}

// Writer returns the Encoder's writer.
func (e *Encoder) Writer() io.Writer {
	return e.w
}

// Reset discards any buffered data, resets all state, and switches the writer to write to w.
func (e *Encoder) Reset(w io.Writer) {
	e.ResetDict(w, nil)
}

// ResetDict is like Reset, but also resets the dict.
func (e *Encoder) ResetDict(w io.Writer, dict map[string]int) {
	e.resetWriter(w)
	e.flags = 0
	e.structTag = ""
	e.dict = dict
}

func (e *Encoder) WithDict(dict map[string]int, fn func(*Encoder) error) error {
	oldDict := e.dict
	e.dict = dict
	err := fn(e)
	e.dict = oldDict
	return err
}

func (e *Encoder) resetWriter(w io.Writer) {
	if bw, ok := w.(writer); ok {
		e.w = bw
	} else {
		e.w = newByteWriter(w)
	}
}

// SetSortMapKeys causes the Encoder to encode map keys in increasing order.
// Supported map types are:
//   - map[string]string
//   - map[string]interface{}
func (e *Encoder) SetSortMapKeys(on bool) *Encoder {
	if on {
		e.flags |= sortMapKeysFlag
	} else {
		e.flags &= ^sortMapKeysFlag
	}
	return e
}

// SetCustomStructTag causes the Encoder to use a custom struct tag as
// fallback option if there is no ljpack tag.
func (e *Encoder) SetCustomStructTag(tag string) {
	e.structTag = tag
}

// SetOmitEmpty causes the Encoder to omit empty values by default.
func (e *Encoder) SetOmitEmpty(on bool) {
	if on {
		e.flags |= omitEmptyFlag
	} else {
		e.flags &= ^omitEmptyFlag
	}
}

// UseArrayEncodedStructs causes the Encoder to encode Go structs as ljpack arrays.
func (e *Encoder) UseArrayEncodedStructs(on bool) {
	if on {
		e.flags |= arrayEncodedStructsFlag
	} else {
		e.flags &= ^arrayEncodedStructsFlag
	}
}

// UseCompactEncoding causes the Encoder to chose the most compact encoding.
// For example, it allows to encode small Go int64 as ljpack int8 saving 7 bytes.
func (e *Encoder) UseCompactInts(on bool) {
	if on {
		e.flags |= useCompactIntsFlag
	} else {
		e.flags &= ^useCompactIntsFlag
	}
}

// UseCompactFloats causes the Encoder to chose a compact integer encoding
// for floats that can be represented as integers.
func (e *Encoder) UseCompactFloats(on bool) {
	if on {
		e.flags |= useCompactFloatsFlag
	} else {
		e.flags &= ^useCompactFloatsFlag
	}
}

// UseInternedStrings causes the Encoder to intern strings.
func (e *Encoder) UseInternedStrings(on bool) {
	if on {
		e.flags |= useInternedStringsFlag
	} else {
		e.flags &= ^useInternedStringsFlag
	}
}

func (e *Encoder) Encode(v interface{}) error {
	switch v := v.(type) {
	case nil:
		return e.EncodeNull()
	case string:
		return e.EncodeString(v)
	case []byte:
		return e.EncodeBytes(v)
	case int:
		return e.EncodeInt(int64(v))
	case int64:
		return e.encodeFFIInt64Cond(v)
	case uint:
		return e.EncodeFFIUint(uint64(v))
	case uint64:
		return e.encodeFFIUint64Cond(v)
	case bool:
		return e.EncodeBool(v)
	case float32:
		return e.EncodeDouble(float64(v))
	case float64:
		return e.EncodeDouble(v)
	case time.Duration:
		return e.encodeFFIInt64Cond(int64(v))
	case time.Time:
		return e.EncodeFFIInt64(v.Unix())
	}
	return e.EncodeValue(reflect.ValueOf(v))
}

func (e *Encoder) EncodeMulti(v ...interface{}) error {
	for _, vv := range v {
		if err := e.Encode(vv); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) EncodeValue(v reflect.Value) error {
	fn := getEncoder(v.Type())
	return fn(e, v)
}

func (e *Encoder) EncodeEmpty() error {
	return e.writeCode(ljpcode.EmptyTable)
}

func (e *Encoder) EncodeNil() error {
	return e.writeCode(ljpcode.Nil)
}

func (e *Encoder) EncodeNull() error {
	return e.writeCode(ljpcode.Null)
}

func (e *Encoder) EncodeBool(value bool) error {
	if value {
		return e.writeCode(ljpcode.True)
	}
	return e.writeCode(ljpcode.False)
}

func (e *Encoder) EncodeDuration(d time.Duration) error {
	return e.EncodeInt(int64(d))
}

func (e *Encoder) writeCode(c byte) error {
	return e.w.WriteByte(c)
}

func (e *Encoder) write(b []byte) error {
	_, err := e.w.Write(b)
	return err
}

func (e *Encoder) writeString(s string) error {
	_, err := e.w.Write(stringToBytes(s))
	return err
}
