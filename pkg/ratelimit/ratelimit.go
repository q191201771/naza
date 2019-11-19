// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package ratelimit

import (
	"sync"
	"time"
)

type RateLimit struct {
	num    int
	option Option

	neededWait time.Duration
	mu         sync.Mutex
	last       time.Time
}

type Option struct {
	Duration time.Duration
}

var defaultOption = Option{
	Duration: 1 * time.Second,
}

type ModOption func(option *Option)

func New(num int, modOptions ...ModOption) *RateLimit {
	option := defaultOption
	for _, fn := range modOptions {
		fn(&option)
	}
	return &RateLimit{
		num:        num,
		option:     option,
		neededWait: option.Duration / time.Duration(num),
	}
}

func (rl *RateLimit) Wait() {
	rl.mu.Lock()
	now := time.Now()
	if rl.last.IsZero() {
		rl.last = now
		rl.mu.Unlock()
		return
	}

	diff := now.Sub(rl.last)
	if diff > rl.neededWait {
		rl.last = now
		rl.mu.Unlock()
		return
	}

	rl.last = rl.last.Add(rl.neededWait)
	rl.mu.Unlock()

	t := time.NewTimer(rl.neededWait - diff)
	<-t.C
	t.Stop()
	return
}
