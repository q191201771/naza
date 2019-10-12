// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package taskpool

import (
	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazalog"
	"sync"
	"testing"
	"time"
)

var (
	taskNum       = 1000 * 1000
	initWorkerNum = 1000
)

func BenchmarkOriginGo(b *testing.B) {
	var wg sync.WaitGroup
	for j := 0; j < b.N; j++ {
		wg.Add(taskNum)
		for i := 0; i < taskNum; i++ {
			go func() {
				time.Sleep(1 * time.Millisecond)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkTaskPool(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := NewPool(func(option *Option) {
		option.InitWorkerNum = initWorkerNum
	})

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		wg.Add(taskNum)
		for i := 0; i < taskNum; i++ {
			p.Go(func() {
				time.Sleep(1 * time.Millisecond)
				wg.Done()
			})
		}
		wg.Wait()
		idle, busy := p.Status()
		nazalog.Debugf("done, worker num. idle=%d, busy=%d", idle, busy) // 此时还有个别busy也是正常的，因为只是业务方的任务代码执行完了，可能还没回收到idle队列中
		p.KillIdleWorkers()
		idle, busy = p.Status()
		nazalog.Debugf("killed, worker num. idle=%d, busy=%d", idle, busy)
	}
}

func TestTaskPool(t *testing.T) {
	var wg sync.WaitGroup
	p, _ := NewPool(func(option *Option) {
		option.InitWorkerNum = 1
	})

	go func() {
		idle, busy := p.Status()
		nazalog.Debugf("timer, worker num. idle=%d, busy=%d", idle, busy)
		time.Sleep(10 * time.Millisecond)
	}()

	n := 1000
	wg.Add(n)
	for i := 0; i < n; i++ {
		p.Go(func() {
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		})
	}
	wg.Wait()
	idle, busy := p.Status()
	nazalog.Debugf("done, worker num. idle=%d, busy=%d", idle, busy) // 此时还有个别busy也是正常的，因为只是业务方的任务代码执行完了，可能还没回收到idle队列中
	p.KillIdleWorkers()
	idle, busy = p.Status()
	nazalog.Debugf("killed, worker num. idle=%d, busy=%d", idle, busy)
}

func TestGlobal(t *testing.T) {
	err := Init()
	assert.Equal(t, nil, err)
	i, b := Status()
	assert.Equal(t, 0, i)
	assert.Equal(t, 0, b)
	Go(func() {
	})
	KillIdleWorkers()
}

func TestCorner(t *testing.T) {
	_, err := NewPool(func(option *Option) {
		option.InitWorkerNum = -1
	})
	assert.Equal(t, ErrTaskPool, err)
}
