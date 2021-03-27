// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazasync

import "testing"

func TestMutex(t *testing.T) {
	var mu Mutex
	mu.Lock()
	mu.Unlock()
	ch := make(chan struct{}, 1)
	go func() {
		mu.Lock()
		mu.Unlock()
		ch <- struct{}{}
	}()
	<-ch
}

func TestMutex_Corner(t *testing.T) {
	//var mu Mutex
	// case1 递归
	//mu.Lock()
	//mu.Lock()

	// case2 先UnLock
	//mu.Unlock()

	// case3 协程1Lock，协程2Unlock
	//mu.Lock()
	//ch := make(chan struct{}, 1)
	//ch <- struct{}{}
	//go func() {
	//	<-ch
	//	mu.Unlock()
	//}()
	//time.Sleep(200 * time.Millisecond)
}
