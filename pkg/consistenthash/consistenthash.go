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
	"math"
	"sort"
	"strconv"
)

var ErrIsEmpty = errors.New("naza.consistenthash: is empty")

type ConsistentHash interface {
	Add(nodes ...string)
	Del(nodes ...string)
	Get(key string) (node string, err error)

	// @return: 返回的 map 的
	//          key 为添加到内部的 node，
	//          value 为该 node 在环上所占的 point 个数。
	//          我们可以通过各个 node 对应的 point 个数是否接近，来判断各 node 在环上的分布是否均衡。
	//          map 的所有 value 加起来应该等于 (math.MaxUint32 + 1)
	Nodes() map[string]uint64
}

type HashFunc func([]byte) uint32

type Option struct {
	hfn HashFunc
}

var defaultOption = Option{
	hfn: crc32.ChecksumIEEE,
}

type ModOption func(option *Option)

// @param dups: 每个实际的 node 转变成多少个环上的节点，必须大于等于1
// @param modOptions: 可修改内部的哈希函数，比如替换成murmur32的开源实现，可以这样：
//   import "github.com/spaolacci/murmur3"
//   import "github.com/q191201771/naza/pkg/consistenthash"
//
//   ch := consistenthash.New(1000, func(option *Option) {
//     option.hfn = func(bytes []byte) uint32 {
//       h := murmur3.New32()
//       h.Write(bytes)
//       return h.Sum32()
//     }
//   })
func New(dups int, modOptions ...ModOption) ConsistentHash {
	option := defaultOption
	for _, fn := range modOptions {
		fn(&option)
	}

	return &consistentHash{
		point2node: make(map[uint32]string),
		dups:       dups,
		option:     option,
	}
}

type consistentHash struct {
	point2node map[uint32]string
	points     []uint32
	dups       int
	option     Option
}

func (ch *consistentHash) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < ch.dups; i++ {
			point := ch.hash2point(virtualKey(node, i))
			ch.point2node[point] = node
			ch.points = append(ch.points, point)
		}
	}
	sortSlice(ch.points)
}

func (ch *consistentHash) Del(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < ch.dups; i++ {
			point := ch.hash2point(virtualKey(node, i))
			delete(ch.point2node, point)
		}
	}

	ch.points = nil
	for k := range ch.point2node {
		ch.points = append(ch.points, k)
	}
	sortSlice(ch.points)
}

func (ch *consistentHash) Get(key string) (node string, err error) {
	if len(ch.points) == 0 {
		return "", ErrIsEmpty
	}

	point := ch.hash2point(key)
	// 从数组中找出满足 point 值 >= key 所对应 point 值的最小的元素
	index := sort.Search(len(ch.points), func(i int) bool {
		return ch.points[i] >= point
	})

	if index == len(ch.points) {
		index = 0
	}

	return ch.point2node[ch.points[index]], nil
}

func (ch *consistentHash) Nodes() map[string]uint64 {
	if len(ch.points) == 0 {
		return nil
	}
	ret := make(map[string]uint64)
	prev := uint64(0)
	for _, point := range ch.points {
		node := ch.point2node[point]
		ret[node] = ret[node] + uint64(point) - prev
		prev = uint64(point)
	}

	// 最后一个 node 到终点位置的 point 都归入第一个 node
	point := ch.points[len(ch.points)-1]
	node := ch.point2node[point]
	ret[node] = ret[node] + uint64(math.MaxUint32-point+1)
	return ret
}

func (ch *consistentHash) hash2point(key string) uint32 {
	return ch.option.hfn([]byte(key))
}

func virtualKey(node string, index int) string {
	return node + strconv.Itoa(index)
}

func sortSlice(a []uint32) {
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
}
