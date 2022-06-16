# LuaJIT string.buffer encoding for Golang

[![Build Status](https://travis-ci.org/fffonion/ljpack.svg)](https://travis-ci.org/fffonion/ljpack)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/fffonion/ljpack)](https://pkg.go.dev/github.com/fffonion/ljpack)


## Resources

- [Reference](https://pkg.go.dev/github.com/fffonion/ljpack)
- [Examples](https://pkg.go.dev/github.com/fffonion/ljpack#pkg-examples)

## Notes

The encoding format for LuaJIT string.buffer is not a formalized structure, the format
could change at any time. Thus use this project with your own risk.

Supported:

- nil, false, true, userdata NULL
- int (int32), double (float64)
- Empty table, hash, 0-based array, 1-based array
- FFI int64, uint64, complex
- string, interned string

Work in Progress:

- lightud32, lightud64
- Mixed table, Metatable dict entry

Not supported

- Non-string value as hash keys

## Features

- Primitives, arrays, maps, structs, time.Time and interface{}.
- Appengine \*datastore.Key and datastore.Cursor.
- [CustomEncoder]/[CustomDecoder] interfaces for custom encoding.
- Renaming fields via `ljpack:"my_field_name"` and alias via `ljpack:"alias:another_name"`.
- Omitting individual empty fields via `ljpack:",omitempty"` tag or all
  [empty fields in a struct](https://pkg.go.dev/github.com/fffonion/ljpack#example-Marshal-OmitEmpty).
- [Map keys sorting](https://pkg.go.dev/github.com/fffonion/ljpack#Encoder.SetSortMapKeys).
- Encoding/decoding all
  [structs as arrays](https://pkg.go.dev/github.com/fffonion/ljpack#Encoder.UseArrayEncodedStructs)
  or
  [individual structs](https://pkg.go.dev/github.com/fffonion/ljpack#example-Marshal-AsArray).
- [Encoder.SetCustomStructTag] with [Decoder.SetCustomStructTag] can turn ljpack into drop-in
  replacement for any tag.
- Simple but very fast and efficient
  [queries](https://pkg.go.dev/github.com/fffonion/ljpack#example-Decoder.Query).

[customencoder]: https://pkg.go.dev/github.com/fffonion/ljpack#CustomEncoder
[customdecoder]: https://pkg.go.dev/github.com/fffonion/ljpack#CustomDecoder
[encoder.setcustomstructtag]:
  https://pkg.go.dev/github.com/fffonion/ljpack#Encoder.SetCustomStructTag
[decoder.setcustomstructtag]:
  https://pkg.go.dev/github.com/fffonion/ljpack#Decoder.SetCustomStructTag

## Installation

ljpack supports 2 last Go versions and requires support for
[Go modules](https://github.com/golang/go/wiki/Modules). So make sure to initialize a Go module:

```shell
go mod init github.com/my/repo
```

And then install ljpack:

```shell
go get github.com/fffonion/ljpack
```

## Quickstart

```go
import "github.com/fffonion/ljpack"

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
```

## Credits

- Forked from https://github.com/vmihailenco/msgpack

