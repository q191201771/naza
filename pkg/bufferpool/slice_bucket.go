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

type SliceBucket struct {
	m    sync.Mutex
	core []*bytes.Buffer
}

func NewSliceBucket() *SliceBucket {
	return new(SliceBucket)
}

func (b *SliceBucket) Get() *bytes.Buffer {
	b.m.Lock()
	defer b.m.Unlock()
	if len(b.core) == 0 {
		return nil
	}
	buf := b.core[len(b.core)-1]
	b.core = b.core[:len(b.core)-1]
	buf.Reset()
	return buf
}

func (b *SliceBucket) Put(buf *bytes.Buffer) {
	b.m.Lock()
	defer b.m.Unlock()
	b.core = append(b.core, buf)
}
