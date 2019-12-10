// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package snowflake_test

import (
	"sync"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazalog"
	"github.com/q191201771/naza/pkg/snowflake"
)

func TestNew(t *testing.T) {
	n, err := snowflake.New(0, 0)
	assert.Equal(t, nil, err)
	id, err := n.Gen()
	assert.Equal(t, nil, err)
	nazalog.Debug(id)
	id, err = n.Gen()
	assert.Equal(t, nil, err)
	nazalog.Debug(id)
}

func TestNegative(t *testing.T) {
	n, err := snowflake.New(0, 0)
	assert.Equal(t, nil, err)
	id, err := n.Gen(1288834974657 + (1 << 41))
	assert.Equal(t, true, id < 0)
}

func TestAlwaysPositive(t *testing.T) {
	n, err := snowflake.New(0, 0, func(option *snowflake.Option) {
		option.AlwaysPositive = true
	})
	assert.Equal(t, nil, err)
	id, err := n.Gen(1288834974657 + (1 << 41))
	assert.Equal(t, true, id == 0)
	id, err = n.Gen(1288834974657 + (1 << 41) + 0x1234)
	assert.Equal(t, true, id >= 0)
}

func TestErrInitial(t *testing.T) {
	var (
		n   *snowflake.Node
		err error
		id  int64
	)

	n, err = snowflake.New(0, 0, func(option *snowflake.Option) {
		option.SequenceBits = 64
	})
	assert.Equal(t, snowflake.ErrInitial, err)

	n, err = snowflake.New(0, 0, func(option *snowflake.Option) {
		option.WorkerIDBits = 64
	})
	assert.Equal(t, snowflake.ErrInitial, err)

	n, err = snowflake.New(0, 0, func(option *snowflake.Option) {
		option.DataCenterIDBits = 64
	})
	assert.Equal(t, snowflake.ErrInitial, err)

	n, err = snowflake.New(0, 0, func(option *snowflake.Option) {
		option.DataCenterIDBits = 31
		option.WorkerIDBits = 31
		option.SequenceBits = 31
	})
	assert.Equal(t, snowflake.ErrInitial, err)

	n, err = snowflake.New(100, 0, func(option *snowflake.Option) {
		option.DataCenterIDBits = 1
	})
	assert.Equal(t, snowflake.ErrInitial, err)

	n, err = snowflake.New(0, 100, func(option *snowflake.Option) {
		option.WorkerIDBits = 1
	})
	assert.Equal(t, snowflake.ErrInitial, err)

	if n != nil {
		id, err = n.Gen()
		assert.Equal(t, nil, err)
		nazalog.Debug(id)
	}
}

func TestErrGen(t *testing.T) {
	var (
		n   *snowflake.Node
		err error
		id  int64
	)

	n, err = snowflake.New(0, 0)
	assert.Equal(t, nil, err)
	assert.IsNotNil(t, n)

	id, err = n.Gen(2)
	assert.Equal(t, nil, err)
	id, err = n.Gen(1)
	assert.Equal(t, snowflake.ErrGen, err)
	nazalog.Debug(id)
}

func TestMT(t *testing.T) {
	var (
		n   *snowflake.Node
		err error
	)

	n, err = snowflake.New(0, 0, func(option *snowflake.Option) {
		option.SequenceBits = 1
	})
	assert.Equal(t, nil, err)
	assert.IsNotNil(t, n)

	ii := 16
	jj := 16
	var mu sync.Mutex
	m := make(map[int64]struct{})
	var wg sync.WaitGroup
	wg.Add(ii * jj)

	for i := 0; i < ii; i++ {
		go func() {
			for j := 0; j < jj; j++ {
				id, err := n.Gen()
				assert.Equal(t, nil, err)
				mu.Lock()
				m[id] = struct{}{}
				mu.Unlock()
				wg.Done()
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, ii*jj, len(m))
}

func BenchmarkNode_Gen(b *testing.B) {
	var (
		n   *snowflake.Node
		err error
	)

	n, err = snowflake.New(0, 0)
	assert.Equal(b, nil, err)
	assert.IsNotNil(b, n)

	var dummy int64
	for i := 0; i < b.N; i++ {
		id, _ := n.Gen()
		dummy += id
	}
	nazalog.Debug(dummy)
}
