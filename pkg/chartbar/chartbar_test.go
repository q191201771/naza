// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package chartbar

import (
	"fmt"
	"testing"
)

func TestWithItems(t *testing.T) {
	var v []Item

	// 测试中文
	v = []Item{
		//{Name: "China", Num: 1},
		{Name: "中", Num: 22},
		{Name: "中国", Num: 333},
		{Name: "中国啊", Num: 4444},
	}
	fmt.Println(DefaultCtx.WithItems(v))

	// 测试负数
	v = []Item{
		{Name: "q", Num: -10},
		{Name: "w", Num: -8},
		{Name: "e", Num: -4},
		{Name: "r", Num: -1},
		{Name: "t", Num: 0},
		{Name: "y", Num: 2},
		{Name: "u", Num: 3},
		{Name: "i", Num: 6},
		{Name: "o", Num: 10},
	}
	fmt.Println(DefaultCtx.WithItems(v))

	// 测试负数且超出MaxBarLength
	v = []Item{
		{Name: "q", Num: -100},
		{Name: "w", Num: -8},
		{Name: "e", Num: -4},
		{Name: "r", Num: -1},
		{Name: "t", Num: 0},
		{Name: "y", Num: 2},
		{Name: "u", Num: 3},
		{Name: "i", Num: 6},
		{Name: "o", Num: 100},
	}
	fmt.Println(DefaultCtx.WithItems(v))

	// 测试NumLimit范围
	v = []Item{
		{Name: "a", Num: 1},
		{Name: "bb", Num: 3},
		{Name: "ccc", Num: 6},
		{Name: "dddd", Num: 10},
		{Name: "eeeee", Num: 15},
		{Name: "ffffff", Num: 21},
		{Name: "ggggggg", Num: 28},
		{Name: "hhhhhhhh", Num: 36},
		{Name: "jjjjjjjjj", Num: 45},
		{Name: "kkkkkkkkkk", Num: 55},
		{Name: "lllllllllll", Num: 66},
	}
	fmt.Println(NewCtx(func(option *Option) {
		option.PrefixNumLimit = 3
	}).WithItems(v))
	fmt.Println(NewCtx(func(option *Option) {
		option.SuffixNumLimit = 4
	}).WithItems(v))
	fmt.Println(NewCtx(func(option *Option) {
		option.PrefixNumLimit = 3
		option.SuffixNumLimit = 4
	}).WithItems(v))

	// 测试hide
	fmt.Println(NewCtx(func(option *Option) {
		option.HideName = true
	}).WithItems(v))
	fmt.Println(NewCtx(func(option *Option) {
		option.HideNum = true
	}).WithItems(v))
	fmt.Println(NewCtx(func(option *Option) {
		option.HideName = true
		option.HideNum = true
	}).WithItems(v))

	// 测试小的正整数，看0如何绘制
	v = []Item{
		{Name: "a", Num: 0},
		{Name: "bb", Num: 3},
		{Name: "ccc", Num: 6},
		{Name: "dddd", Num: 10},
	}
	fmt.Println(NewCtx().WithItems(v))

	// 测试浮点数
	v = []Item{
		{Name: "bb", Num: 3.3},
		{Name: "ccc", Num: 6.66},
		{Name: "dddd", Num: 10.101},
	}
	fmt.Println(NewCtx().WithItems(v))

	// TODO(chef): 更精细的绘制 "　▏▎▍▌▋▊▉█"
}
