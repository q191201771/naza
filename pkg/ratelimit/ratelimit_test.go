// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package ratelimit_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/q191201771/naza/pkg/nazalog"
	"github.com/q191201771/naza/pkg/ratelimit"
)

var (
	duration = 100 * time.Millisecond
	num      = 50
)

func TestNew(t *testing.T) {
	rl := ratelimit.New(num)
	rl = ratelimit.New(num, func(option *ratelimit.Option) {
		option.Duration = duration
	})
	rl.Wait()
}

func TestRateLimit_Wait(t *testing.T) {
	rl := ratelimit.New(num, func(option *ratelimit.Option) {
		option.Duration = duration
	})
	b := time.Now()
	for i := 0; i < num; i++ {
		rl.Wait()
	}
	nazalog.Debugf("cost:%v", time.Now().Sub(b))
}

func TestRateLimit_Wait2(t *testing.T) {
	rl := ratelimit.New(num, func(option *ratelimit.Option) {
		option.Duration = duration
	})
	b := time.Now()
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func(ii int) {
			rl.Wait()
			wg.Done()
		}(i)
	}
	wg.Wait()
	nazalog.Debugf("cost:%v", time.Now().Sub(b))
}

func TestRateLimit_Wait3(t *testing.T) {
	rand.Seed(time.Now().Unix())
	rl := ratelimit.New(num, func(option *ratelimit.Option) {
		option.Duration = duration
	})
	b := time.Now()
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func(ii int) {
			time.Sleep(time.Duration(rand.Int63n(int64(duration) / 2)))
			rl.Wait()
			wg.Done()
		}(i)
	}
	wg.Wait()
	nazalog.Debugf("cost:%v", time.Now().Sub(b))
}

func TestRateLimit_Wait4(t *testing.T) {
	rl := ratelimit.New(num, func(option *ratelimit.Option) {
		option.Duration = duration
	})
	b := time.Now()
	for i := 0; i < num; i++ {
		time.Sleep(time.Duration(int64(duration) * 2 / int64(num)))
		rl.Wait()
	}
	nazalog.Debugf("cost:%v", time.Now().Sub(b))
}

func BenchmarkRateLimit_Wait(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rl := ratelimit.New(num)
		for i := 0; i < num; i++ {
			rl.Wait()
		}
	}
}
