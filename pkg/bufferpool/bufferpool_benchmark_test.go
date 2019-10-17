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
	"math/rand"
	"testing"
	"time"
)

var bp BufferPool
var count int32

func size() int {
	return random(0, 128*1024)

	//return 128 * 1024

	//ss := []int{1000, 2000, 5000}
	//////ss := []int{128, 1024, 4096, 16384}
	//atomic.AddInt32(&count, 1)
	//return ss[count % 3]

	//count++
	//if count > 128 * 1024 {
	//	count = 1
	//}
	//return count
}

func random(l, r int) int {
	return l + (rand.Int() % (r - l))
}

func originFunc() {
	var buf bytes.Buffer
	size := size()
	buf.Grow(size)
}

func bufferPoolFunc() {
	size := size()
	buf := bp.Get(size)
	buf.Grow(size)
	bp.Put(buf)
}

func BenchmarkOrigin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		originFunc()
	}
}

func BenchmarkBufferPool(b *testing.B) {
	bp = NewBufferPool(StategyMultiStdPoolBucket)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bufferPoolFunc()
	}
	//nazalog.Debugf("%+v", bp.RetrieveStatus())
}

//func BenchmarkOriginParallel(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		b.RunParallel(func(pb *testing.PB) {
//			for pb.Next() {
//				originFunc()
//			}
//		})
//	}
//}
//
//func BenchmarkBufferPoolParallel(b *testing.B) {
//	bp = NewBufferPool()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		b.RunParallel(func(pb *testing.PB) {
//			for pb.Next() {
//				bufferPoolFunc()
//			}
//		})
//	}
//	//nazalog.Debugf("%+v", bp.RetrieveStatus())
//}

func init() {
	rand.Seed(time.Now().Unix())
}
