// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package slicebytepool

import "github.com/q191201771/naza/pkg/nazaatomic"

var (
	minSize = 1024
	maxSize = 1073741824
)

type sliceBytePool struct {
	strategy        Strategy
	capToFreeBucket map[int]Bucket
	status          statusAtomic
}

type statusAtomic struct {
	getCount  nazaatomic.Int64
	putCount  nazaatomic.Int64
	hitCount  nazaatomic.Int64
	sizeBytes nazaatomic.Int64
}

func (bp *sliceBytePool) Get(size int) []byte {
	bp.status.getCount.Increment()

	ss := up2power(size)
	if ss < minSize {
		ss = minSize
	}
	bucket := bp.capToFreeBucket[ss]

	buf := bucket.Get(size)
	if buf == nil {
		buf = make([]byte, size, ss)
		return buf
	}

	bp.status.hitCount.Increment()
	bp.status.sizeBytes.Sub(int64(cap(buf)))
	return buf
}

func (bp *sliceBytePool) Put(buf []byte) {
	c := cap(buf)
	bp.status.putCount.Increment()
	bp.status.sizeBytes.Add(int64(c))

	size := down2power(c)
	if size < minSize {
		size = minSize
	}

	bucket := bp.capToFreeBucket[size]

	bucket.Put(buf)
}

func (bp *sliceBytePool) RetrieveStatus() Status {
	return Status{
		getCount:  bp.status.getCount.Load(),
		putCount:  bp.status.putCount.Load(),
		hitCount:  bp.status.hitCount.Load(),
		sizeBytes: bp.status.sizeBytes.Load(),
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
