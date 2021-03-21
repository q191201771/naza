// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package unique

import (
	"sync"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
)

func TestGenUniqueKey(t *testing.T) {
	m := make(map[string]struct{})
	var mutex sync.Mutex
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(j int) {
			var uk string
			if j%2 == 0 {
				uk = GenUniqueKey("hello")
			} else {
				uk = GenUniqueKey("world")
			}
			mutex.Lock()
			m[uk] = struct{}{}
			mutex.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	assert.Equal(t, 1000, len(m))
}

func TestSingleGenerator_GenUniqueKey(t *testing.T) {
	si1 := NewSingleGenerator("one")
	si2 := NewSingleGenerator("two")
	var uk1, uk2 string
	for i := 0; i < 1000; i++ {
		uk1 = si1.GenUniqueKey()
		uk2 = si2.GenUniqueKey()
	}
	assert.Equal(t, "one1000", uk1)
	assert.Equal(t, "two1000", uk2)
}

func BenchmarkGenUniqueKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenUniqueKey("benchmark")
	}
}
