// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazabytes

import (
	"bytes"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
)

var inbuf = bytes.Repeat([]byte{'1', '2', '3', '4'}, 5678)
var instr = string(inbuf)

func TestSliceByteToStringTmp(t *testing.T) {
	str := Bytes2StringRef(inbuf)
	assert.Equal(t, instr, str)
}

func TestStringToSliceByteTmp(t *testing.T) {
	buf := String2BytesRef(instr)
	assert.Equal(t, inbuf, buf)
}

func TestSubSliceSafety(t *testing.T) {
	var b []byte
	assert.Equal(t, nil, Sub(b, 0, 1))
	assert.Equal(t, nil, Sub(b, 0, 2))
	b = []byte{1}
	assert.Equal(t, b, Sub(b, 0, 1))
	assert.Equal(t, b, Sub(b, 0, 2))
	b = []byte{1, 2}
	assert.Equal(t, []byte{1}, Sub(b, 0, 1))
	assert.Equal(t, b, Sub(b, 0, 2))
	assert.Equal(t, b, Sub(b, 0, 3))

	assert.Equal(t, []byte{2}, Sub(b, 1, 1))
	assert.Equal(t, []byte{2}, Sub(b, 1, 2))
	assert.Equal(t, nil, Sub(b, 2, 1))
}

func BenchmarkSliceByteToStringOrigin(b *testing.B) {
	var str string

	for i := 0; i < b.N; i++ {
		str = string(inbuf)
	}
	assert.Equal(b, instr, str)
}

func BenchmarkSliceByteToStringTmp(b *testing.B) {
	var str string

	for i := 0; i < b.N; i++ {
		str = Bytes2StringRef(inbuf)
	}
	assert.Equal(b, instr, str)
}

func BenchmarkStringToSliceByteOrigin(b *testing.B) {
	var buf []byte

	for i := 0; i < b.N; i++ {
		buf = []byte(instr)
	}
	assert.Equal(b, buf, inbuf)
}

func BenchmarkStringToSliceByteTmp(b *testing.B) {
	var buf []byte

	for i := 0; i < b.N; i++ {
		buf = String2BytesRef(instr)
	}
	assert.Equal(b, buf, inbuf)
}
