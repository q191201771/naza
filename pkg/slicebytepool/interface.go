// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package slicebytepool

type SliceBytePool interface {
	Get(size int) []byte

	Put(buf []byte)

	RetrieveStatus() Status
}

type Status struct {
	getCount  int64
	putCount  int64
	hitCount  int64
	sizeBytes int64
}

type Strategy int

const (
	StrategyMultiStdPoolBucket = iota + 1

	StrategyMultiSlicePoolBucket
)

type Bucket interface {
	Get(size int) []byte

	Put(buf []byte)
}

func NewSliceBytePool(strategy Strategy) SliceBytePool {
	var capToFreeBucket map[int]Bucket

	switch strategy {
	case StrategyMultiStdPoolBucket:
		capToFreeBucket = make(map[int]Bucket)
		for i := minSize; i <= maxSize; i <<= 1 {
			capToFreeBucket[i] = NewStdPoolBucket()
		}
	case StrategyMultiSlicePoolBucket:
		capToFreeBucket = make(map[int]Bucket)
		for i := minSize; i <= maxSize; i <<= 1 {
			capToFreeBucket[i] = NewSliceBucket()
		}
	}

	return &sliceBytePool{
		strategy:        strategy,
		capToFreeBucket: capToFreeBucket,
	}
}
