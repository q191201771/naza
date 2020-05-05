// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazabits

// TODO chef 这个package的性能可以优化

type BitReader struct {
	core  []byte
	index uint
	pos   uint // 从左往右
}

func NewBitReader(b []byte) BitReader {
	return BitReader{
		core: b,
	}
}

func (br *BitReader) ReadBit() uint8 {
	res := GetBit8(br.core[br.index], 7-br.pos)
	br.pos++
	if br.pos == 8 {
		br.pos = 0
		br.index++
	}
	return res
}

// 目前限制一次最多读8位，大于8位可以分多次读，或者以后把限制放开
func (br *BitReader) ReadBits(n uint) (r uint8) {
	var i uint
	for i = 0; i < n; i++ {
		r = (r << 1) | br.ReadBit()
	}
	return
}

// ----------------------------------------------------------------------------

type BitWriter struct {
	core  []byte
	index int
	pos   uint // 从左往右
}

func NewBitWriter(b []byte) BitWriter {
	return BitWriter{
		core: b,
	}
}

// @param b: 当b不为0和1时，取b的最低位
func (bw *BitWriter) WriteBit(b uint8) {
	bw.core[bw.index] |= (b & 0x1) << (7 - bw.pos)
	bw.pos++
	if bw.pos == 8 {
		bw.pos = 0
		bw.index++
	}
}

// 目前限制一次最多写8位，大于8位可以分多次写，或者以后把限制放开
func (bw *BitWriter) WriteBits(n int, v uint8) {
	for i := n - 1; i >= 0; i-- {
		bw.WriteBit(v >> uint8(i))
	}
}

// ----------------------------------------------------------------------------

// @param pos: 取值范围 [0, 7]，0表示最低位
func GetBit8(v uint8, pos uint) uint8 {
	return GetBits8(v, pos, 1)
}

// @param pos: 取值范围 [0, 7]，0表示最低位
// @param n:   取多少位， 取值范围 [1, 8]
//
// 举例，GetBits8(105, 2, 4) = 10（即1010）
//   v: 0110 1001
// pos:       2
//   n:   .. ..
//
func GetBits8(v uint8, pos uint, n uint) uint8 {
	return v >> pos & m[n]
}

func GetBit16(v []byte, pos uint) uint8 {
	if pos < 8 {
		return GetBit8(v[1], pos)
	}
	return GetBit8(v[0], pos-8)
}

func GetBits16(v []byte, pos uint, n uint) uint16 {
	if pos < 8 {
		if pos+n < 9 {
			return uint16(GetBits8(v[1], pos, n))
		}
		return uint16(GetBits8(v[1], pos, 8-pos)) | uint16(GetBits8(v[0], 0, pos+n-8))<<(8-pos)
	}

	return uint16(GetBits8(v[0], pos-8, n))
}

var m []uint8

func init() {
	m = []uint8{0, 1, 3, 7, 15, 31, 63, 127, 255} // 0 is dummy
}
