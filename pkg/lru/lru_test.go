// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package lru_test

import (
	"testing"

	"github.com/q191201771/naza/pkg/assert"

	"github.com/q191201771/naza/pkg/lru"
)

func TestLRU(t *testing.T) {
	l := lru.New(3)
	l.Put("chef", 1)
	l.Put("yoko", 2)
	l.Put("tom", 3)
	l.Put("jerry", 4) // 超过容器大小，淘汰最老的`chef`

	v, exist := l.Get("chef")
	assert.Equal(t, false, exist)

	v, exist = l.Get("yoko")
	assert.Equal(t, true, exist)
	assert.Equal(t, 2, v.(int))

	l.Put("garfield", 5) // 超过容器大小，注意，由于`yoko`刚才读取时会更新热度，所以淘汰的是`tom`

	v, exist = l.Get("yoko")
	assert.Equal(t, true, exist)
	assert.Equal(t, 2, v.(int))

	v, exist = l.Get("tom")
	assert.Equal(t, false, exist)

	l = lru.New(3)
	v, exist = l.Get("notexist")
	assert.Equal(t, false, exist)
	assert.Equal(t, 0, l.Size())

	l.Put("chef", 60)
	assert.Equal(t, 1, l.Size())

	v, exist = l.Get("chef")
	assert.Equal(t, true, exist)
	assert.Equal(t, 60, v.(int))
	assert.Equal(t, 1, l.Size())

	v, exist = l.Get("ne")
	assert.Equal(t, false, exist)
	assert.Equal(t, 1, l.Size())

	l.Put("yoko", 100)
	assert.Equal(t, 2, l.Size())

	l.Put("coco", 33)
	assert.Equal(t, 3, l.Size())

	l.Put("dad", 44)
	assert.Equal(t, 3, l.Size())

	isNewPut := l.Put("coco", 1000)
	assert.Equal(t, false, isNewPut)
}
