// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"bytes"
	"github.com/q191201771/naza/pkg/bufferpool"
	"github.com/q191201771/naza/pkg/nazalog"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

var runNum = 1000 * 1000
var bp bufferpool.BufferPool
var count int32

func size() int {
	//return 128 * 1024

	//ss := []int{1000, 2000, 5000}
	//////ss := []int{128, 1024, 4096, 16384}
	//atomic.AddInt32(&count, 1)
	//return ss[count % 3]

	return random(0, 4*1024)

	//count++
	//if count > 128 * 1024 {
	//	count = 1
	//}
	//return count
}

func random(l, r int) int {
	return l + (rand.Int() % (r - l))
}

func originFunc() {
	var buf bytes.Buffer
	size := size()
	buf.Grow(size)
	time.Sleep(time.Duration(random(0, 200)) * time.Millisecond)
}

func bufferPoolFunc() {
	size := size()
	buf := bp.Get(size)
	buf.Grow(size)
	time.Sleep(time.Duration(random(0, 200)) * time.Millisecond)
	bp.Put(buf)
}

func main() {
	f, _ := os.Create("/tmp/demo.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	rand.Seed(time.Now().Unix())

	bp = bufferpool.NewBufferPool()

	go func() {
		for {
			nazalog.Debugf("time. %+v", bp.RetrieveStatus())
			time.Sleep(1 * time.Second)

		}
	}()

	var wg sync.WaitGroup
	wg.Add(runNum)
	nazalog.Debug("> loop.")
	for i := 0; i < runNum; i++ {
		go func() {
			//originFunc()
			bufferPoolFunc()
			wg.Done()
		}()
	}
	wg.Wait()
	nazalog.Debugf("%+v", bp.RetrieveStatus())
	nazalog.Debug("< loop.")
}
