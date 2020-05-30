// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package lru

import "container/list"

type LRU struct {
	c int                           // capacity
	m map[interface{}]*list.Element // mapping key -> index
	l *list.List                    // value
}

type pair struct {
	k interface{}
	v interface{}
}

func New(capacity int) *LRU {
	return &LRU{
		c: capacity,
		m: make(map[interface{}]*list.Element),
		l: list.New(),
	}
}

// 注意：
// 1. 无论插入前，元素是否已经存在，插入后，元素都会存在于lru容器中
// 2. 插入元素时，也会更新热度（不管插入前元素是否已经存在）
// @return 插入前元素已经存在则返回false
func (lru *LRU) Put(k interface{}, v interface{}) bool {
	var (
		exist bool
		e     *list.Element
	)
	e, exist = lru.m[k]
	if exist {
		lru.l.Remove(e)
		delete(lru.m, k)
	}

	// 头部更热
	e = lru.l.PushFront(pair{k, v})
	lru.m[k] = e

	if lru.l.Len() > lru.c {
		k = lru.l.Back().Value.(pair).k

		lru.l.Remove(lru.l.Back())
		delete(lru.m, k)
	}

	return !exist
}

func (lru *LRU) Get(k interface{}) (v interface{}, exist bool) {
	e, exist := lru.m[k]
	if !exist {
		return nil, false
	}
	pair := e.Value.(pair)
	lru.l.MoveToFront(e)
	return pair.v, true
}

func (lru *LRU) Size() int {
	return lru.l.Len()
}
