// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package bufferpool

import (
	"bytes"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
)

func TestBufferPool(t *testing.T) {
	// TODO chef: assert result

	strategyList := []Strategy{
		StrategySingleStdPoolBucket,
		StrategySingleSlicePoolBucket,
		StategyMultiStdPoolBucket,
		StategyMultiSlicePoolBucket,
	}

	for _, s := range strategyList {
		bp := NewBufferPool(s)
		buf := &bytes.Buffer{}
		bp.Get(128)
		bp.Put(buf)
		buf = bp.Get(128)
		buf.Grow(4096)
		bp.Put(buf)
		buf = bp.Get(4096)
		bp.Put(buf)
		bp.RetrieveStatus()
	}
}

func TestGlobal(t *testing.T) {
	buf := Get(128)
	Put(buf)
	RetrieveStatus()
}

func TestSliceBucket(t *testing.T) {
	sb := NewSliceBucket()
	buf := sb.Get()
	assert.Equal(t, nil, buf)
	sb.Put(&bytes.Buffer{})
	buf = sb.Get()
	assert.IsNotNil(t, buf)
}

func TestUp2power(t *testing.T) {
	assert.Equal(t, 2, up2power(0))
	assert.Equal(t, 2, up2power(1))
	assert.Equal(t, 2, up2power(2))
	assert.Equal(t, 4, up2power(3))
	assert.Equal(t, 4, up2power(4))
	assert.Equal(t, 8, up2power(5))
	assert.Equal(t, 8, up2power(6))
	assert.Equal(t, 8, up2power(7))
	assert.Equal(t, 8, up2power(8))
	assert.Equal(t, 16, up2power(9))
	assert.Equal(t, 1024, up2power(1023))
	assert.Equal(t, 1024, up2power(1024))
	assert.Equal(t, 2048, up2power(1025))
	assert.Equal(t, 1073741824, up2power(1073741824-1))
	assert.Equal(t, 1073741824, up2power(1073741824))
	assert.Equal(t, 1073741824+1, up2power(1073741824+1))
	assert.Equal(t, 2047483647-1, up2power(2047483647-1))
	assert.Equal(t, 2047483647, up2power(2047483647))
}

func TestDown2power(t *testing.T) {
	assert.Equal(t, 2, down2power(0))
	assert.Equal(t, 2, down2power(1))
	assert.Equal(t, 2, down2power(2))
	assert.Equal(t, 2, down2power(3))
	assert.Equal(t, 4, down2power(4))
	assert.Equal(t, 4, down2power(5))
	assert.Equal(t, 4, down2power(6))
	assert.Equal(t, 4, down2power(7))
	assert.Equal(t, 8, down2power(8))
	assert.Equal(t, 8, down2power(9))
	assert.Equal(t, 512, down2power(1023))
	assert.Equal(t, 1024, down2power(1024))
	assert.Equal(t, 1024, down2power(1025))
	assert.Equal(t, 1073741824>>1, down2power(1073741824-1))
	assert.Equal(t, 1073741824, down2power(1073741824))
	assert.Equal(t, 1073741824, down2power(1073741824+1))
	assert.Equal(t, 1073741824, down2power(2047483647-1))
	assert.Equal(t, 1073741824, down2power(2047483647))
}
