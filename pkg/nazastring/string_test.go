// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazastring

import (
	"bytes"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
)

var inbuf = bytes.Repeat([]byte{'1', '2', '3', '4'}, 5678)
var instr = string(inbuf)

func TestSliceByteToStringTmp(t *testing.T) {
	str := SliceByteToStringTmp(inbuf)
	assert.Equal(t, instr, str)
}

func TestStringToSliceByteTmp(t *testing.T) {
	buf := StringToSliceByteTmp(instr)
	assert.Equal(t, inbuf, buf)
}

func TestDumpSliceByte(t *testing.T) {
	golden := []byte{1, 2, 3, 4, 5}
	ret := DumpSliceByte(golden)
	assert.Equal(t, "[]byte{0x01, 0x02, 0x03, 0x04, 0x05}", ret)
}

func TestSubSliceSafety(t *testing.T) {
	var b []byte
	assert.Equal(t, nil, SubSliceSafety(b, 1))
	assert.Equal(t, nil, SubSliceSafety(b, 2))
	b = []byte{1}
	assert.Equal(t, b, SubSliceSafety(b, 1))
	assert.Equal(t, b, SubSliceSafety(b, 2))
	b = []byte{1, 2}
	assert.Equal(t, []byte{1}, SubSliceSafety(b, 1))
	assert.Equal(t, b, SubSliceSafety(b, 2))
	assert.Equal(t, b, SubSliceSafety(b, 3))
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
		str = SliceByteToStringTmp(inbuf)
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
		buf = StringToSliceByteTmp(instr)
	}
	assert.Equal(b, buf, inbuf)
}
