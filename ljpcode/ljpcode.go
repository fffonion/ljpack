package ljpcode

var (
	PosFixedNumHigh byte = 0x7f
	NegFixedNumLow  byte = 0xe0

	Nil byte = 0x00

	False byte = 0x01
	True  byte = 0x02

	Null            byte = 0x03
	LightUserData32 byte = 0x04
	LightUserData64 byte = 0x05

	Int    byte = 0x06
	Double byte = 0x07
	// Num byte = 0x07

	EmptyTable          byte = 0x08
	Hash                byte = 0x09
	ZeroBasedArray      byte = 0x0a
	MixedZeroBasedArray byte = 0x0b
	OneBasedArray       byte = 0x0c
	MixedOneBasedArray  byte = 0x0d
	Metatable           byte = 0x0e

	FFIInt64   byte = 0x10
	FFIUint64  byte = 0x11
	FFIComplex byte = 0x12

	String         byte = 0x20
	StringInterned byte = 0x0f
)

func IsString(c byte) bool {
	return c >= String || // zero length string also ok
		c == StringInterned
}

func IsMap(c byte) bool {
	return c == Hash || c == MixedOneBasedArray || c == MixedZeroBasedArray
}

func IsArray(c byte) bool {
	return c == OneBasedArray || c == ZeroBasedArray || c == MixedOneBasedArray || c == MixedZeroBasedArray
}
