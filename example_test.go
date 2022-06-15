package ljpack_test

import (
	"bytes"
	"fmt"

	"github.com/fffonion/ljpack"
)

func ExampleMarshal() {
	type Item struct {
		Foo string
	}

	b, err := ljpack.Marshal(&Item{Foo: "bar"})
	if err != nil {
		panic(err)
	}

	var item Item
	err = ljpack.Unmarshal(b, &item)
	if err != nil {
		panic(err)
	}
	fmt.Println(item.Foo)
	// Output: bar
}

func ExampleMarshal_mapStringInterface() {
	in := map[string]interface{}{"foo": 1, "hello": "world"}
	b, err := ljpack.Marshal(in)
	if err != nil {
		panic(err)
	}

	var out map[string]interface{}
	err = ljpack.Unmarshal(b, &out)
	if err != nil {
		panic(err)
	}

	fmt.Println("foo =", out["foo"])
	fmt.Println("hello =", out["hello"])

	// Output:
	// foo = 1
	// hello = world
}

func ExampleDecoder_SetMapDecoder() {
	buf := new(bytes.Buffer)

	enc := ljpack.NewEncoder(buf)
	in := map[string]string{"hello": "world"}
	err := enc.Encode(in)
	if err != nil {
		panic(err)
	}

	dec := ljpack.NewDecoder(buf)

	// Causes decoder to produce map[string]string instead of map[string]interface{}.
	dec.SetMapDecoder(func(d *ljpack.Decoder) (interface{}, error) {
		n, err := d.DecodeMapLen()
		if err != nil {
			return nil, err
		}

		m := make(map[string]string, n)
		for i := 0; i < n; i++ {
			mk, err := d.DecodeString()
			if err != nil {
				return nil, err
			}

			mv, err := d.DecodeString()
			if err != nil {
				return nil, err
			}

			m[mk] = mv
		}
		return m, nil
	})

	out, err := dec.DecodeInterface()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", out)
	// Output: map[string]string{"hello":"world"}
}

func ExampleDecoder_Query() {
	b, err := ljpack.Marshal([]map[string]interface{}{
		{"id": 1, "attrs": map[string]interface{}{"phone": 12345}},
		{"id": 2, "attrs": map[string]interface{}{"phone": 54321}},
	})
	if err != nil {
		panic(err)
	}

	dec := ljpack.NewDecoder(bytes.NewBuffer(b))
	values, err := dec.Query("*.attrs.phone")
	if err != nil {
		panic(err)
	}
	fmt.Println("phones are", values)

	dec.Reset(bytes.NewBuffer(b))
	values, err = dec.Query("1.attrs.phone")
	if err != nil {
		panic(err)
	}
	fmt.Println("2nd phone is", values[0])
	// Output: phones are [12345 54321]
	// 2nd phone is 54321
}

func ExampleEncoder_UseArrayEncodedStructs() {
	type Item struct {
		Foo string
		Bar string
	}

	var buf bytes.Buffer
	enc := ljpack.NewEncoder(&buf)
	enc.UseArrayEncodedStructs(true)

	err := enc.Encode(&Item{Foo: "foo", Bar: "bar"})
	if err != nil {
		panic(err)
	}

	dec := ljpack.NewDecoder(&buf)
	v, err := dec.DecodeInterface()
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
	// Output: [foo bar]
}

func ExampleMarshal_asArray() {
	type Item struct {
		_ljpack struct{} `ljpack:",as_array"`
		Foo      string
		Bar      string
	}

	var buf bytes.Buffer
	enc := ljpack.NewEncoder(&buf)
	err := enc.Encode(&Item{Foo: "foo", Bar: "bar"})
	if err != nil {
		panic(err)
	}

	dec := ljpack.NewDecoder(&buf)
	v, err := dec.DecodeInterface()
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
	// Output: [foo bar]
}

func ExampleMarshal_omitEmpty() {
	type Item struct {
		Foo string
		Bar string
	}

	item := &Item{
		Foo: "hello",
	}
	b, err := ljpack.Marshal(item)
	if err != nil {
		panic(err)
	}
	fmt.Printf("item: %q\n", b)

	type ItemOmitEmpty struct {
		_ljpack struct{} `ljpack:",omitempty"`
		Foo      string
		Bar      string
	}

	itemOmitEmpty := &ItemOmitEmpty{
		Foo: "hello",
	}
	b, err = ljpack.Marshal(itemOmitEmpty)
	if err != nil {
		panic(err)
	}
	fmt.Printf("item2: %q\n", b)

	// Output: item: "\x82\xa3Foo\xa5hello\xa3Bar\xa0"
	// item2: "\x81\xa3Foo\xa5hello"
}

func ExampleMarshal_escapedNames() {
	og := map[string]interface{}{
		"something:special": uint(123),
		"hello, world":      "hello!",
	}
	raw, err := ljpack.Marshal(og)
	if err != nil {
		panic(err)
	}

	type Item struct {
		SomethingSpecial uint   `ljpack:"'something:special'"`
		HelloWorld       string `ljpack:"'hello, world'"`
	}
	var item Item
	if err := ljpack.Unmarshal(raw, &item); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", item)
	//output: ljpack_test.Item{SomethingSpecial:0x7b, HelloWorld:"hello!"}
}
