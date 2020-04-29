// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazabits_test

import (
	"testing"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazabits"
)

func TestGetBit8(t *testing.T) {
	assert.Equal(t, 0, nazabits.GetBit8(0, 0))
	assert.Equal(t, 1, nazabits.GetBit8(1<<1, 1))
	assert.Equal(t, 1, nazabits.GetBit8(1<<1+1<<2, 2))
	assert.Equal(t, 1, nazabits.GetBit8(1<<3+1<<4+1<<5, 3))
	assert.Equal(t, 0, nazabits.GetBit8(0, 4))
	assert.Equal(t, 1, nazabits.GetBit8(1<<5, 5))
	assert.Equal(t, 0, nazabits.GetBit8(1+1<<5+1<<7, 6))
	assert.Equal(t, 0, nazabits.GetBit8(0, 7))
}

func TestGetBits8(t *testing.T) {
	// 0110 1001 = 1 + 8 + 32 + 64 = 105
	v := uint8(105)
	assert.Equal(t, 1, nazabits.GetBits8(v, 0, 1))
	assert.Equal(t, 0, nazabits.GetBits8(v, 7, 1))
	assert.Equal(t, 1, nazabits.GetBits8(v, 0, 2))
	assert.Equal(t, 1, nazabits.GetBits8(v, 6, 2))
	assert.Equal(t, 1, nazabits.GetBits8(v, 0, 3))
	assert.Equal(t, 5, nazabits.GetBits8(v, 3, 3))
	assert.Equal(t, 10, nazabits.GetBits8(v, 2, 4))
	assert.Equal(t, 26, nazabits.GetBits8(v, 2, 5))
	assert.Equal(t, 26, nazabits.GetBits8(v, 2, 6))
	assert.Equal(t, 105, nazabits.GetBits8(v, 0, 7))
	assert.Equal(t, 105, nazabits.GetBits8(v, 0, 8))
}

func TestGetBit16(t *testing.T) {
	v := []byte{48, 57} // 12345 {0011 0000, 0011 1001}
	assert.Equal(t, 1, nazabits.GetBit16(v, 0))
	assert.Equal(t, 0, nazabits.GetBit16(v, 1))
	assert.Equal(t, 0, nazabits.GetBit16(v, 2))
	assert.Equal(t, 1, nazabits.GetBit16(v, 3))
	assert.Equal(t, 1, nazabits.GetBit16(v, 4))
	assert.Equal(t, 1, nazabits.GetBit16(v, 5))
	assert.Equal(t, 0, nazabits.GetBit16(v, 6))
	assert.Equal(t, 0, nazabits.GetBit16(v, 7))
	assert.Equal(t, 0, nazabits.GetBit16(v, 8))
	assert.Equal(t, 0, nazabits.GetBit16(v, 9))
	assert.Equal(t, 0, nazabits.GetBit16(v, 10))
	assert.Equal(t, 0, nazabits.GetBit16(v, 11))
	assert.Equal(t, 1, nazabits.GetBit16(v, 12))
	assert.Equal(t, 1, nazabits.GetBit16(v, 13))
	assert.Equal(t, 0, nazabits.GetBit16(v, 14))
	assert.Equal(t, 0, nazabits.GetBit16(v, 15))
}

func TestGetBits16(t *testing.T) {
	v := []byte{48, 57} // 12345 {0011 0000, 0011 1001}
	assert.Equal(t, 12345, nazabits.GetBits16(v, 0, 16))
	assert.Equal(t, 1, nazabits.GetBits16(v, 0, 1))
	assert.Equal(t, 0, nazabits.GetBits16(v, 1, 2))
	assert.Equal(t, 6, nazabits.GetBits16(v, 2, 3))
	assert.Equal(t, 7, nazabits.GetBits16(v, 3, 4))
	assert.Equal(t, 3, nazabits.GetBits16(v, 4, 5))
	assert.Equal(t, 1, nazabits.GetBits16(v, 5, 6))
	assert.Equal(t, 64, nazabits.GetBits16(v, 6, 7))
	assert.Equal(t, 96, nazabits.GetBits16(v, 7, 8))
	assert.Equal(t, 0, nazabits.GetBits16(v, 8, 1))
	assert.Equal(t, 0, nazabits.GetBits16(v, 9, 2))
	assert.Equal(t, 4, nazabits.GetBits16(v, 10, 3))
	assert.Equal(t, 6, nazabits.GetBits16(v, 11, 4))
	assert.Equal(t, 3, nazabits.GetBits16(v, 12, 3))
	assert.Equal(t, 1, nazabits.GetBits16(v, 13, 2))
	assert.Equal(t, 0, nazabits.GetBits16(v, 14, 1))
	assert.Equal(t, 0, nazabits.GetBits16(v, 15, 1))
}
