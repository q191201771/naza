// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package bufferpool

import "bytes"

// TODO chef: 合适的释放接口

type BufferPool interface {
	Get(size int) *bytes.Buffer
	Put(buf *bytes.Buffer)
	RetrieveStatus() Status
}

type Status struct {
	getCount    int64
	putCount    int64
	hitCount    int64
	mallocCount int64
	sizeBytes   int64
}

func NewBufferPool() BufferPool {
	capToFreeBucket := make(map[int]*item)
	for i := minSize; i <= maxSize; i <<= 1 {
		capToFreeBucket[i] = new(item)
	}

	return &bufferPool{
		capToFreeBucket: capToFreeBucket,
	}
}
