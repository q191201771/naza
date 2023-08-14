// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// Package chartbar 控制台绘制ascii柱状图
package chartbar

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
	MaxBarLength    int    // 柱状图形的最大长度
	DrawIconBlock   string // 柱状图实体绘制内容
	DrawIconPadding string // 柱状图空余部分绘制内容
	HideName        bool   // 是否隐藏图形旁边的Name字段
	HideNum         bool   // 是否隐藏图形旁边的Num字段

	Order          Order // 排序方式
	PrefixNumLimit int   // 只显示前`PrefixNumLimit`个元素，注意，可以和`SuffixNumLimit`同时使用
	SuffixNumLimit int   // 只显示后`SuffixNumLimit`个元素
}

var defaultOption = Option{
	// 50 "▇" " "
	// 18 "口" "　"
	MaxBarLength:    50,
	DrawIconBlock:   "▇",
	DrawIconPadding: " ",
	HideName:        false, //
	HideNum:         false,

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

// NewCtxWith
//
// 在`ctx`参数基础上使用`modOptions`生成新的 Ctx
func NewCtxWith(ctx *Ctx, modOptions ...ModOption) *Ctx {
	option := ctx.option
	for _, fn := range modOptions {
		fn(&option)
	}
	return &Ctx{
		option: option,
	}
}
