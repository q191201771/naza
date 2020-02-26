// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package ratelimit

import (
	"errors"
	"sync"
	"time"
)

var ErrResourceNotAvailable = errors.New("naza.ratelimit: resource not available")

// 漏桶
type LeakyBucket struct {
	intervalMSec int64

	mu       sync.Mutex
	lastTick int64
}

// @param intervalMSec 多长时间以上，允许获取到一个资源，单位毫秒
func NewLeakyBucket(intervalMSec int) *LeakyBucket {
	return &LeakyBucket{
		intervalMSec: int64(intervalMSec),
		// 注意，第一次获取资源，需要与创建对象时的时间点做比较
		lastTick: time.Now().UnixNano() / 1e6,
	}
}

// 尝试获取资源，获取成功返回nil，获取失败返回ErrResourceNotAvailable
// 如果获取失败，上层可自由选择多久后重试或丢弃本次任务
func (lb *LeakyBucket) TryAquire() error {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	nowMSec := time.Now().UnixNano() / 1e6

	// 距离上次获取成功时间超过了间隔阈值，返回成功
	if nowMSec-lb.lastTick > lb.intervalMSec {
		lb.lastTick = nowMSec
		return nil
	}

	return ErrResourceNotAvailable
}

// 阻塞直到获取到资源
func (lb *LeakyBucket) WaitUntilAquire() {
	lb.mu.Lock()
	nowMSec := time.Now().UnixNano() / 1e6

	diff := nowMSec - lb.lastTick
	if diff > lb.intervalMSec {
		lb.lastTick = nowMSec
		lb.mu.Unlock()
		return
	}

	// 没有达到间隔，我们更新lastTick再出锁，使得其他想获取资源的协程以新的lastTick作为判断条件
	lb.lastTick += lb.intervalMSec
	lb.mu.Unlock()

	// 我们不需要等整个interval间隔，因为可能已经过去了一段时间了，
	// 注意，diff是根据更新前的lastTick计算得到的
	time.Sleep(time.Duration(lb.intervalMSec-diff) * time.Millisecond)
	return
}

// 最快可获取到资源距离当前的时长， 但是不保证获取时一定能抢到
// 返回0，说明可以获取，返回非0，则是对应的时长，单位毫秒
func (lb *LeakyBucket) MaybeAvailableIntervalMSec() int64 {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	nowMSec := time.Now().UnixNano() / 1e6

	if nowMSec-lb.lastTick > lb.intervalMSec {
		return 0
	}

	return lb.lastTick + lb.intervalMSec - nowMSec
}
