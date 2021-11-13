// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// Package chartbar 控制台绘制ascii柱状图
//
package chartbar

// TODO(chef): 如果总数小于绘制长度并且都是正整数，可以考虑按原始值而非比例绘制

const (
	OrderOrigin    Order = iota + 1 // 原始序
	OrderAscCount                   // 按计数值升序排序
	OrderDescCount                  // 按计数值降序排序
	OrderAscName                    // 按字段名称升序排序
	OrderDescName                   // 按字段名称降序排序

	NoNumLimit = -1
)

type Item struct {
	Name string  // key
	Num  float64 // value

	count int // bar
}

type Order int

type Option struct {
	MaxBarLength    int
	DrawIconBlock   string
	DrawIconPadding string

	Order          Order
	PrefixNumLimit int
	SuffixNumLimit int
}

var defaultOption = Option{
	// 50 "▇" " "
	// 18 "口" "　"
	MaxBarLength:    50,  // MaxBarLength 柱状图形的最大长度
	DrawIconBlock:   "▇", // 柱状图实体绘制内容
	DrawIconPadding: " ", // 柱状图空余部分绘制内容

	Order:          OrderDescCount,
	PrefixNumLimit: NoNumLimit,
	SuffixNumLimit: NoNumLimit,
}

// ---------------------------------------------------------------------------------------------------------------------

var DefaultCtx = NewCtx()

type ModOption func(option *Option)

func NewCtx(modOptions ...ModOption) *Ctx {
	option := defaultOption
	for _, fn := range modOptions {
		fn(&option)
	}
	return &Ctx{
		option: option,
	}
}

func NewCtxWith(ctx *Ctx, modOptions ...ModOption) *Ctx {
	option := ctx.option
	for _, fn := range modOptions {
		fn(&option)
	}
	return &Ctx{
		option: option,
	}
}
