// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package taskpool

type Worker struct {
	taskChan chan Task
	p        *pool
}

func NewWorker(p *pool) *Worker {
	return &Worker{
		taskChan: make(chan Task, 1),
		p:        p,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			task, ok := <-w.taskChan
			if !ok {
				break
			}
			task()
			w.p.markIdle(w)
		}
	}()
}

func (w *Worker) Stop() {
	close(w.taskChan)
}

func (w *Worker) Go(task Task) {
	w.taskChan <- task
}
