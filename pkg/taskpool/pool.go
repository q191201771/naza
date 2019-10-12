// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package taskpool

import (
	"sync"
	"sync/atomic"
)

type pool struct {
	idleWorkerNum int32
	busyWorkerNum int32

	m              sync.Mutex
	idleWorkerList []*Worker
}

func (p *pool) Go(task Task) {
	var w *Worker
	p.m.Lock()
	if len(p.idleWorkerList) != 0 {
		w = p.idleWorkerList[len(p.idleWorkerList)-1]
		p.idleWorkerList = p.idleWorkerList[0 : len(p.idleWorkerList)-1]
		atomic.AddInt32(&p.idleWorkerNum, -1)
		atomic.AddInt32(&p.busyWorkerNum, 1)
	}
	p.m.Unlock()
	if w == nil {
		w = NewWorker(p)
		w.Start()
		atomic.AddInt32(&p.busyWorkerNum, 1)
	}
	w.Go(task)
}

func (p *pool) KillIdleWorkers() {
	p.m.Lock()
	for i := range p.idleWorkerList {
		p.idleWorkerList[i].Stop()
	}
	atomic.AddInt32(&p.idleWorkerNum, int32(-len(p.idleWorkerList)))
	p.idleWorkerList = p.idleWorkerList[0:0]
	p.m.Unlock()
}

func (p *pool) Status() (idleWorkerNum int, busyWorkerNum int) {
	idleWorkerNum = int(atomic.LoadInt32(&p.idleWorkerNum))
	busyWorkerNum = int(atomic.LoadInt32(&p.busyWorkerNum))
	return
}

func (p *pool) markIdle(w *Worker) {
	p.m.Lock()
	atomic.AddInt32(&p.idleWorkerNum, 1)
	atomic.AddInt32(&p.busyWorkerNum, -1)
	p.idleWorkerList = append(p.idleWorkerList, w)
	p.m.Unlock()
}
