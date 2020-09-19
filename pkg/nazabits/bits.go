// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazabits

import "errors"

// TODO
// - 这个package的性能可以优化
// - BitReader考虑增加skip接口
// - BitReader考虑增加返回目前位置接口

var ErrNazaBits = errors.New("nazabits: fxxk")

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

func (br *BitReader) ReadBit() (uint8, error) {
	if br.index >= uint(len(br.core)) {
		return 0, ErrNazaBits
	}
	res := GetBit8(br.core[br.index], 7-br.pos)
	br.pos++
	if br.pos == 8 {
		br.pos = 0
		br.index++
	}
	return res, nil
}

// @param n: 取值范围 [1, 8]
func (br *BitReader) ReadBits8(n uint) (r uint8, err error) {
	var t uint8
	for i := uint(0); i < n; i++ {
		t, err = br.ReadBit()
		if err != nil {
			return
		}
		r = (r << 1) | t
	}
	return
}

// @param n: 取值范围 [1, 16]
func (br *BitReader) ReadBits16(n uint) (r uint16, err error) {
	var t uint8
	for i := uint(0); i < n; i++ {
		t, err = br.ReadBit()
		if err != nil {
			return
		}
		r = (r << 1) | uint16(t)
	}
	return
}

// @param n: 取值范围 [1, 32]
func (br *BitReader) ReadBits32(n uint) (r uint32, err error) {
	var t uint8
	for i := uint(0); i < n; i++ {
		t, err = br.ReadBit()
		if err != nil {
			return
		}
		r = (r << 1) | uint32(t)
	}
	return
}

// @param n: 取值范围 [1, 64]
func (br *BitReader) ReadBits64(n uint) (r uint64, err error) {
	var t uint8
	for i := uint(0); i < n; i++ {
		t, err = br.ReadBit()
		if err != nil {
			return
		}
		r = (r << 1) | uint64(t)
	}
	return
}

// @param n: 读取多少个字节
func (br *BitReader) ReadBytes(n uint) (r []byte, err error) {
	var t uint8
	for i := uint(0); i < n; i++ {
		t, err = br.ReadBits8(8)
		if err != nil {
			return
		}
		r = append(r, t)
	}
	return
}

// 0阶指数哥伦布编码
func (br *BitReader) ReadGolomb() (v uint32, err error) {
	var t uint8
	var n uint
	var m uint32
	for {
		t, err = br.ReadBit()
		if err != nil {
			return
		}
		if t == 0 {
			n++
		} else {
			break
		}
	}
	m, err = br.ReadBits32(n)
	if err != nil {
		return
	}
	v = 1<<n + m - 1
	return
}

// ----------------------------------------------------------------------------

// TODO chef: BitWriter没有对写越界做检查，由调用方保证这一点，后续可能会加上检查

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
	if b&0x1 == 1 {
		bw.core[bw.index] |= 1 << (7 - bw.pos)
	} else {
		bw.core[bw.index] &= ^(1 << (7 - bw.pos))
	}
	bw.pos++
	if bw.pos == 8 {
		bw.pos = 0
		bw.index++
	}
}

// 将<v>的低<n>位写入
// @param n: 取值范围 [1, 8]
func (bw *BitWriter) WriteBits8(n uint, v uint8) {
	for i := n - 1; ; i-- {
		bw.WriteBit(v >> i & 0x1)
		if i == 0 {
			break
		}
	}
}

func (bw *BitWriter) WriteBits16(n uint, v uint16) {
	for i := n - 1; ; i-- {
		bw.WriteBit(uint8(v >> i & 0x1))
		if i == 0 {
			break
		}
	}
}

// ----------------------------------------------------------------------------

// TODO chef: func GetBitX和func GetBitsX没有对写越界做检查，由调用方保证这一点，后续可能会加上检查

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
