// Copyright 2026
// license that can be found in the LICENSE file.

package tcp

type Header struct{}

// ParseHeader
// ParseHeader save slices from data. You should copy data before parse
// to avoid hold full original data in memory
// Warning! TODO not implemented
func ParseHeader(data []byte) (*Header, error) {
	return &Header{}, nil
}

func (h *Header) HeaderLen() int {
	// todo need implementation
	return 0
}

func (h *Header) String() string {
	// todo need implementation
	return "TCP header: Not implemented yet"
}
