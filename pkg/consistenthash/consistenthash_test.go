// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package consistenthash

import (
	"hash/crc32"
	"math"
	"strconv"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazalog"
)

func TestConsistentHash(t *testing.T) {
	ch := New(1024)
	_, err := ch.Get("aaa")
	assert.Equal(t, ErrIsEmpty, err)

	ch.Add("127.0.0.1")
	ch.Add("0.0.0.0", "8.8.8.8")
	ch.Del("127.0.0.1", "8.8.8.8")
	ch.Add("114.114.114.114", "255.255.255.255", "1.1.1.1", "2.2.2.2", "3.3.3.3")
	exptectedNodes := []string{
		"0.0.0.0",
		"114.114.114.114",
		"255.255.255.255",
		"1.1.1.1",
		"2.2.2.2",
		"3.3.3.3",
	}
	actualNodes := ch.Nodes()
	assert.Equal(t, len(exptectedNodes), len(actualNodes))
	for _, en := range exptectedNodes {
		_, ok := actualNodes[en]
		assert.Equal(t, true, ok)
	}

	counts := make(map[string]int)
	for i := 0; i < 16384; i++ {
		node, err := ch.Get(strconv.Itoa(i))
		assert.Equal(t, nil, err)
		counts[node]++
	}
	nazalog.Debugf("%+v", counts)
}

func TestConsistentHash_Nodes(t *testing.T) {
	nodesGolden := []string{
		"0.0.0.0",
		"114.114.114.114",
		"255.255.255.255",
		"1.1.1.1",
		"2.2.2.2",
		"3.3.3.3",
	}
	j := 1024 // 1
	k := 1024 // 16384
	for i := j; i <= k; i = i << 1 {
		nazalog.Debugf("-----%d-----", i)
		ch := New(i)
		ch.Add(nodesGolden...)
		nodes := ch.Nodes()
		count := uint64(0)
		for k, v := range nodes {
			count += uint64(v)
			nazalog.Debugf("%s: %+v", k, float32(v)/float32(math.MaxUint32+1))
		}
		assert.Equal(t, uint64(math.MaxUint32+1), count)
	}
}

func TestCorner(t *testing.T) {
	ch := New(1, func(option *Option) {
		option.hfn = crc32.ChecksumIEEE
	})

	ch = New(1)
	nodes := ch.Nodes()
	assert.Equal(t, nil, nodes)

	ch = New(1)
	ch.Add("127.0.0.1")
	nodes = ch.Nodes()
	exptectedNodes := map[string]uint64{
		"127.0.0.1": math.MaxUint32 + 1,
	}
	assert.Equal(t, exptectedNodes, nodes)
}
