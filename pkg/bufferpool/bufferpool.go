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
	"sync/atomic"
)

var (
	minSize = 1024
	maxSize = 1073741824
)

type bufferPool struct {
	strategy        Strategy
	singleBucket    Bucket
	capToFreeBucket map[int]Bucket
	status          Status
}

func (bp *bufferPool) Get(size int) *bytes.Buffer {
	atomic.AddInt64(&bp.status.getCount, 1)

	var bucket Bucket
	if bp.strategy == StategyMultiStdPoolBucket || bp.strategy == StategyMultiSlicePoolBucket {
		ss := up2power(size)
		if ss < minSize {
			ss = minSize
		}
		bucket = bp.capToFreeBucket[ss]
	} else {
		bucket = bp.singleBucket
	}

	buf := bucket.Get()
	if buf == nil {
		return &bytes.Buffer{}
	}

	atomic.AddInt64(&bp.status.hitCount, 1)
	atomic.AddInt64(&bp.status.sizeBytes, int64(-buf.Cap()))
	return buf
}

func (bp *bufferPool) Put(buf *bytes.Buffer) {
	c := buf.Cap()
	atomic.AddInt64(&bp.status.putCount, 1)
	atomic.AddInt64(&bp.status.sizeBytes, int64(c))

	var bucket Bucket
	if bp.strategy == StategyMultiStdPoolBucket || bp.strategy == StategyMultiSlicePoolBucket {
		size := down2power(c)
		if size < minSize {
			size = minSize
		}

		bucket = bp.capToFreeBucket[size]
	} else {
		bucket = bp.singleBucket
	}

	bucket.Put(buf)
}

func (bp *bufferPool) RetrieveStatus() Status {
	return Status{
		getCount:  atomic.LoadInt64(&bp.status.getCount),
		putCount:  atomic.LoadInt64(&bp.status.putCount),
		hitCount:  atomic.LoadInt64(&bp.status.hitCount),
		sizeBytes: atomic.LoadInt64(&bp.status.sizeBytes),
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
