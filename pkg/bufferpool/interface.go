// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package bufferpool

import "bytes"

type BufferPool interface {
	// 获取一个预估容量为<size>的Buffer对象，不同的策略返回的Buffer对象实际容量可能会有不同，但是不会返回nil
	Get(size int) *bytes.Buffer

	// 将Buffer对象放回池中
	Put(buf *bytes.Buffer)

	// 获取池当前状态
	RetrieveStatus() Status
}

type Status struct {
	getCount  int64 // 调用Get方法的次数
	putCount  int64 // 调用Put方法的次数
	hitCount  int64 // 调用Get方法时，池内存在满足条件的空闲Buffer，这种情况的计数
	sizeBytes int64 // 池内所有空闲Buffer占用的内存大小，单位字节
}

type Strategy int

const (
	// 直接使用一个标准库中的sync.Pool
	StrategySingleStdPoolBucket Strategy = iota + 1

	// 直接使用一个切片存储所有的Buffer对象
	StrategySingleSlicePoolBucket

	// 按Buffer对象的容量哈希到不同的桶中，每个桶是一个sync.Pool
	StategyMultiStdPoolBucket

	// 按Buffer对象的容量哈希到不同的桶中，每个桶是一个切片
	StategyMultiSlicePoolBucket
)

type Bucket interface {
	// 桶内没有Buffer对象时，返回nil
	Get() *bytes.Buffer

	Put(buf *bytes.Buffer)
}

func NewBufferPool(strategy Strategy) BufferPool {
	var (
		singleBucket    Bucket
		capToFreeBucket map[int]Bucket
	)

	switch strategy {
	case StrategySingleStdPoolBucket:
		singleBucket = NewStdPoolBucket()
	case StrategySingleSlicePoolBucket:
		singleBucket = NewSliceBucket()
	case StategyMultiStdPoolBucket:
		capToFreeBucket = make(map[int]Bucket)
		for i := minSize; i <= maxSize; i <<= 1 {
			capToFreeBucket[i] = NewStdPoolBucket()
		}
	case StategyMultiSlicePoolBucket:
		capToFreeBucket = make(map[int]Bucket)
		for i := minSize; i <= maxSize; i <<= 1 {
			capToFreeBucket[i] = NewSliceBucket()
		}
	}

	return &bufferPool{
		strategy:        strategy,
		singleBucket:    singleBucket,
		capToFreeBucket: capToFreeBucket,
	}
}
