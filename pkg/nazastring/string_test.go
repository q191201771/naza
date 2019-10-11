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
	"github.com/q191201771/naza/pkg/assert"
	"testing"
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
