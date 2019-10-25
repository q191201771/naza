// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package consistenthash

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
)

type ConsistentHash interface {
	Add(nodes ...string)
	Del(nodes ...string)
	Get(key string) (node string, err error)
	Nodes() map[string]struct{}
}

var ErrIsEmpty = errors.New("naza.consistenthash: is empty")

func New(dups int) ConsistentHash {
	return &consistentHash{
		point2node: make(map[int]string),
		dups:       dups,
	}
}

type consistentHash struct {
	point2node map[int]string
	points     []int
	dups       int
}

func (ch *consistentHash) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < ch.dups; i++ {
			point := hash2point(virtualKey(node, i))
			ch.point2node[point] = node
			ch.points = append(ch.points, point)
		}
	}
	sort.Ints(ch.points)
}

func (ch *consistentHash) Del(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < ch.dups; i++ {
			point := hash2point(virtualKey(node, i))
			delete(ch.point2node, point)
		}
	}

	ch.points = nil
	for k := range ch.point2node {
		ch.points = append(ch.points, k)
	}
	sort.Ints(ch.points)
}

func (ch *consistentHash) Get(key string) (node string, err error) {
	if len(ch.points) == 0 {
		return "", ErrIsEmpty
	}

	point := hash2point(key)
	index := sort.Search(len(ch.points), func(i int) bool {
		return ch.points[i] >= point
	})

	if index == len(ch.points) {
		index = 0
	}

	return ch.point2node[ch.points[index]], nil
}

func (ch *consistentHash) Nodes() map[string]struct{} {
	ret := make(map[string]struct{})
	for _, v := range ch.point2node {
		ret[v] = struct{}{}
	}
	return ret
}

func virtualKey(node string, index int) string {
	return node + strconv.Itoa(index)
}

func hash2point(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}
