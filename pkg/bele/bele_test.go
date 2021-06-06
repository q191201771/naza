// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package bele

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
)

func TestBeUint16(t *testing.T) {
	vector := []struct {
		input  []byte
		output uint16
	}{
		{input: []byte{0, 0}, output: 0},
		{input: []byte{0, 1}, output: 1},
		{input: []byte{0, 255}, output: 255},
		{input: []byte{1, 0}, output: 256},
		{input: []byte{255, 0}, output: 255 * 256},
		{input: []byte{12, 34}, output: 12*256 + 34},
	}

	for i := 0; i < len(vector); i++ {
		assert.Equal(t, vector[i].output, BeUint16(vector[i].input))
	}
}

func TestBeUint24(t *testing.T) {
	vector := []struct {
		input  []byte
		output uint32
	}{
		{input: []byte{0, 0, 0}, output: 0},
		{input: []byte{0, 0, 1}, output: 1},
		{input: []byte{0, 1, 0}, output: 256},
		{input: []byte{1, 0, 0}, output: 1 * 256 * 256},
		{input: []byte{12, 34, 56}, output: 12*256*256 + 34*256 + 56},
	}

	for i := 0; i < len(vector); i++ {
		assert.Equal(t, vector[i].output, BeUint24(vector[i].input))
	}
}

func TestBeUint32(t *testing.T) {
	vector := []struct {
		input  []byte
		output uint32
	}{
		{input: []byte{0, 0, 0, 0}, output: 0},
		{input: []byte{0, 0, 1, 0}, output: 1 * 256},
		{input: []byte{0, 1, 0, 0}, output: 1 * 256 * 256},
		{input: []byte{1, 0, 0, 0}, output: 1 * 256 * 256 * 256},
		{input: []byte{12, 34, 56, 78}, output: 12*256*256*256 + 34*256*256 + 56*256 + 78},
	}

	for i := 0; i < len(vector); i++ {
		assert.Equal(t, vector[i].output, BeUint32(vector[i].input))
	}
}

func TestBeUint64(t *testing.T) {
	vector := []struct {
		input  []byte
		output uint64
	}{
		{input: []byte{0, 0, 0, 0, 0, 0, 0, 0}, output: 0},
		{input: []byte{0, 0, 0, 0, 1, 0, 0, 0}, output: 1 * 256 * 256 * 256},
		{input: []byte{0, 0, 0, 0, 12, 34, 56, 78}, output: 12*256*256*256 + 34*256*256 + 56*256 + 78},
		{input: []byte{0, 12, 34, 56, 78, 0, 0, 0}, output: 12*256*256*256*256*256*256 + 34*256*256*256*256*256 + 56*256*256*256*256 + 78*256*256*256},
	}

	for i := 0; i < len(vector); i++ {
		assert.Equal(t, vector[i].output, BeUint64(vector[i].input))
	}
}

func TestBeFloat64(t *testing.T) {
	vector := []int{
		1,
		0xFF,
		0xFFFF,
		0xFFFFFF,
	}
	for i := 0; i < len(vector); i++ {
		b := &bytes.Buffer{}
		err := binary.Write(b, binary.BigEndian, float64(vector[i]))
		assert.Equal(t, nil, err)
		assert.Equal(t, vector[i], int(BeFloat64(b.Bytes())))
	}
}

func TestLeUint32(t *testing.T) {
	vector := []struct {
		input  []byte
		output uint32
	}{
		{input: []byte{0, 0, 0, 0}, output: 0},
		{input: []byte{0, 0, 1, 0}, output: 1 * 256 * 256},
		{input: []byte{0, 1, 0, 0}, output: 1 * 256},
		{input: []byte{1, 0, 0, 0}, output: 1},
		{input: []byte{12, 34, 56, 78}, output: 12 + 34*256 + 56*256*256 + 78*256*256*256},
	}

	for i := 0; i < len(vector); i++ {
		assert.Equal(t, vector[i].output, LeUint32(vector[i].input))
	}
}

func TestBePutUint16(t *testing.T) {
	b := make([]byte, 2)
	BePutUint16(b, 1)
	assert.Equal(t, []byte{0, 1}, b)
}

func TestBePutUint64(t *testing.T) {
	b := make([]byte, 8)
	BePutUint64(b, 1)
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 1}, b)
}

func TestBePutUint24(t *testing.T) {
	vector := []struct {
		input  uint32
		output []byte
	}{
		{input: 0, output: []byte{0, 0, 0}},
		{input: 1, output: []byte{0, 0, 1}},
		{input: 256, output: []byte{0, 1, 0}},
		{input: 1 * 256 * 256, output: []byte{1, 0, 0}},
		{input: 12*256*256 + 34*256 + 56, output: []byte{12, 34, 56}},
	}

	out := make([]byte, 3)
	for i := 0; i < len(vector); i++ {
		BePutUint24(out, vector[i].input)
		assert.Equal(t, vector[i].output, out)
	}
}

func TestBePutUint32(t *testing.T) {
	vector := []struct {
		input  uint32
		output []byte
	}{
		{input: 0, output: []byte{0, 0, 0, 0}},
		{input: 1 * 256, output: []byte{0, 0, 1, 0}},
		{input: 1 * 256 * 256, output: []byte{0, 1, 0, 0}},
		{input: 1 * 256 * 256 * 256, output: []byte{1, 0, 0, 0}},
		{input: 12*256*256*256 + 34*256*256 + 56*256 + 78, output: []byte{12, 34, 56, 78}},
	}

	out := make([]byte, 4)
	for i := 0; i < len(vector); i++ {
		BePutUint32(out, vector[i].input)
		assert.Equal(t, vector[i].output, out)
	}
}

func TestLePutUint32(t *testing.T) {
	vector := []struct {
		input  uint32
		output []byte
	}{
		{input: 0, output: []byte{0, 0, 0, 0}},
		{input: 1 * 256 * 256, output: []byte{0, 0, 1, 0}},
		{input: 1 * 256, output: []byte{0, 1, 0, 0}},
		{input: 1, output: []byte{1, 0, 0, 0}},
		{input: 78*256*256*256 + 56*256*256 + 34*256 + 12, output: []byte{12, 34, 56, 78}},
	}

	out := make([]byte, 4)
	for i := 0; i < len(vector); i++ {
		LePutUint32(out, vector[i].input)
		assert.Equal(t, vector[i].output, out)
	}
}

func TestWriteBeUint24(t *testing.T) {
	vector := []struct {
		input  uint32
		output []byte
	}{
		{input: 0, output: []byte{0, 0, 0}},
		{input: 1, output: []byte{0, 0, 1}},
		{input: 256, output: []byte{0, 1, 0}},
		{input: 1 * 256 * 256, output: []byte{1, 0, 0}},
		{input: 12*256*256 + 34*256 + 56, output: []byte{12, 34, 56}},
	}

	for i := 0; i < len(vector); i++ {
		out := &bytes.Buffer{}
		err := WriteBeUint24(out, vector[i].input)
		assert.Equal(t, nil, err)
		assert.Equal(t, vector[i].output, out.Bytes())
	}
}

func TestWriteBe(t *testing.T) {
	vector := []struct {
		input  interface{}
		output []byte
	}{
		{input: uint32(1), output: []byte{0, 0, 0, 1}},
		{input: uint64(1), output: []byte{0, 0, 0, 0, 0, 0, 0, 1}},
	}
	for i := 0; i < len(vector); i++ {
		out := &bytes.Buffer{}
		err := WriteBe(out, vector[i].input)
		assert.Equal(t, nil, err)
		assert.Equal(t, vector[i].output, out.Bytes())
	}
}

func TestWriteLe(t *testing.T) {
	vector := []struct {
		input  interface{}
		output []byte
	}{
		{input: uint32(1), output: []byte{1, 0, 0, 0}},
		{input: uint64(1), output: []byte{1, 0, 0, 0, 0, 0, 0, 0}},
	}
	for i := 0; i < len(vector); i++ {
		out := &bytes.Buffer{}
		err := WriteLe(out, vector[i].input)
		assert.Equal(t, nil, err)
		assert.Equal(t, vector[i].output, out.Bytes())
	}
}

func TestReadBytes(t *testing.T) {
	var buf bytes.Buffer
	buf.Write([]byte{'1', '2', '3'})
	b, err := ReadBytes(&buf, 2)
	assert.Equal(t, []byte{'1', '2'}, b)
	assert.Equal(t, nil, err)
	b, err = ReadBytes(&buf, 2)
	assert.Equal(t, []byte{'3', 0}, b)
	assert.IsNotNil(t, err)
	b, err = ReadBytes(&buf, 2)
	assert.Equal(t, nil, b)
	assert.IsNotNil(t, err)
}

func TestRead(t *testing.T) {
	var err error
	b := &bytes.Buffer{}
	_, err = ReadUint8(b)
	assert.IsNotNil(t, err)
	_, err = ReadBeUint16(b)
	assert.IsNotNil(t, err)
	_, err = ReadBeUint24(b)
	assert.IsNotNil(t, err)
	_, err = ReadBeUint32(b)
	assert.IsNotNil(t, err)
	_, err = ReadBeUint64(b)
	assert.IsNotNil(t, err)
	_, err = ReadLeUint32(b)
	assert.IsNotNil(t, err)

	b.Write([]byte{1})
	i8, err := ReadUint8(b)
	assert.Equal(t, uint8(1), i8)
	assert.Equal(t, nil, err)
	b.Write([]byte{1, 2})
	i16, err := ReadBeUint16(b)
	assert.Equal(t, BeUint16([]byte{1, 2}), i16)
	assert.Equal(t, nil, err)
	b.Write([]byte{1, 2, 3})
	i24, err := ReadBeUint24(b)
	assert.Equal(t, BeUint24([]byte{1, 2, 3}), i24)
	assert.Equal(t, nil, err)
	b.Write([]byte{1, 2, 3, 4})
	i32, err := ReadBeUint32(b)
	assert.Equal(t, BeUint32([]byte{1, 2, 3, 4}), i32)
	assert.Equal(t, nil, err)
	b.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	i64, err := ReadBeUint64(b)
	assert.Equal(t, BeUint64([]byte{1, 2, 3, 4, 5, 6, 7, 8}), i64)
	assert.Equal(t, nil, err)

	b.Write([]byte{1, 0, 0, 0})
	i32, err = ReadLeUint32(b)
	assert.Equal(t, uint32(1), i32)
	assert.Equal(t, nil, err)
}

func TestReadString(t *testing.T) {
	var buf bytes.Buffer
	str, err := ReadString(&buf, 2)
	assert.Equal(t, "", str)
	assert.IsNotNil(t, err)
}

func BenchmarkBeFloat64(b *testing.B) {
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, float64(123.4))
	for i := 0; i < b.N; i++ {
		BeFloat64(buf.Bytes())
	}
}

func BenchmarkBePutUint24(b *testing.B) {
	out := make([]byte, 3)
	for i := 0; i < b.N; i++ {
		BePutUint24(out, uint32(i))
	}
}

func BenchmarkBeUint24(b *testing.B) {
	buf := []byte{1, 2, 3}
	for i := 0; i < b.N; i++ {
		BeUint24(buf)
	}
}

func BenchmarkWriteBe(b *testing.B) {
	out := &bytes.Buffer{}
	in := uint64(123)
	for i := 0; i < b.N; i++ {
		_ = WriteBe(out, in)
	}
}
