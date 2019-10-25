// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package slicebytepool

import (
	"sync"
)

type StdPoolBucket struct {
	core *sync.Pool
}

func NewStdPoolBucket() *StdPoolBucket {
	return &StdPoolBucket{
		core: new(sync.Pool),
	}
}

func (b *StdPoolBucket) Get(size int) []byte {
	v := b.core.Get()
	if v == nil {
		return nil
	}
	vv := v.([]byte)
	return vv[0:size]
}

func (b *StdPoolBucket) Put(buf []byte) {
	b.core.Put(buf)
}
