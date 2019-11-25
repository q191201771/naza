// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package consistenthash

import (
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
	exptectedNodes := map[string]struct{}{
		"0.0.0.0":         {},
		"114.114.114.114": {},
		"255.255.255.255": {},
		"1.1.1.1":         {},
		"2.2.2.2":         {},
		"3.3.3.3":         {},
	}
	actualNodes := ch.Nodes()
	assert.Equal(t, exptectedNodes, actualNodes)

	counts := make(map[string]int)
	for i := 0; i < 16384; i++ {
		node, err := ch.Get(strconv.Itoa(i))
		assert.Equal(t, nil, err)
		counts[node]++
	}
	nazalog.Debugf("%+v", counts)
}
