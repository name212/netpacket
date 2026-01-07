// Copyright 2026
// license that can be found in the LICENSE file.

package netpacket

import (
	"errors"
	"fmt"
)

var (
	EmptyPayloadErr      = errors.New("empty payload")
	NotImplementedErr    = errors.New("not implemented yet")
	ShortDataErr         = errors.New("data too short")
	CannotParseHeaderErr = errors.New("cannot parse header")
)

func WrapShortDataErr(err error) error {
	return fmt.Errorf("%w for %w", ShortDataErr, err)
}

func WrapCannotParseHeaderErr(err error) error {
	return fmt.Errorf("%w: %w", CannotParseHeaderErr, err)
}

func WrapNotImplementedErr(err error) error {
	return fmt.Errorf("%w: %w", NotImplementedErr, err)
}
