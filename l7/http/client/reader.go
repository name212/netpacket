// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"bytes"
	"io"
)


type readerClose struct {
	io.Reader
}

func newReaderCloseBytes(b []byte) *readerClose {
	return &readerClose{
		Reader: bytes.NewReader(b),
	}
}

func newReaderClose(b io.Reader) *readerClose {
	return &readerClose{
		Reader: b,
	}
}

func (r *readerClose) Close() error {
	var err error
	if closer, ok := r.Reader.(io.Closer); ok {
		err = closer.Close()
	}
	r.Reader = nil
	return err
}