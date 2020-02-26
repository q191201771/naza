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

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazalog"
	"github.com/q191201771/naza/pkg/ratelimit"
)

func TestNewTokenBucket(t *testing.T) {
	ratelimit.NewTokenBucket(2000, 1000, 1000)
}

func TestTokenBucket_TryAquire(t *testing.T) {
	var (
		tb  *ratelimit.TokenBucket
		err error
	)

	tb = ratelimit.NewTokenBucket(2000, 1000, 1000)
	err = tb.TryAquire()
	assert.Equal(t, ratelimit.ErrTokenNotEnough, err)
	err = tb.TryAquire()
	assert.Equal(t, ratelimit.ErrTokenNotEnough, err)

	tb = ratelimit.NewTokenBucket(2000, 1, 1000)
	time.Sleep(10 * time.Millisecond)
	err = tb.TryAquire()
	assert.Equal(t, nil, err)
	err = tb.TryAquire()
	assert.Equal(t, nil, err)
}

func TestTokenBucket_WaitUntilAquire(t *testing.T) {
	var tb *ratelimit.TokenBucket

	tb = ratelimit.NewTokenBucket(2000, 1000, 1000)
	tb.WaitUntilAquire()
	tb.WaitUntilAquire()
}

func TestTokenBucket_Dispose(t *testing.T) {
	var (
		tb  *ratelimit.TokenBucket
		err error
	)

	tb = ratelimit.NewTokenBucket(2000, 1, 1000)
	time.Sleep(10 * time.Millisecond)
	err = tb.TryAquireWithNum(1)
	assert.Equal(t, nil, err)
	tb.WaitUntilAquireWithNum(1)
	tb.Dispose()
}

func TestTokenBucket_panic(t *testing.T) {
	defer func() {
		nazalog.Debug(recover())
	}()
	tb := ratelimit.NewTokenBucket(1, 1, 1)
	tb.TryAquireWithNum(100)
}
