// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"

	"github.com/q191201771/naza/pkg/slicebytepool"

	"github.com/q191201771/naza/pkg/nazalog"
)

var bp slicebytepool.SliceBytePool

//var count int32
var doneCount uint32
var tmpSliceByte []byte

var gorutineNum = 1000
var loopNum = 1000
var sleepMSec = time.Duration(10) * time.Millisecond

func size() int {
	return random(1, 128*1024)

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
	size := size()
	buf := make([]byte, size)
	tmpSliceByte = buf
	atomic.AddUint32(&doneCount, 1)
	time.Sleep(sleepMSec)
}

func bufferPoolFunc() {
	size := size()
	buf := bp.Get(size)
	tmpSliceByte = buf
	time.Sleep(sleepMSec)
	bp.Put(buf)
	atomic.AddUint32(&doneCount, 1)
}

func main() {
	strategy := parseFlag()
	nazalog.Debugf("strategy: %d", strategy)

	//cpufd, err := os.Create("/tmp/cpu.prof")
	//nazalog.FatalIfErrorNotNil(err)
	//pprof.StartCPUProfile(cpufd)
	//defer pprof.StopCPUProfile()

	rand.Seed(time.Now().Unix())

	if strategy != 3 {
		bp = slicebytepool.NewSliceBytePool(slicebytepool.Strategy(strategy))
	}

	go func() {
		for {
			if strategy != 3 {
				nazalog.Debugf("time. done=%d, pool=%+v", atomic.LoadUint32(&doneCount), bp.RetrieveStatus())
				time.Sleep(1 * time.Second)
			} else {
				nazalog.Debugf("time. done=%d", atomic.LoadUint32(&doneCount))
				time.Sleep(1 * time.Second)
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(gorutineNum * loopNum)
	nazalog.Debug("> loop.")
	for i := 0; i < gorutineNum; i++ {
		go func() {
			if strategy != 3 {
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
	memfd, err := os.Create(fmt.Sprintf("/tmp/mem%d.prof", strategy))
	nazalog.Assert(nil, err)
	_ = pprof.WriteHeapProfile(memfd)
	_ = memfd.Close()
	nazalog.Debug("> GC.")
	runtime.GC()
	nazalog.Debug("< GC.")
	if strategy != 3 {
		nazalog.Debugf("%+v", bp.RetrieveStatus())
	}
	nazalog.Debug("< loop.")
}

func parseFlag() int {
	strategy := flag.Int("t", 0, "type: 1. multi std pool 2. multi slice pool 3. origin")
	flag.Parse()
	if *strategy < 1 || *strategy > 3 {
		flag.Usage()
		os.Exit(1)
	}
	return *strategy
}
