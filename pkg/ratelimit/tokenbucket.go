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
	"fmt"
	"sync"
	"time"

	"github.com/q191201771/naza/pkg/nazaatomic"
)

var ErrTokenNotEnough = errors.New("naza.ratelimit: token not enough")

// 令牌桶
type TokenBucket struct {
	capacity                  int
	prodTokenInterval         time.Duration
	prodTokenNumEveryInterval int

	disposeFlag nazaatomic.Bool

	mu        sync.Mutex
	available int
	cond      *sync.Cond
}

// @param capacity: 桶容量大小
// @param prodTokenIntervalMSec: 生产令牌的时间间隔，单位毫秒
// @param prodTokenNumEveryInterval: 每次生产多少个令牌
func NewTokenBucket(capacity int, prodTokenIntervalMSec int, prodTokenNumEveryInterval int) *TokenBucket {
	tb := &TokenBucket{
		capacity:                  capacity,
		prodTokenInterval:         time.Duration(time.Duration(prodTokenIntervalMSec) * time.Millisecond),
		prodTokenNumEveryInterval: prodTokenNumEveryInterval,
	}
	tb.cond = sync.NewCond(&tb.mu)
	tb.asyncProdToken()
	return tb
}

func (tb *TokenBucket) TryAquire() error {
	return tb.TryAquireWithNum(1)
}

func (tb *TokenBucket) WaitUntilAquire() {
	tb.WaitUntilAquireWithNum(1)
}

// 尝试获取相应数量的令牌，获取成功返回nil，获取失败返回ErrTokenNotEnough
// 如果获取失败，上层可自由选择多久后重试或丢弃本次任务
func (tb *TokenBucket) TryAquireWithNum(num int) error {
	tb.checkAquireNum(num)

	tb.mu.Lock()
	defer tb.mu.Unlock()
	if tb.available >= num {
		tb.available -= num
		return nil
	}

	return ErrTokenNotEnough
}

// 阻塞直到获取到相应数量的令牌
func (tb *TokenBucket) WaitUntilAquireWithNum(num int) {
	tb.checkAquireNum(num)

	for {
		tb.mu.Lock()
		if tb.available >= num {
			tb.available -= num
			tb.mu.Unlock()
			return
		}

		// 等待下次令牌生产时被唤醒
		// wait的内部会将自身添加到事件监听队列中然后释放锁，当接收到事件时，内部会重新获取锁然后返回
		tb.cond.Wait()
		tb.mu.Unlock()
	}
}

// 销毁令牌桶
func (tb *TokenBucket) Dispose() {
	tb.disposeFlag.Store(true)
}

func (tb *TokenBucket) asyncProdToken() {
	go func() {
		t := time.NewTicker(tb.prodTokenInterval)
		defer t.Stop()
		for {
			if tb.disposeFlag.Load() {
				break
			}
			select {
			case <-t.C:
				tb.mu.Lock()
				tb.available += tb.prodTokenNumEveryInterval
				if tb.available > tb.capacity {
					tb.available = tb.capacity
				}
				// It is allowed but not required for the caller to hold c.L
				// during the call.
				tb.cond.Broadcast()
				tb.mu.Unlock()
			}
		}
	}()
}

func (tb *TokenBucket) checkAquireNum(num int) {
	if num > tb.capacity {
		panic(fmt.Sprintf("aquire num should not bigger than capacity. num=%d, capacity=%d", num, tb.capacity))
	}
}
