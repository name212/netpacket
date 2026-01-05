// Copyright 2026
// license that can be found in the LICENSE file.

package utils

func CheckBit(b uint8, n uint8) bool {
	shiftedByte := b >> n
	return (shiftedByte & 1) == 1
}
