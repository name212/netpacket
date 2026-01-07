// Copyright 2026
// license that can be found in the LICENSE file.

package netpacket

import "fmt"

type (
	PayloadExtractor func([]byte) ([]byte, error)
	Kind             string
)

type Kinder interface {
	Kind() Kind
}

type Header interface {
	fmt.Stringer
	HeaderLen() int
}

type Packet[T Header] interface {
	fmt.Stringer
	Kinder

	GetHeader() *T

	GetHeaderData() []byte
	GetPayload() []byte
}
