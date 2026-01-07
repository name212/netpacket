// Copyright 2026
// license that can be found in the LICENSE file.

package netpacket

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyPayload      = errors.New("empty payload")
	ErrNotImplemented    = errors.New("not implemented yet")
	ErrShortData         = errors.New("data too short")
	ErrCannotParseHeader = errors.New("cannot parse header")
)

func WrapShortDataErr(err error) error {
	return fmt.Errorf("%w for %w", ErrShortData, err)
}

func WrapCannotParseHeaderErr(err error) error {
	return fmt.Errorf("%w: %w", ErrCannotParseHeader, err)
}

func WrapNotImplementedErr(err error) error {
	return fmt.Errorf("%w: %w", ErrNotImplemented, err)
}
