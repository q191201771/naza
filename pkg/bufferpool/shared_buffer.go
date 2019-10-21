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

	"github.com/q191201771/naza/pkg/nazaatomic"
)

type SharedBuffer struct {
	*bytes.Buffer
	pool  BufferPool
	count *nazaatomic.Uint32
}

func NewSharedBufferDefault(size int) *SharedBuffer {
	return NewSharedBuffer(defaultPool, size)
}

func NewSharedBuffer(pool BufferPool, size int) *SharedBuffer {
	return &SharedBuffer{
		Buffer: pool.Get(size),
		count:  new(nazaatomic.Uint32),
	}
}

func (sb *SharedBuffer) Ref() *SharedBuffer {
	sb.count.Increment()
	return sb
}

func (sb *SharedBuffer) ReleaseIfNeeded() {
	if sb.count.Decrement() == 0 {
		sb.pool.Put(sb.Buffer)
	}
}
