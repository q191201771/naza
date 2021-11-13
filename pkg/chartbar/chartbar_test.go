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

	// 测试范围
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

	// TODO(chef): 更精细的绘制 "　▏▎▍▌▋▊▉█"
}
