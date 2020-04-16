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
