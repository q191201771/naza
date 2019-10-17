// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// TODO
// - channel 通信替换成其他方式是否有可能提高性能

package taskpool

import (
	"errors"
)

var ErrTaskPool = errors.New("naza.taskpool: fxxk")

type Task func()

type Pool interface {
	// 向池内放入任务
	Go(task Task)

	// 获取当前空闲worker和工作worker的数量。注意，这只是一个瞬时值
	Status() (idleWorkerNum int, busyWorkerNum int)

	// 关闭池内所有的空闲worker(协程)
	KillIdleWorkers()
}

type Option struct {
	InitWorkerNum int // 创建池对象时，预先开启的worker(协程)数量，如果为0，则不预先开启
}

var defaultOption = Option{
	InitWorkerNum: 0,
}

type ModOption func(option *Option)

func NewPool(modOptions ...ModOption) (Pool, error) {
	option := defaultOption

	for _, fn := range modOptions {
		fn(&option)
	}

	if err := validate(option); err != nil {
		return nil, err
	}

	var p pool
	for i := 0; i < option.InitWorkerNum; i++ {
		w := NewWorker(&p)
		w.Start()
		p.idleWorkerList = append(p.idleWorkerList, w)
	}
	p.idleWorkerNum.Add(uint32(option.InitWorkerNum))
	return &p, nil
}

func validate(option Option) error {
	if option.InitWorkerNum < 0 {
		return ErrTaskPool
	}
	return nil
}
