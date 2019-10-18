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

func (b *StdPoolBucket) Get() *bytes.Buffer {
	v := b.core.Get()
	if v == nil {
		return nil
	}
	return v.(*bytes.Buffer)

}

func (b *StdPoolBucket) Put(buf *bytes.Buffer) {
	b.core.Put(buf)
}
