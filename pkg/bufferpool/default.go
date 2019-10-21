// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package bufferpool

import "bytes"

var defaultPool BufferPool

func Get(size int) *bytes.Buffer {
	return defaultPool.Get(size)
}

func Put(buf *bytes.Buffer) {
	defaultPool.Put(buf)
}

func RetrieveStatus() Status {
	return defaultPool.RetrieveStatus()
}

func Init(strategy Strategy) {
	defaultPool = NewBufferPool(strategy)
}

func init() {
	Init(StrategyMultiStdPoolBucket)
}
