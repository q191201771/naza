// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package bitrate

import (
	"sync"
	"time"
)

// 先来个最简单的，后续再精细化配置：
//
// - 收包时间目前只能由内部获取当前时间，应提供接口支持外部传入
// - 返回的rate单位固定为 kbit/s
// - 不需要存储Time结构体，存毫秒级的 unix 时间戳

type Bitrate struct {
	windowMS int

	mu          sync.Mutex
	bucketSlice []bucket
}

type bucket struct {
	n int
	t time.Time
}

func NewBitrate(windowMS int) *Bitrate {
	return &Bitrate{
		windowMS: windowMS,
	}
}

func (b *Bitrate) Add(bytes int) {
	now := time.Now()
	b.mu.Lock()
	defer b.mu.Unlock()

	b.sweepStale(now)

	b.bucketSlice = append(b.bucketSlice, bucket{
		n: bytes,
		t: now,
	})
}

// @return 返回值单位 kbit/s
func (b *Bitrate) Rate() int {
	now := time.Now()
	b.mu.Lock()
	defer b.mu.Unlock()

	b.sweepStale(now)
	var total int
	for i := range b.bucketSlice {
		total += b.bucketSlice[i].n
	}

	// total * 8 / 1000 * 1000 / b.windowMS
	return total * 8 / b.windowMS
}

func (b *Bitrate) sweepStale(now time.Time) {
	for i := range b.bucketSlice {
		if now.Sub(b.bucketSlice[i].t) > time.Duration(b.windowMS)*time.Millisecond {
			b.bucketSlice = b.bucketSlice[1:]
		} else {
			break
		}
	}
}
