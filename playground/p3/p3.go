// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Origin struct {
	a uint64
	b uint64
}

type WithPadding struct {
	a uint64
	_ [56]byte
	b uint64
	_ [56]byte
}

var num = 1000 * 1000

func OriginParallel() {
	var v Origin

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.a, 1)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.b, 1)
		}
		wg.Done()
	}()

	wg.Wait()
	_ = v.a + v.b
}

func WithPaddingParallel() {
	var v WithPadding

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.a, 1)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.b, 1)
		}
		wg.Done()
	}()

	wg.Wait()
	_ = v.a + v.b
}

func main() {
	var b time.Time

	b = time.Now()
	OriginParallel()
	fmt.Printf("OriginParallel. Cost=%+v.\n", time.Now().Sub(b))

	b = time.Now()
	WithPaddingParallel()
	fmt.Printf("WithPaddingParallel. Cost=%+v.\n", time.Now().Sub(b))
}
