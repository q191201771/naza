// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// 非阻塞协程池，协程数量可动态增长，可配置最大协程并发数量，可手动释放空闲的协程
package taskpool

import (
	"errors"
)

// TODO
// - channel 通信替换成其他方式是否有可能提高性能

var ErrTaskPool = errors.New("naza.taskpool: fxxk")

type TaskFn func(param ...interface{})

type Status struct {
	TotalWorkerNum int // 总协程数量
	IdleWorkerNum  int // 空闲协程数量
	BlockTaskNum   int // 等待执行的任务数。注意，只在协程数量有最大限制的情况下，该值才可能不为0，具体见Option.MaxWorkerNum
}

type Pool interface {
	// 向池内放入任务
	// 非阻塞函数，不会等待task执行
	Go(task TaskFn, param ...interface{})

	// 获取当前的状态，注意，只是一个瞬时值
	GetCurrentStatus() Status

	// 关闭池内所有的空闲协程
	KillIdleWorkers()
}

type Option struct {
	// 创建池对象时，预先开启的worker(协程)数量，如果为0，则不预先开启。只是一个小优化
	InitWorkerNum int

	// - 如果为0，则无协程数量限制。向池中添加任务时如果无空闲协程，会无条件创建新的协程。
	// - 如果不为0，则池内总协程数量达到阈值后，将不再创建新的协程。此时任务会被缓存，等待有空闲协程时才被执行。
	//   可用来控制任务的最大并发数
	MaxWorkerNum int
}

var defaultOption = Option{
	InitWorkerNum: 0,
	MaxWorkerNum:  0,
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

	return newPool(option), nil
}

func validate(option Option) error {
	if option.InitWorkerNum < 0 {
		return ErrTaskPool
	}
	if option.MaxWorkerNum < 0 {
		return ErrTaskPool
	}
	if option.MaxWorkerNum > 0 && option.InitWorkerNum > option.MaxWorkerNum {
		return ErrTaskPool
	}
	return nil
}
