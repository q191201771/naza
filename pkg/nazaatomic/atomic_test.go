// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazaatomic

import (
	"testing"

	"github.com/q191201771/naza/pkg/assert"
)

func TestBool(t *testing.T) {
	var v Bool
	assert.Equal(t, false, v.Load())
	v.Store(true)
	assert.Equal(t, true, v.Load())
	assert.Equal(t, false, v.CompareAndSwap(false, true))
	assert.Equal(t, true, v.Load())
	assert.Equal(t, true, v.CompareAndSwap(true, false))
	assert.Equal(t, false, v.Load())
	assert.Equal(t, false, v.Swap(true))
	assert.Equal(t, true, v.Load())
}

func TestInt32(t *testing.T) {
	var v Int32
	assert.Equal(t, int32(0), v.Load())
	v.Store(123)
	assert.Equal(t, int32(123), v.Load())
	v.Add(100)
	assert.Equal(t, int32(223), v.Load())
	v.Sub(200)
	assert.Equal(t, int32(23), v.Load())
	v.Increment()
	assert.Equal(t, int32(24), v.Load())
	v.Decrement()
	assert.Equal(t, int32(23), v.Load())
	assert.Equal(t, false, v.CompareAndSwap(0, 100))
	assert.Equal(t, int32(23), v.Load())
	assert.Equal(t, true, v.CompareAndSwap(23, 100))
	assert.Equal(t, int32(100), v.Load())
	assert.Equal(t, int32(100), v.Swap(200))
	assert.Equal(t, int32(200), v.Load())

	// 非常规操作
	v.Add(-50)
	assert.Equal(t, int32(150), v.Load())
	v.Sub(-60)
	assert.Equal(t, int32(210), v.Load())

	// 越界
	v.Store(0)
	v.Add(2147483647)
	assert.Equal(t, int32(2147483647), v.Load())
	v.Add(1)
	assert.Equal(t, int32(-2147483648), v.Load())
	v.Add(1)
	assert.Equal(t, int32(-2147483647), v.Load())
}

func TestUint32(t *testing.T) {
	var v Uint32
	assert.Equal(t, uint32(0), v.Load())
	v.Store(123)
	assert.Equal(t, uint32(123), v.Load())
	v.Add(100)
	assert.Equal(t, uint32(223), v.Load())
	v.Sub(200)
	assert.Equal(t, uint32(23), v.Load())
	v.Increment()
	assert.Equal(t, uint32(24), v.Load())
	v.Decrement()
	assert.Equal(t, uint32(23), v.Load())
	assert.Equal(t, false, v.CompareAndSwap(0, 100))
	assert.Equal(t, uint32(23), v.Load())
	assert.Equal(t, true, v.CompareAndSwap(23, 100))
	assert.Equal(t, uint32(100), v.Load())
	assert.Equal(t, uint32(100), v.Swap(200))
	assert.Equal(t, uint32(200), v.Load())

	// 越界
	v.Store(0)
	v.Add(4294967295)
	assert.Equal(t, uint32(4294967295), v.Load())
	v.Add(1)
	assert.Equal(t, uint32(0), v.Load())
	v.Add(1)
	assert.Equal(t, uint32(1), v.Load())
}

func TestInt64(t *testing.T) {
	var v Int64
	assert.Equal(t, int64(0), v.Load())
	v.Store(123)
	assert.Equal(t, int64(123), v.Load())
	v.Add(100)
	assert.Equal(t, int64(223), v.Load())
	v.Sub(200)
	assert.Equal(t, int64(23), v.Load())
	v.Increment()
	assert.Equal(t, int64(24), v.Load())
	v.Decrement()
	assert.Equal(t, int64(23), v.Load())
	assert.Equal(t, false, v.CompareAndSwap(0, 100))
	assert.Equal(t, int64(23), v.Load())
	assert.Equal(t, true, v.CompareAndSwap(23, 100))
	assert.Equal(t, int64(100), v.Load())
	assert.Equal(t, int64(100), v.Swap(200))
	assert.Equal(t, int64(200), v.Load())

	// 非常规操作
	v.Add(-50)
	assert.Equal(t, int64(150), v.Load())
	v.Sub(-60)
	assert.Equal(t, int64(210), v.Load())
}

func TestUint64(t *testing.T) {
	var v Uint64
	assert.Equal(t, uint64(0), v.Load())
	v.Store(123)
	assert.Equal(t, uint64(123), v.Load())
	v.Add(100)
	assert.Equal(t, uint64(223), v.Load())
	v.Sub(200)
	assert.Equal(t, uint64(23), v.Load())
	v.Increment()
	assert.Equal(t, uint64(24), v.Load())
	v.Decrement()
	assert.Equal(t, uint64(23), v.Load())
	assert.Equal(t, false, v.CompareAndSwap(0, 100))
	assert.Equal(t, uint64(23), v.Load())
	assert.Equal(t, true, v.CompareAndSwap(23, 100))
	assert.Equal(t, uint64(100), v.Load())
	assert.Equal(t, uint64(100), v.Swap(200))
	assert.Equal(t, uint64(200), v.Load())
}
