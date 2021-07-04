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
)

const (
	OrderOrigin = iota + 1 // 原始序
	OrderAsc               // 升序
	OrderDesc              // 降序
)

type Item struct {
	Name string  // key
	Num  float64 // value

	count int // bar
}

var (
	//barList  = "▏▎▍▌▋▊▉█"

	// config
	maxLength = 50
	order     = OrderDesc
)

func WithItems(items []Item) (string, error) {
	// 最大的画满柱状条，其他的按与最大占比画
	maxNum := calcMaxNum(items)
	for i := range items {
		items[i].count = int(math.Round(items[i].Num * float64(maxLength) / maxNum))
		// 最小可能和最大的比太小了
		if items[i].count == 0 {
			items[i].count = 1
		}
	}

	if order == OrderDesc {
		sort.Slice(items, func(i, j int) bool {
			return items[i].count > items[j].count
		})
	}

	maxNameLength := calcMaxNameLen(items)
	maxCountLength := calcMaxCount(items)
	tmpl := fmt.Sprintf("  %%%ds | %%-%ds | %%0.2f\n", maxNameLength, maxCountLength)
	_ = maxNameLength
	var out string
	for _, item := range items {
		bar := strings.Repeat("█", item.count)
		out += fmt.Sprintf(tmpl, item.Name, bar, item.Num)
	}
	return out, nil
}

func WithMap(m map[string]int) (string, error) {
	var items []Item

	for k, v := range m {
		item := Item{
			Name: k,
			Num:  float64(v),
		}
		items = append(items, item)
	}

	return WithItems(items)
}

func WithMapFloat(m map[string]float64) (string, error) {
	var items []Item

	for k, v := range m {
		item := Item{
			Name: k,
			Num:  v,
		}
		items = append(items, item)
	}

	return WithItems(items)
}

func WithCsv(filename string) (string, error) {
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

	return WithItems(items)
}

func isFloat(v string) bool {
	return strings.Contains(v, ".")
}

func calcMaxNum(items []Item) float64 {
	var max float64
	for _, item := range items {
		if item.Num > max {
			max = item.Num
		}
	}
	return max
}

func calcMaxNameLen(items []Item) int {
	var max int
	for _, item := range items {
		if len(item.Name) > max {
			max = len(item.Name)
		}
	}
	return max
}

func calcMaxCount(items []Item) int {
	var max int
	for _, item := range items {
		if item.count > max {
			max = item.count
		}
	}
	return max
}
