// Copyright 2026
// license that can be found in the LICENSE file.

package content

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	LenHeader = "Content-Length"
)

var (
	ErrInvalidLen = fmt.Errorf("invalid len")
	ErrNoLen      = fmt.Errorf("no len")
)

func LenStr(r *http.Response) string {
	if r == nil || len(r.Header) == 0 {
		return ""
	}

	return r.Header.Get(LenHeader)
}

func Len[T int | int64](r *http.Response) T {
	l, _ := parseLen(LenStr(r))
	return T(l)
}

func LenRequired[T int | int64](r *http.Response) (T, error) {
	l, err := parseLen(LenStr(r))
	return T(l), err
}

func Len32(r *http.Response) int {
	return Len[int](r)
}

func Len32Required(r *http.Response) (int, error) {
	return LenRequired[int](r)
}

func Len64(r *http.Response) int64 {
	return Len[int64](r)
}

func Len64Required(r *http.Response) (int64, error) {
	return LenRequired[int64](r)
}

func SetLenHeader[T int | int64](w http.ResponseWriter, l T) {
	if l < 0 {
		l = 0
	}

	w.Header().Add(LenHeader, fmt.Sprintf("%d", l))
}

func AppendLenHeader[T int | int64](dest map[string]string, l T) {
	if l < 0 {
		l = 0
	}
	dest[LenHeader] = fmt.Sprintf("%d", l)
}

func parseLen(lStr string) (int64, error) {
	if lStr == "" {
		return 0, fmt.Errorf("no content len %w", ErrInvalidLen)
	}

	l, err := strconv.ParseInt(lStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse content len %w: %w", ErrInvalidLen, err)
	}

	if l == 0 {
		return 0, ErrNoLen
	}

	if l < 0 {
		return 0, fmt.Errorf("%w negative content len", ErrInvalidLen)
	}

	return l, nil
}
