// +build !appengine

package ljpack

import (
	"unsafe"
)

// bytesToString converts byte slice to string.
func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// stringToBytes converts string to byte slice.
func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
