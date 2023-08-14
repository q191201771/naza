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

// WithItems
//
// @param items: 注意，内部不会修改切片底层数据的值以及顺序
func (ctx *Ctx) WithItems(items []Item) string {
	// 拷贝一份，避免修改外部切片的原始顺序
	if ctx.option.Order != OrderOrigin {
		clone := make([]Item, len(items))
		copy(clone, items)
		items = clone
	}

	// 排序
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

	// 选取需要的元素
	var newItems []Item
	dataops.SliceLimit(items, ctx.option.PrefixNumLimit, ctx.option.SuffixNumLimit, func(index int) {
		newItems = append(newItems, items[index])
	})
	items = newItems

	var (
		maxLengthOfCount int // count柱状最长画多长
		maxLengthOfNum   int // num字段多长
	)

	minN, maxN := dataops.SliceMinMax(items, func(i, j int) bool {
		return items[i].Num < items[j].Num
	})
	minNum := minN.(Item).Num
	maxNum := maxN.(Item).Num

	isAllInteger := dataops.SliceAllOf(items, func(originItem interface{}) bool {
		return isInteger(originItem.(Item).Num)
	})

	if isAllInteger && (int(maxNum-minNum) < ctx.option.MaxBarLength) {
		// 如果都是整数，且实际最大值最小值的差值小于柱状最大长度限制

		for i := range items {
			if minNum >= 0.00 {
				// 都是正整数,按原始值绘制
				items[i].count = int(items[i].Num)
			} else {
				// 最小的负值画1
				items[i].count = int(items[i].Num - minNum + 1)
			}
		}
	} else {
		for i := range items {
			if minNum > 0.00 {
				// 都是正数的情况，最大的画满柱状条，其他的按与最大占比画
				// round四舍五入
				items[i].count = int(math.Round(items[i].Num * float64(ctx.option.MaxBarLength) / maxNum))
			} else {
				// 有负数的情况，最小的负数画1，最大的画满
				items[i].count = int(math.Round((items[i].Num - minNum) * float64(ctx.option.MaxBarLength) / (maxNum - minNum)))
			}

			// 最小可能和最大的比太小了
			if items[i].count == 0 {
				items[i].count = 1
			}
		}
	}

	maxLengthOfCount = dataops.SliceMax(items, func(i, j int) bool {
		return items[i].count < items[j].count
	}).(Item).count

	maxn := len(fmt.Sprintf("%0.2f", maxNum))
	minn := len(fmt.Sprintf("%0.2f", minNum))
	if maxn > minn {
		maxLengthOfNum = maxn
	} else {
		maxLengthOfNum = minn
	}

	maxLengthOfName := len(dataops.SliceMax(items, func(i, j int) bool {
		return len(items[i].Name) < len(items[j].Name)
	}).(Item).Name)

	var tmpl string
	var tmplNum string
	if isAllInteger {
		// -3是因为整数不需要小数点和小数点的后两位
		tmplNum = fmt.Sprintf("%%%d.0f", maxLengthOfNum-3)
	} else {
		tmplNum = fmt.Sprintf("%%%d.2f", maxLengthOfNum)
	}
	if !ctx.option.HideNum && !ctx.option.HideName {
		tmpl = fmt.Sprintf("%s | %%s%%s | %%-%ds\n", tmplNum, maxLengthOfName)
	} else if !ctx.option.HideNum && ctx.option.HideName {
		tmpl = fmt.Sprintf("%s | %%s%%s\n", tmplNum)
	} else if ctx.option.HideNum && !ctx.option.HideName {
		tmpl = fmt.Sprintf("%%s%%s | %%-%ds\n", maxLengthOfName)
	} else {
		tmpl = "%s\n"
	}
	var out string
	for _, item := range items {
		bar := strings.Repeat(ctx.option.DrawIconBlock, item.count)
		padding := strings.Repeat(ctx.option.DrawIconPadding, maxLengthOfCount-item.count)

		if !ctx.option.HideNum && !ctx.option.HideName {
			out += fmt.Sprintf(tmpl, item.Num, bar, padding, item.Name)
		} else if !ctx.option.HideNum && ctx.option.HideName {
			out += fmt.Sprintf(tmpl, item.Num, bar, padding)
		} else if ctx.option.HideNum && !ctx.option.HideName {
			out += fmt.Sprintf(tmpl, bar, padding, item.Name)
		} else {
			out += fmt.Sprintf(tmpl, bar)
		}
	}
	return out
}

func (ctx *Ctx) WithAnySlice(a interface{}, iterateTransFn func(originItem interface{}) Item, modOptions ...ModOption) string {
	var items []Item
	dataops.IterateInterfaceAsSlice(a, func(iterItem interface{}) bool {
		items = append(items, iterateTransFn(iterItem))
		return true
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
