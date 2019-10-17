// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package bufferpool

import "bytes"

var global BufferPool

func Get(size int) *bytes.Buffer {
	return global.Get(size)
}

func Put(buf *bytes.Buffer) {
	global.Put(buf)
}

func RetrieveStatus() Status {
	return global.RetrieveStatus()
}

func Init(strategy Strategy) {
	global = NewBufferPool(strategy)
}

func init() {
	Init(StategyMultiStdPoolBucket)
}
