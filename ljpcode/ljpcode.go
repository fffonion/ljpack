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

	// Uint8  byte = 0xcc
	// Uint16 byte = 0xcd
	// Uint32 byte = 0xce
	// Uint64 byte = 0xcf

	//Int8  byte = 0xd0
	//Int16 byte = 0xd1
	//Int32 byte = 0xd2
	FFIInt64   byte = 0x10
	FFIUint64  byte = 0x11
	FFIComplex byte = 0x12

	String    byte = 0x20
	StringKey byte = 0x0f

	// FixedStrLow  byte = 0xa0
	// FixedStrHigh byte = 0xbf
	// FixedStrMask byte = 0x1f
	// Str8         byte = 0xd9
	// Str16        byte = 0xda
	// Str32        byte = 0xdb

	// Bin8  byte = 0xc4
	// Bin16 byte = 0xc5
	// Bin32 byte = 0xc6

	// FixedArrayLow  byte = 0x90
	// FixedArrayHigh byte = 0x9f
	// FixedArrayMask byte = 0xf
	// Array16        byte = 0xdc
	// Array32        byte = 0xdd

	// FixedMapLow  byte = 0x80
	// FixedMapHigh byte = 0x8f
	// FixedMapMask byte = 0xf
	// Map16        byte = 0xde
	// Map32        byte = 0xdf

	// FixExt1  byte = 0xd4
	// FixExt2  byte = 0xd5
	// FixExt4  byte = 0xd6
	// FixExt8  byte = 0xd7
	// FixExt16 byte = 0xd8
	// Ext8     byte = 0xc7
	// Ext16    byte = 0xc8
	// Ext32    byte = 0xc9
)

func IsString(c byte) bool {
	return c >= String // zero length string also ok
}

func IsMap(c byte) bool {
	return c == Hash || c == MixedOneBasedArray || c == MixedZeroBasedArray
}

func IsArray(c byte) bool {
	return c == OneBasedArray || c == ZeroBasedArray || c == MixedOneBasedArray || c == MixedZeroBasedArray
}
