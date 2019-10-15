// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package bufferpool

import "bytes"

// 使用sync.Pool作为底层bucket实现时，池内自由选择合适的时机，自动释放空闲Buffer

type BufferPool interface {
	// 获取一个已经预申请大于<size>大小的Buffer对象，如果池内没有满足条件的空闲Buffer，会向Go内存管理模块申请
	Get(size int) *bytes.Buffer

	// 将Buffer对象放回池中
	Put(buf *bytes.Buffer)

	// 获取池当前状态
	RetrieveStatus() Status
}

type Status struct {
	getCount    int64 // 调用Get方法的次数
	putCount    int64 // 调用Put方法的次数
	hitCount    int64 // 调用Get方法时，池内存在满足条件的空闲Buffer，这种情况的计数
	mallocCount int64 // 调用Get方法时，池内不存在满足条件的空闲Buffer，向Go内存管理模块申请，这种情况的计数
	sizeBytes   int64 // 池内所有空闲Buffer占用的内存大小，单位字节
}

func NewBufferPool() BufferPool {
	capToFreeBucket := make(map[int]Bucket)
	for i := minSize; i <= maxSize; i <<= 1 {
		capToFreeBucket[i] = NewStdPoolBucket()
		//capToFreeBucket[i] = NewSliceBucket()
	}

	return &bufferPool{
		capToFreeBucket: capToFreeBucket,
	}
}
