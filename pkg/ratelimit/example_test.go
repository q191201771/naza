// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package ratelimit_test

import (
	"testing"
	"time"

	"github.com/q191201771/naza/pkg/nazalog"
	"github.com/q191201771/naza/pkg/ratelimit"
)

var runExample = true

func TestLeakyBucket(t *testing.T) {
	if !runExample {
		return
	}

	lb := ratelimit.NewLeakyBucket(100)
	for i := 0; i < 16; i++ {
		go func(j int) {
			for k := 0; k < 16; k++ {
				err := lb.TryAquire()
				if err == nil {
					nazalog.Debugf("TryAquire succ. goroutine=%d, index=%d", j, k)
				} else {
					//nazalog.Debugf("TryAquire fail. goroutine=%d, index=%d, err=%v", j, k, err)
					time.Sleep(time.Duration(j*k) * time.Millisecond)
				}
				//time.Sleep(time.Duration(j * k) * time.Millisecond)
			}
		}(i)
	}
	time.Sleep(2 * time.Second)
	//2020/02/25 23:39:13.896802 DEBUG TryAquire succ. goroutine=10, index=5 - example_test.go:23
	//2020/02/25 23:39:13.999801 DEBUG TryAquire succ. goroutine=7, index=8 - example_test.go:23
	//2020/02/25 23:39:14.113313 DEBUG TryAquire succ. goroutine=11, index=8 - example_test.go:23
	//2020/02/25 23:39:14.216280 DEBUG TryAquire succ. goroutine=11, index=10 - example_test.go:23
	//2020/02/25 23:39:14.334289 DEBUG TryAquire succ. goroutine=8, index=12 - example_test.go:23
	//2020/02/25 23:39:14.439413 DEBUG TryAquire succ. goroutine=8, index=14 - example_test.go:23
	//2020/02/25 23:39:14.574673 DEBUG TryAquire succ. goroutine=14, index=11 - example_test.go:23
	//2020/02/25 23:39:14.742899 DEBUG TryAquire succ. goroutine=14, index=13 - example_test.go:23
	//2020/02/25 23:39:14.914085 DEBUG TryAquire succ. goroutine=12, index=14 - example_test.go:23
	//2020/02/25 23:39:15.172121 DEBUG TryAquire succ. goroutine=15, index=14 - example_test.go:23

	lb2 := ratelimit.NewLeakyBucket(100)
	for i := 0; i < 4; i++ {
		go func(j int) {
			for k := 0; k < 4; k++ {
				lb2.WaitUntilAquire()
				nazalog.Debugf("< lb.WaitUntilAquire. goroutine=%d, index=%d", j, k)
			}
		}(i)
	}
	time.Sleep(2 * time.Second)
	//2020/02/25 23:40:11.275685 DEBUG < lb.WaitUntilAquire. goroutine=3, index=0 - example_test.go:49
	//2020/02/25 23:40:11.374789 DEBUG < lb.WaitUntilAquire. goroutine=0, index=0 - example_test.go:49
	//2020/02/25 23:40:11.473445 DEBUG < lb.WaitUntilAquire. goroutine=1, index=0 - example_test.go:49
	//2020/02/25 23:40:11.571714 DEBUG < lb.WaitUntilAquire. goroutine=2, index=0 - example_test.go:49
	//2020/02/25 23:40:11.670913 DEBUG < lb.WaitUntilAquire. goroutine=3, index=1 - example_test.go:49
	//2020/02/25 23:40:11.772003 DEBUG < lb.WaitUntilAquire. goroutine=0, index=1 - example_test.go:49
	//2020/02/25 23:40:11.871239 DEBUG < lb.WaitUntilAquire. goroutine=1, index=1 - example_test.go:49
	//2020/02/25 23:40:11.973307 DEBUG < lb.WaitUntilAquire. goroutine=2, index=1 - example_test.go:49
	//2020/02/25 23:40:12.075015 DEBUG < lb.WaitUntilAquire. goroutine=3, index=2 - example_test.go:49
	//2020/02/25 23:40:12.173357 DEBUG < lb.WaitUntilAquire. goroutine=0, index=2 - example_test.go:49
	//2020/02/25 23:40:12.270387 DEBUG < lb.WaitUntilAquire. goroutine=1, index=2 - example_test.go:49
	//2020/02/25 23:40:12.370509 DEBUG < lb.WaitUntilAquire. goroutine=2, index=2 - example_test.go:49
	//2020/02/25 23:40:12.475001 DEBUG < lb.WaitUntilAquire. goroutine=3, index=3 - example_test.go:49
	//2020/02/25 23:40:12.571062 DEBUG < lb.WaitUntilAquire. goroutine=0, index=3 - example_test.go:49
	//2020/02/25 18:40:12.672385 DEBUG < lb.WaitUntilAquire. goroutine=1, index=3 - example_test.go:49
	//2020/02/25 18:40:12.770939 DEBUG < lb.WaitUntilAquire. goroutine=2, index=3 - example_test.go:49
}

func TestTokenBucket(t *testing.T) {
	if !runExample {
		return
	}

	tb := ratelimit.NewTokenBucket(100, 1000, 10)
	for i := 0; i < 4; i++ {
		go func(j int) {
			for k := 0; k < 4; k++ {
				tb.WaitUntilAquire()
				nazalog.Debugf("< tb.WaitUntilAquire. goroutine=%d, index=%d", j, k)
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}
	time.Sleep(2 * time.Second)
	//2020/02/25 19:02:33.453207 DEBUG < tb.WaitUntilAquire. goroutine=2, index=0 - example_test.go:82
	//2020/02/25 19:02:33.453302 DEBUG < tb.WaitUntilAquire. goroutine=1, index=0 - example_test.go:82
	//2020/02/25 19:02:33.453414 DEBUG < tb.WaitUntilAquire. goroutine=0, index=0 - example_test.go:82
	//2020/02/25 19:02:33.453535 DEBUG < tb.WaitUntilAquire. goroutine=3, index=0 - example_test.go:82
	//2020/02/25 19:02:33.557602 DEBUG < tb.WaitUntilAquire. goroutine=2, index=1 - example_test.go:82
	//2020/02/25 19:02:33.557881 DEBUG < tb.WaitUntilAquire. goroutine=0, index=1 - example_test.go:82
	//2020/02/25 19:02:33.557968 DEBUG < tb.WaitUntilAquire. goroutine=1, index=1 - example_test.go:82
	//2020/02/25 19:02:33.558043 DEBUG < tb.WaitUntilAquire. goroutine=3, index=1 - example_test.go:82
	//2020/02/25 19:02:33.661411 DEBUG < tb.WaitUntilAquire. goroutine=0, index=2 - example_test.go:82
	//2020/02/25 19:02:33.661495 DEBUG < tb.WaitUntilAquire. goroutine=2, index=2 - example_test.go:82
	//2020/02/25 19:02:34.451624 DEBUG < tb.WaitUntilAquire. goroutine=2, index=3 - example_test.go:82
	//2020/02/25 19:02:34.451715 DEBUG < tb.WaitUntilAquire. goroutine=0, index=3 - example_test.go:82
	//2020/02/25 19:02:34.451787 DEBUG < tb.WaitUntilAquire. goroutine=3, index=2 - example_test.go:82
	//2020/02/25 19:02:34.451841 DEBUG < tb.WaitUntilAquire. goroutine=1, index=2 - example_test.go:82
	//2020/02/25 19:02:34.551910 DEBUG < tb.WaitUntilAquire. goroutine=3, index=3 - example_test.go:82
	//2020/02/25 19:02:34.556604 DEBUG < tb.WaitUntilAquire. goroutine=1, index=3 - example_test.go:82

}

func TestRateLimiter(t *testing.T) {
	tb := ratelimit.NewTokenBucket(1, 1, 1)
	lb := ratelimit.NewLeakyBucket(1)
	var rl ratelimit.RateLimiter
	rl = tb
	rl.WaitUntilAquire()
	rl = lb
	rl.TryAquire()
}
