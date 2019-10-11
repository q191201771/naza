// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package mockwriter

import (
	"bytes"
	"errors"
)

// TODO chef: 可以添加一个接口，获取内部 buffer 的内容

type WriterType uint8

const (
	WriterTypeDoNothing WriterType = iota
	WriterTypeReturnError
	WriterTypeIntoBuffer
)

var (
	mockWriterErr = errors.New("mockwriter: a mock error")
)

type MockWriter struct {
	t     WriterType
	ts    map[uint32]WriterType
	count uint32
	b     bytes.Buffer
}

func NewMockWriter(t WriterType) *MockWriter {
	return &MockWriter{
		t: t,
	}
}

// 为某些写操作指定特定的类型，次数从 0 开始计数
func (w *MockWriter) SetSpecificType(ts map[uint32]WriterType) {
	w.ts = ts
}

func (w *MockWriter) Write(b []byte) (int, error) {
	t, exist := w.ts[w.count]
	w.count++
	if !exist {
		t = w.t
	}
	switch t {
	case WriterTypeDoNothing:
		return len(b), nil
	case WriterTypeReturnError:
		return 0, mockWriterErr
		//case WriterTypeIntoBuffer:
		//	return w.b.Write(b)
	}

	return w.b.Write(b)
	//panic("never reach here.")
}
