// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package bitrate_test

import (
	"testing"
	"time"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/bitrate"
)

func TestBitrate(t *testing.T) {
	var b bitrate.Bitrate
	b = bitrate.New(func(option *bitrate.Option) {
		option.WindowMS = 10
	})
	b.Add(1000)
	r := b.Rate()
	assert.Equal(t, float32(800), r)
	time.Sleep(100 * time.Millisecond)
	b.Rate()
}

func TestUnit(t *testing.T) {
	golden := map[bitrate.Unit]float32{
		bitrate.UnitBitPerSec:   800 * 1000,
		bitrate.UnitBytePerSec:  100 * 1000,
		bitrate.UnitKBitPerSec:  800,
		bitrate.UnitKBytePerSec: 100,
	}
	for k, v := range golden {
		b := bitrate.New(func(option *bitrate.Option) {
			option.WindowMS = 10
			option.Unit = k
		})
		b.Add(1000)
		r := b.Rate()
		assert.Equal(t, v, r)
	}
}

func TestOutsizeNow(t *testing.T) {
	var b bitrate.Bitrate
	b = bitrate.New(func(option *bitrate.Option) {
		option.WindowMS = 10
	})
	now := time.Now().UnixNano() / 1e6
	b.Add(1000, now)
	r := b.Rate(now)
	assert.Equal(t, float32(800), r)
}
