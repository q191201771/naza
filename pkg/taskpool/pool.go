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
)

type taskWrapper struct {
	taskFn TaskFn
	param  []interface{}
}

type pool struct {
	maxWorkerNum int

	m              sync.Mutex
	totalWorkerNum int
	idleWorkerList []*worker
	blockTaskList  []taskWrapper
}

func newPool(option Option) *pool {
	p := pool{
		maxWorkerNum: option.MaxWorkerNum,
	}
	for i := 0; i < option.InitWorkerNum; i++ {
		p.newWorker()
	}
	return &p
}

func (p *pool) Go(task TaskFn, param ...interface{}) {
	tw := taskWrapper{
		taskFn: task,
		param:  param,
	}
	var w *worker
	p.m.Lock()
	if len(p.idleWorkerList) != 0 {
		// 还有空闲worker

		w = p.idleWorkerList[len(p.idleWorkerList)-1]
		p.idleWorkerList = p.idleWorkerList[0 : len(p.idleWorkerList)-1]
		w.Go(tw)
	} else {
		// 无空闲worker

		if p.maxWorkerNum == 0 ||
			(p.maxWorkerNum > 0 && p.totalWorkerNum < p.maxWorkerNum) {
			// 无最大worker限制，或还未达到限制

			p.newWorkerWithTask(tw)
		} else {
			// 已达到限制

			p.blockTaskList = append(p.blockTaskList, tw)
		}
	}
	p.m.Unlock()
}

func (p *pool) KillIdleWorkers() {
	p.m.Lock()
	p.totalWorkerNum = p.totalWorkerNum - len(p.idleWorkerList)
	for i := range p.idleWorkerList {
		p.idleWorkerList[i].Stop()
	}
	p.idleWorkerList = p.idleWorkerList[0:0]
	p.m.Unlock()
}

func (p *pool) GetCurrentStatus() Status {
	p.m.Lock()
	defer p.m.Unlock()
	return Status{
		TotalWorkerNum: p.totalWorkerNum,
		IdleWorkerNum:  len(p.idleWorkerList),
		BlockTaskNum:   len(p.blockTaskList),
	}
}

func (p *pool) newWorker() *worker {
	w := NewWorker(p)
	w.Start()
	p.idleWorkerList = append(p.idleWorkerList, w)
	p.totalWorkerNum++
	return w
}

func (p *pool) newWorkerWithTask(task taskWrapper) {
	w := NewWorker(p)
	w.Start()
	w.Go(task)
	p.totalWorkerNum++
}

func (p *pool) onIdle(w *worker) {
	p.m.Lock()
	if len(p.blockTaskList) == 0 {
		// 没有等待执行的任务

		p.idleWorkerList = append(p.idleWorkerList, w)
	} else {
		t := p.blockTaskList[0]
		p.blockTaskList = p.blockTaskList[1:]
		w.Go(t)
	}
	p.m.Unlock()
}
