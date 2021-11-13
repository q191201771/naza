// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package mock

import (
	"fmt"
	"testing"
	"time"

	"github.com/q191201771/naza/pkg/assert"
)

func TestClock(t *testing.T) {
	var (
		c     Clock
		timer *Timer
		ch    time.Time
		flag  bool
	)

	// 测试Now
	{
		c = NewStdClock()
		fmt.Printf("%+v\n", c.Now())
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("%+v\n", c.Now())

		c = NewFakeClock()
		fmt.Printf("%+v\n", c.Now())
		c.Add(10 * time.Millisecond)
		fmt.Printf("%+v\n", c.Now())
	}

	// 简单测试Timer
	{
		c = NewStdClock()
		fmt.Printf("%+v\n", c.Now())
		timer = c.NewTimer(100 * time.Millisecond)
		ch = <-timer.C
		fmt.Printf("%+v\n", ch)

		c = NewFakeClock()
		fmt.Printf("%+v\n", c.Now())
		timer = c.NewTimer(100 * time.Millisecond)
		c.Add(100 * time.Millisecond)
		ch = <-timer.C
		fmt.Printf("%+v\n", ch)
	}

	// 测试Set
	{
		c = NewFakeClock()
		fmt.Printf("%+v\n", c.Now())
		c.Set(time.Date(2000, 1, 2, 3, 4, 5, 6, time.Local))
		fmt.Printf("%+v\n", c.Now())
		c.Set(time.Now())
		fmt.Printf("%+v\n", c.Now())
	}

	// 测试Timer::Stop
	{
		c = NewStdClock()
		timer = c.NewTimer(100 * time.Millisecond)
		flag = timer.Stop()
		assert.Equal(t, true, flag)

		c = NewFakeClock()
		timer = c.NewTimer(100 * time.Millisecond)
		flag = timer.Stop()
		assert.Equal(t, true, flag)
	}

	// 测试超时后Timer::Stop
	{
		c = NewStdClock()
		timer = c.NewTimer(1 * time.Millisecond)
		time.Sleep(100 * time.Millisecond)
		flag = timer.Stop()
		assert.Equal(t, false, flag)

		c = NewFakeClock()
		timer = c.NewTimer(1 * time.Millisecond)
		c.Add(100 * time.Millisecond)
		flag = timer.Stop()
		assert.Equal(t, false, flag)
	}

	// 测试Timer::Stop后再Timer::Stop
	{
		c = NewStdClock()
		timer = c.NewTimer(100 * time.Millisecond)
		flag = timer.Stop()
		assert.Equal(t, true, flag)
		flag = timer.Stop()
		assert.Equal(t, false, flag)

		c = NewFakeClock()
		timer = c.NewTimer(100 * time.Millisecond)
		flag = timer.Stop()
		assert.Equal(t, true, flag)
		flag = timer.Stop()
		assert.Equal(t, false, flag)
	}

	// 测试Timer::Reset
	{
		c = NewStdClock()
		timer = c.NewTimer(100 * time.Millisecond)
		flag = timer.Reset(1 * time.Millisecond)
		assert.Equal(t, true, flag)
		time.Sleep(100 * time.Millisecond)
		flag = timer.Reset(100 * time.Millisecond)
		assert.Equal(t, false, flag)

		c = NewFakeClock()
		timer = c.NewTimer(100 * time.Millisecond)
		flag = timer.Reset(1 * time.Millisecond)
		assert.Equal(t, true, flag)
		c.Add(100 * time.Millisecond)
		flag = timer.Reset(100 * time.Millisecond)
		assert.Equal(t, false, flag)
	}

	// 测试多次Add，多个Timer
	{
		c = NewFakeClock()
		t1 := c.NewTimer(100 * time.Millisecond)
		t2 := c.NewTimer(100 * time.Millisecond)
		t3 := c.NewTimer(200 * time.Millisecond)
		c.Add(80 * time.Millisecond)
		c.Add(50 * time.Millisecond)
		ch = <-t1.C
		ch = <-t2.C
		c.Add(300 * time.Millisecond)
		ch = <-t3.C
	}

	// corner
	{
		c = NewStdClock()
		c.Add(1 * time.Millisecond)
		c.Set(time.Now())
	}
}
