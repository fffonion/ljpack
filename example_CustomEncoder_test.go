package ljpack_test

import (
	"fmt"

	"github.com/fffonion/ljpack"
)

type customStruct struct {
	S string
	N int
}

var _ ljpack.CustomEncoder = (*customStruct)(nil)
var _ ljpack.CustomDecoder = (*customStruct)(nil)

func (s *customStruct) EncodeMsgpack(enc *ljpack.Encoder) error {
	return enc.EncodeMulti(s.S, s.N)
}

func (s *customStruct) DecodeMsgpack(dec *ljpack.Decoder) error {
	return dec.DecodeMulti(&s.S, &s.N)
}

func ExampleCustomEncoder() {
	b, err := ljpack.Marshal(&customStruct{S: "hello", N: 42})
	if err != nil {
		panic(err)
	}

	var v customStruct
	err = ljpack.Unmarshal(b, &v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", v)
	// Output: ljpack_test.customStruct{S:"hello", N:42}
}
