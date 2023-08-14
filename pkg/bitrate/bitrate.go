// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// Package bitrate 平滑计算比特率（码率）
package bitrate

import (
	"sync"
	"time"
)

type Bitrate interface {
	// Add
	//
	// @param nowUnixMs: 变参，可选择从外部传入当前 unix 时间戳，单位毫秒
	//
	Add(bytes int, nowUnixMs ...int64)

	Rate(nowUnixMs ...int64) float32
}

type Unit uint8

const (
	UnitBitPerSec Unit = iota + 1
	UnitBytePerSec
	UnitKbitPerSec
	UnitKbytePerSec
)

// TODO chef: 考虑支持配置是否在内部使用锁
type Option struct {
	WindowMs int
	Unit     Unit
}

var defaultOption = Option{
	WindowMs: 1000,
	Unit:     UnitKbitPerSec,
}

type ModOption func(option *Option)

func New(modOptions ...ModOption) Bitrate {
	option := defaultOption
	for _, fn := range modOptions {
		fn(&option)
	}
	return &bitrate{
		option: option,
	}
}

type bitrate struct {
	option Option

	mu          sync.Mutex
	bucketSlice []bucket
}

type bucket struct {
	n int
	t int64 // unix 时间戳，单位毫秒
}

func (b *bitrate) Add(bytes int, nowUnixMs ...int64) {
	var now int64
	if len(nowUnixMs) == 0 {
		now = time.Now().UnixNano() / 1e6
	} else {
		now = nowUnixMs[0]
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.sweepStale(now)
	b.bucketSlice = append(b.bucketSlice, bucket{
		n: bytes,
		t: now,
	})
}

func (b *bitrate) Rate(nowUnixMs ...int64) float32 {
	var now int64
	if len(nowUnixMs) == 0 {
		now = time.Now().UnixNano() / 1e6
	} else {
		now = nowUnixMs[0]
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.sweepStale(now)
	var total int
	for i := range b.bucketSlice {
		total += b.bucketSlice[i].n
	}

	var ret float32
	switch b.option.Unit {
	case UnitBitPerSec:
		ret = float32(total*8*1000) / float32(b.option.WindowMs)
	case UnitBytePerSec:
		ret = float32(total*1000) / float32(b.option.WindowMs)
	case UnitKbitPerSec:
		ret = float32(total*8) / float32(b.option.WindowMs)
	case UnitKbytePerSec:
		ret = float32(total) / float32(b.option.WindowMs)
	}
	return ret
}

func (b *bitrate) sweepStale(now int64) {
	i := 0
	l := len(b.bucketSlice)
	for ; i < l; i++ {
		if now-b.bucketSlice[i].t <= int64(b.option.WindowMs) {
			break
		}
	}
	if i == l {
		b.bucketSlice = nil
	} else {
		b.bucketSlice = b.bucketSlice[i:]
	}
}
