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
	status          Status
	capToFreeBucket map[int]*item
}

type item struct {
	m    sync.Mutex
	core []*bytes.Buffer
}

func (bp *bufferPool) Get(size int) *bytes.Buffer {
	atomic.AddInt64(&bp.status.getCount, 1)
	ss := up2power(size)
	if ss < minSize {
		ss = minSize
	}

	bucket := bp.capToFreeBucket[ss]
	bucket.m.Lock()
	if len(bucket.core) == 0 {
		bucket.m.Unlock()
		var buf bytes.Buffer
		buf.Grow(ss)
		atomic.AddInt64(&bp.status.mallocCount, 1)
		return &buf
	} else {
		buf := bucket.core[len(bucket.core)-1]
		bucket.core = bucket.core[:len(bucket.core)-1]
		bucket.m.Unlock()
		buf.Reset()
		atomic.AddInt64(&bp.status.hitCount, 1)
		atomic.AddInt64(&bp.status.sizeBytes, int64(-buf.Cap()))
		return buf
	}
}

func (bp *bufferPool) Put(buf *bytes.Buffer) {
	c := buf.Cap()
	atomic.AddInt64(&bp.status.putCount, 1)
	atomic.AddInt64(&bp.status.sizeBytes, int64(c))
	size := down2power(c)
	if size < minSize {
		size = minSize
	}

	bucket := bp.capToFreeBucket[size]
	bucket.m.Lock()
	bucket.core = append(bucket.core, buf)
	bucket.m.Unlock()
}

func (bp *bufferPool) RetrieveStatus() Status {
	return Status{
		getCount:    atomic.LoadInt64(&bp.status.getCount),
		putCount:    atomic.LoadInt64(&bp.status.putCount),
		hitCount:    atomic.LoadInt64(&bp.status.hitCount),
		mallocCount: atomic.LoadInt64(&bp.status.mallocCount),
		sizeBytes:   atomic.LoadInt64(&bp.status.sizeBytes),
	}
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
