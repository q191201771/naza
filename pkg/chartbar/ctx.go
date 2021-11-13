// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package chartbar

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/q191201771/naza/pkg/dataops"
)

type Ctx struct {
	option Option
}

//func (ctx *Ctx) ModOptions(modOptions ...ModOption) *Ctx {
//	for _, fn := range modOptions {
//		fn(&ctx.option)
//	}
//	return ctx
//}

// WithItems
//
// @param items: 注意，内部不会修改切片底层数据的值以及顺序
//
func (ctx *Ctx) WithItems(items []Item) string {
	// 拷贝一份，避免修改外部切片的原始顺序
	if ctx.option.Order != OrderOrigin {
		clone := make([]Item, len(items))
		copy(clone, items)
		items = clone
	}

	switch ctx.option.Order {
	case OrderOrigin:
	// noop
	case OrderAscCount:
		sort.Slice(items, func(i, j int) bool {
			return items[i].Num < items[j].Num
		})
	case OrderDescCount:
		sort.Slice(items, func(i, j int) bool {
			return items[i].Num > items[j].Num
		})
	case OrderAscName:
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name < items[j].Name
		})
	case OrderDescName:
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name > items[j].Name
		})
	}

	var newItems []Item
	dataops.SliceLimit(items, ctx.option.PrefixNumLimit, ctx.option.SuffixNumLimit, func(index int) {
		newItems = append(newItems, items[index])
	})
	items = newItems

	var (
		maxCountLength int // count柱状最长画多长
		maxLengthOfNum int // num字段多长
	)
	minNum, maxNum := calcMinMaxNum(items)
	if minNum > 0.00 {
		// 都是正数的情况，最大的画满柱状条，其他的按与最大占比画
		for i := range items {
			// round四舍五入
			items[i].count = int(math.Round(items[i].Num * float64(ctx.option.MaxBarLength) / maxNum))
			// 最小可能和最大的比太小了
			if items[i].count == 0 {
				items[i].count = 1
			}
		}
		maxCountLength = calcMaxCount(items)
		maxLengthOfNum = len(fmt.Sprintf("%0.2f", maxNum))
	} else {
		// 有负数的情况，最小的负数画1，最大的画满
		for i := range items {
			items[i].count = int(math.Round((items[i].Num - minNum) * float64(ctx.option.MaxBarLength) / (maxNum - minNum)))
			if items[i].count == 0 {
				items[i].count = 1
			}
		}
		maxCountLength = calcMaxCount(items)
		maxn := len(fmt.Sprintf("%0.2f", maxNum))
		minn := len(fmt.Sprintf("%0.2f", minNum))
		if maxn > minn {
			maxLengthOfNum = maxn
		} else {
			maxLengthOfNum = minn
		}
	}

	maxLengthOfName := calcMaxLengthOfName(items)
	//tmpl := fmt.Sprintf("%%%d.2f | %%-%ds | %%-%ds\n", maxLengthOfNum, maxCountLength, maxLengthOfName)
	tmpl := fmt.Sprintf("%%%d.2f | %%s%%s | %%-%ds\n", maxLengthOfNum, maxLengthOfName)
	var out string
	for _, item := range items {
		bar := strings.Repeat(ctx.option.DrawIconBlock, item.count)
		padding := strings.Repeat(ctx.option.DrawIconPadding, maxCountLength-item.count)
		out += fmt.Sprintf(tmpl, item.Num, bar, padding, item.Name)
	}
	return out
}

func (ctx *Ctx) WithAnySlice(a interface{}, iterateTransFn func(originItem interface{}) Item, modOptions ...ModOption) string {
	var items []Item
	dataops.IterateInterfaceAsSlice(a, func(iterItem interface{}) {
		items = append(items, iterateTransFn(iterItem))
	})
	return ctx.WithItems(items)
}

func (ctx *Ctx) WithMap(m map[string]int) string {
	var items []Item

	for k, v := range m {
		item := Item{
			Name: k,
			Num:  float64(v),
		}
		items = append(items, item)
	}

	return ctx.WithItems(items)
}

func (ctx *Ctx) WithMapFloat(m map[string]float64) string {
	var items []Item

	for k, v := range m {
		item := Item{
			Name: k,
			Num:  v,
		}
		items = append(items, item)
	}

	return ctx.WithItems(items)
}

func (ctx *Ctx) WithCsv(filename string) (string, error) {
	// 读取
	fp, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	r := csv.NewReader(fp)
	records, err := r.ReadAll()
	if err != nil {
		return "", err
	}

	var items []Item
	for _, line := range records {
		var item Item
		item.Name = line[0]
		item.Num, err = strconv.ParseFloat(line[1], 64)
		if err != nil {
			return "", err
		}
		items = append(items, item)
	}

	return ctx.WithItems(items), nil
}

// ---------------------------------------------------------------------------------------------------------------------

// Num最大值
func calcMinMaxNum(items []Item) (min, max float64) {
	max = math.SmallestNonzeroFloat64
	min = math.MaxFloat64
	for _, item := range items {
		if item.Num > max {
			max = item.Num
		}
		if item.Num < min {
			min = item.Num
		}
	}
	return
}

// count最大值
func calcMaxCount(items []Item) int {
	var max int
	for _, item := range items {
		if item.count > max {
			max = item.count
		}
	}
	return max
}

func calcMaxLengthOfName(items []Item) int {
	var max int
	for _, item := range items {
		if len(item.Name) > max {
			max = len(item.Name)
		}
	}
	return max
}

// ---------------------------------------------------------------------------------------------------------------------
