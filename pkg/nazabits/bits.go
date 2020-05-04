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
	index int
	pos   int // 从左往右
}

func NewBitReader(b []byte) BitReader {
	return BitReader{
		core: b,
	}
}

func (br *BitReader) ReadBit() int {
	res := GetBit8(br.core[br.index], 7-br.pos)
	br.pos++
	if br.pos == 8 {
		br.pos = 0
		br.index++
	}
	return res
}

// 从高位（左边）开始读取
func (br *BitReader) ReadBits(n int) int {
	var res int
	for i := 0; i < n; i++ {
		res = (res << 1) | br.ReadBit()
	}
	return res
}

// ----------------------------------------------------------------------------

// @param pos: 取值范围 [0, 7]，0表示最低位
func GetBit8(v uint8, pos int) int {
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
func GetBits8(v uint8, pos int, n int) int {
	return int(v >> uint(pos) & m[n])
}

func GetBit16(v []byte, pos int) int {
	if pos < 8 {
		return GetBit8(v[1], pos)
	}
	return GetBit8(v[0], pos-8)
}

func GetBits16(v []byte, pos int, n int) int {
	if pos < 8 {
		if pos+n < 9 {
			return GetBits8(v[1], pos, n)
		}
		return GetBits8(v[1], pos, 8-pos) | GetBits8(v[0], 0, pos+n-8)<<(8-uint8(pos))
	}

	return GetBits8(v[0], pos-8, n)
}

var m []uint8

func init() {
	m = []uint8{0, 1, 3, 7, 15, 31, 63, 127, 255} // 0 is dummy
}
