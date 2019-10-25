// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package slicebytepool

var defaultPool SliceBytePool

func Get(size int) []byte {
	return defaultPool.Get(size)
}

func Put(buf []byte) {
	defaultPool.Put(buf)
}

func RetrieveStatus() Status {
	return defaultPool.RetrieveStatus()
}

func Init(strategy Strategy) {
	defaultPool = NewSliceBytePool(strategy)
}

func init() {
	Init(StrategyMultiSlicePoolBucket)
}
