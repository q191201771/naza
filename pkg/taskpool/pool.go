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

	"github.com/q191201771/naza/pkg/nazaatomic"
)

type pool struct {
	idleWorkerNum nazaatomic.Uint32
	busyWorkerNum nazaatomic.Uint32

	m              sync.Mutex
	idleWorkerList []*Worker
}

func (p *pool) Go(task Task) {
	var w *Worker
	p.m.Lock()
	if len(p.idleWorkerList) != 0 {
		w = p.idleWorkerList[len(p.idleWorkerList)-1]
		p.idleWorkerList = p.idleWorkerList[0 : len(p.idleWorkerList)-1]
		p.idleWorkerNum.Decrement()
		p.busyWorkerNum.Increment()
	}
	p.m.Unlock()
	if w == nil {
		w = NewWorker(p)
		w.Start()
		p.busyWorkerNum.Increment()
	}
	w.Go(task)
}

func (p *pool) KillIdleWorkers() {
	p.m.Lock()
	for i := range p.idleWorkerList {
		p.idleWorkerList[i].Stop()
	}
	p.idleWorkerNum.Sub(uint32(len(p.idleWorkerList)))
	p.idleWorkerList = p.idleWorkerList[0:0]
	p.m.Unlock()
}

func (p *pool) Status() (idleWorkerNum int, busyWorkerNum int) {
	return int(p.idleWorkerNum.Load()), int(p.busyWorkerNum.Load())
}

func (p *pool) markIdle(w *Worker) {
	p.m.Lock()
	p.idleWorkerNum.Increment()
	p.busyWorkerNum.Decrement()
	p.idleWorkerList = append(p.idleWorkerList, w)
	p.m.Unlock()
}
