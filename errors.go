package netpacket

import "errors"

var (
	ShortDataErr = errors.New("data too short")
)
