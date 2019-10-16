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
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"

	"github.com/q191201771/naza/pkg/bufferpool"
	"github.com/q191201771/naza/pkg/nazalog"
)

var bp bufferpool.BufferPool

//var count int32
var doneCount uint32

var gorutineNum = 1000
var loopNum = 1000
var sleepMSec = 10
var usePool = true

//var usePool = false

func size() int {
	return random(0, 128*1024)

	//return 128 * 1024

	//ss := []int{1000, 2000, 5000}
	//////ss := []int{128, 1024, 4096, 16384}
	//atomic.AddInt32(&count, 1)
	//return ss[count % 3]

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
	atomic.AddUint32(&doneCount, 1)
}

func bufferPoolFunc() {
	size := size()
	buf := bp.Get(size)
	buf.Grow(size)
	bp.Put(buf)
	atomic.AddUint32(&doneCount, 1)
}

func main() {
	cpufd, err := os.Create("/tmp/cpu.prof")
	nazalog.FatalIfErrorNotNil(err)
	pprof.StartCPUProfile(cpufd)
	defer pprof.StopCPUProfile()

	rand.Seed(time.Now().Unix())

	bp = bufferpool.NewBufferPool()

	go func() {
		for {
			nazalog.Debugf("time. done=%d, %+v", atomic.LoadUint32(&doneCount), bp.RetrieveStatus())
			time.Sleep(1 * time.Second)

		}
	}()

	var wg sync.WaitGroup
	wg.Add(gorutineNum * loopNum)
	nazalog.Debug("> loop.")
	for i := 0; i < gorutineNum; i++ {
		go func() {
			if usePool {
				for j := 0; j < loopNum; j++ {
					bufferPoolFunc()
					wg.Done()
				}
			} else {
				for j := 0; j < loopNum; j++ {
					originFunc()
					wg.Done()
				}
			}
		}()
	}
	wg.Wait()
	memfd, err := os.Create("/tmp/mem.prof")
	nazalog.FatalIfErrorNotNil(err)
	pprof.WriteHeapProfile(memfd)
	memfd.Close()
	nazalog.Debugf("%+v", bp.RetrieveStatus())
	nazalog.Debug("< loop.")
}
