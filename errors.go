package netpacket

import (
	"errors"
	"fmt"
)

var (
	ShortDataErr         = errors.New("data too short")
	CannotParseHeaderErr = errors.New("cannot parse header")
)

func WrapShortDataErr(msg string) error {
	return fmt.Errorf("%w for %s", ShortDataErr, msg)
}

func WrapCannotParseHeaderErr(msg string) error {
	return fmt.Errorf("%w: %s", CannotParseHeaderErr, msg)
}
