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
	"sync"
	"sync/atomic"
)

var (
	minSize = 1024
	maxSize = 1073741824
)

type bufferPool struct {
	getCount    uint32
	putCount    uint32
	hitCount    uint32
	mallocCount uint32

	capToFreeBucket map[int]*item
}

type item struct {
	m    sync.Mutex
	core []*bytes.Buffer
}

func (bp *bufferPool) Get(size int) *bytes.Buffer {
	atomic.AddUint32(&bp.getCount, 1)
	ss := up2power(size)
	if ss < minSize {
		ss = minSize
	}

	bucket := bp.capToFreeBucket[ss]
	bucket.m.Lock()
	if len(bucket.core) == 0 {
		bucket.m.Unlock()
		return bp.newBuffer(ss)
	} else {
		buf := bucket.core[len(bucket.core)-1]
		bucket.core = bucket.core[:len(bucket.core)-1]
		bucket.m.Unlock()
		buf.Reset()
		atomic.AddUint32(&bp.hitCount, 1)
		return buf
	}
}

func (bp *bufferPool) Put(buf *bytes.Buffer) {
	atomic.AddUint32(&bp.putCount, 1)
	size := down2power(buf.Cap())
	if size < minSize {
		size = minSize
	}

	bucket := bp.capToFreeBucket[size]
	bucket.m.Lock()
	bucket.core = append(bucket.core, buf)
	bucket.m.Unlock()
}

func (bp *bufferPool) newBuffer(n int) *bytes.Buffer {
	var buf bytes.Buffer
	buf.Grow(n)
	atomic.AddUint32(&bp.mallocCount, 1)
	return &buf
}

// @return 范围为 [2, 4, 8, 16, ..., 1073741824]，如果大于等于1073741824，则直接返回n
func up2power(n int) int {
	if n >= maxSize {
		return n
	}

	var i uint32
	for ; n > (2 << i); i++ {
	}
	return 2 << i
}

// @return 范围为 [2, 4, 8, 16, ..., 1073741824]
func down2power(n int) int {
	if n < 2 {
		return 2
	} else if n >= maxSize {
		return maxSize
	}

	var i uint32
	for {
		nn := 2 << i
		if n > nn {
			i++
		} else if n == nn {
			return n
		} else if n < nn {
			return 2 << (i - 1)
		}
	}
}
