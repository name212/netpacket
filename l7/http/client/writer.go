// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"io"
)

type limitWriter struct {
	parent io.Writer
	limit  int64

	written int64
}

func newLimitWriter(w io.Writer, limit []int64) *limitWriter {
	l := int64(0)

	if len(limit) > 0 && limit[0] > 0 {
		l = limit[0]
	}

	return &limitWriter{
		parent:  w,
		limit:   l,
		written: 0,
	}
}

func (w *limitWriter) Write(p []byte) (int, error) {
	l := len(p)
	if l == 0 {
		return 0, nil
	}

	if w.limit == 0 {
		return w.parent.Write(p)
	}

	if w.written >= w.limit {
		return l, nil
	}

	remain := w.limit - w.written

	l64 := int64(len(p))

	if remain >= l64 {
		return w.write(p)
	}

	toWrite := l64 - remain
	return w.write(p[:toWrite])
}

func (w *limitWriter) write(p []byte) (int, error) {
	n, err := w.parent.Write(p)
	w.written += int64(n)
	return n, err
}

func (w *limitWriter) getWritten(n int64) int64 {
	if w.limit == 0 {
		return n
	}

	if w.limit > n {
		return n
	}

	return w.limit
}
