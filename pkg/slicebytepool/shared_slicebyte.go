// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package slicebytepool

import (
	"github.com/q191201771/naza/pkg/nazaatomic"
)

type SharedSliceByte struct {
	Core  []byte
	pool  SliceBytePool
	count nazaatomic.Uint32
}

type SharedSliceByteOption struct {
	pool SliceBytePool
}

var defaultSharedSliceByteOption = SharedSliceByteOption{
	pool: defaultPool,
}

type ModSharedSliceByteOption func(option *SharedSliceByteOption)

func WithPool(pool SliceBytePool) ModSharedSliceByteOption {
	return func(option *SharedSliceByteOption) {
		option.pool = pool
	}
}

func NewSharedSliceByte(size int, modOptions ...ModSharedSliceByteOption) *SharedSliceByte {
	option := defaultSharedSliceByteOption
	for _, fn := range modOptions {
		fn(&option)
	}

	var ssb SharedSliceByte
	ssb.Core = option.pool.Get(size)
	ssb.pool = option.pool
	ssb.count.Store(1)
	return &ssb
}

func WrapSharedSliceByte(b []byte, modOptions ...ModSharedSliceByteOption) *SharedSliceByte {
	option := SharedSliceByteOption{
		pool: defaultPool,
	}
	for _, fn := range modOptions {
		fn(&option)
	}

	var ssb SharedSliceByte
	ssb.Core = b
	ssb.pool = option.pool
	ssb.count.Store(1)
	return &ssb
}

func (ssb *SharedSliceByte) Ref() *SharedSliceByte {
	ssb.count.Increment()
	return ssb
}

func (ssb *SharedSliceByte) ReleaseIfNeeded() {
	if ssb.count.Decrement() == 0 {
		ssb.pool.Put(ssb.Core)
	}
}
