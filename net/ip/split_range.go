// Copyright 2026
// license that can be found in the LICENSE file.

package ip

import (
	"encoding/binary"
	"fmt"
	"iter"
	"math"
	"math/big"
	math_bits "math/bits"
	"net/netip"
)

// SplitRangeToSubnets
// returns all subnets (summarize) from IP addresses range
// based on python ipaddress.summarize_address_range function
// https://github.com/python/cpython/blob/8a5e24a6049338b928f52af79fc610e5b9b2f3fa/Lib/ipaddress.py#L197
// IP's should has same type IPv4 or IPv6 both
func SplitRangeToSubnets(first netip.Addr, second netip.Addr) (iter.Seq[netip.Prefix], error) {
	bitLen := first.BitLen()

	if bitLen != second.BitLen() {
		return nil, fmt.Errorf("addresses should same type, first bit len %d != %d second bit len", bitLen, second.BitLen())
	}

	if first.Compare(second) > 0 {
		return nil, fmt.Errorf("last IP address must be greater than first")
	}

	if first.Is4() {
		return newSplitRangeCollectorV4(&first, &second).Iter(), nil
	}

	return newSplitRangeCollectorV6(&first, &second).Iter(), nil
}

type (
	ipV4 = [4]byte
	ipV6 = [16]byte
)

func minInt(f, s int) int {
	if f <= s {
		return f
	}

	return s
}

func newInt() *big.Int {
	return new(big.Int)
}

func copyInt(i *big.Int) *big.Int {
	res := new(big.Int)
	return res.SetBytes(i.Bytes())
}

var (
	allOnesV4 uint32 = math.MaxUint32
	allOnesV6        = newInt().Sub(
		newInt().Exp(big.NewInt(2), big.NewInt(128), nil),
		big.NewInt(1),
	)
	bIntZero = big.NewInt(0)
	bIntOne  = big.NewInt(1)
)

type splitRangeCollectorV4 struct {
	first  uint32
	second uint32

	ipBits        int
	tempAddrArray ipV4
}

func newSplitRangeCollectorV4(first *netip.Addr, second *netip.Addr) *splitRangeCollectorV4 {
	fs := first.As4()
	ss := second.As4()

	arr := ipV4{}

	return &splitRangeCollectorV4{
		first:         binary.BigEndian.Uint32(fs[:]),
		second:        binary.BigEndian.Uint32(ss[:]),
		ipBits:        first.BitLen(),
		tempAddrArray: arr,
	}
}

func (c *splitRangeCollectorV4) Iter() iter.Seq[netip.Prefix] {
	return func(yield func(netip.Prefix) bool) {
		for c.first <= c.second {
			nbits := minInt(
				c.countRighthandZeroBits(c.first, c.ipBits),
				c.bitsLen(c.second-c.first+1)-1,
			)
			if !yield(
				netip.PrefixFrom(
					netip.AddrFrom4(c.addr(c.first)),
					c.ipBits-nbits,
				),
			) {
				return
			}

			c.first += 1 << nbits

			if c.first+1 == allOnesV4 {
				return
			}
		}
	}
}

func (c *splitRangeCollectorV4) addr(v uint32) ipV4 {
	c.tempAddrArray[0] = byte(v >> 24)
	c.tempAddrArray[1] = byte(v >> 16)
	c.tempAddrArray[2] = byte(v >> 8)
	c.tempAddrArray[3] = byte(v)

	return c.tempAddrArray
}

func (c *splitRangeCollectorV4) bitsLen(n uint32) int {
	return math_bits.Len32(n)
}

func (c *splitRangeCollectorV4) countRighthandZeroBits(num uint32, bits int) int {
	if num == 0 {
		return bits
	}

	n := ^num & (num - 1)

	return minInt(bits, c.bitsLen(n))
}

type splitRangeCollectorV6 struct {
	first  *big.Int
	second *big.Int

	ipBits        int
	tempAddrArray ipV6
}

func newSplitRangeCollectorV6(first *netip.Addr, second *netip.Addr) *splitRangeCollectorV6 {
	fs := newInt()
	ss := newInt()

	first16 := first.As16()
	second16 := second.As16()

	fs = fs.SetBytes(first16[:])
	ss = ss.SetBytes(second16[:])

	arr := ipV6{}

	return &splitRangeCollectorV6{
		first:         fs,
		second:        ss,
		ipBits:        first.BitLen(),
		tempAddrArray: arr,
	}
}

func (c *splitRangeCollectorV6) Iter() iter.Seq[netip.Prefix] {
	return func(yield func(netip.Prefix) bool) {
		for c.first.Cmp(c.second) <= 0 {
			subPlusOne := copyInt(c.second)
			subPlusOne = subPlusOne.Sub(subPlusOne, c.first)
			subPlusOne = subPlusOne.Add(subPlusOne, bIntOne)

			nbits := minInt(
				c.countRighthandZeroBits(c.first, c.ipBits),
				c.bitsLen(subPlusOne)-1,
			)

			if !yield(
				netip.PrefixFrom(
					netip.AddrFrom16(c.addr(c.first)),
					c.ipBits-nbits,
				),
			) {
				return
			}

			appendBits := big.NewInt(1)

			appendBits = appendBits.Lsh(appendBits, uint(nbits))

			c.first = c.first.Add(c.first, appendBits)

			firstPlusOne := copyInt(c.first)
			firstPlusOne = firstPlusOne.Add(firstPlusOne, bIntOne)

			if firstPlusOne.Cmp(allOnesV6) == 0 {
				return
			}
		}
	}
}

func (c *splitRangeCollectorV6) bitsLen(num *big.Int) int {
	return num.BitLen()
}

func (c *splitRangeCollectorV6) addr(num *big.Int) ipV6 {
	for i := 0; i < len(c.tempAddrArray); i++ {
		c.tempAddrArray[i] = 0
	}

	bts := num.Bytes()

	j := 0
	for i := len(c.tempAddrArray) - len(bts); i < len(c.tempAddrArray); i++ {
		c.tempAddrArray[i] = bts[j]
		j++
	}

	return c.tempAddrArray
}

func (c *splitRangeCollectorV6) countRighthandZeroBits(num *big.Int, bits int) int {
	if num.Cmp(bIntZero) == 0 {
		return bits
	}

	numMinusOne := copyInt(num)
	numMinusOne = numMinusOne.Sub(numMinusOne, bIntOne)

	numInvert := copyInt(num)
	numInvert = numInvert.Not(numInvert)

	n := numInvert.And(numInvert, numMinusOne)

	return minInt(bits, c.bitsLen(n))
}
