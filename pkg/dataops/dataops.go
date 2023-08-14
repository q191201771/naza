// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package dataops

import (
	"reflect"
	"time"
)

// Slice2Strings 将任意类型切片转换为字符串切片
//
// @param a:  任意类型切片，比如结构体切片
// @param fn: 业务方编写转换逻辑，内部对原始切片的元素逐个回调给业务方，并通过回调返回值组成转换后的字符串切片
//
// @return ret: 转换后的字符串切片
func Slice2Strings(a interface{}, fn func(originItem interface{}) string) (ret []string) {
	IterateInterfaceAsSlice(a, func(iterItem interface{}) bool {
		ret = append(ret, fn(iterItem))
		return true
	})
	return
}

// Slice2Times 将任意类型切片转换为时间切片
func Slice2Times(a interface{}, fn func(originItem interface{}) time.Time) (ret []time.Time) {
	IterateInterfaceAsSlice(a, func(iterItem interface{}) bool {
		ret = append(ret, fn(iterItem))
		return true
	})
	return
}

// ---------------------------------------------------------------------------------------------------------------------

// SliceUniqueCount 遍历切片`a`，逐个调用`fn`转换为string, 并将所有元素归类计数
func SliceUniqueCount(a interface{}, fn func(originItem interface{}) string) (ret map[string]int) {
	ret = make(map[string]int)
	IterateInterfaceAsSlice(a, func(iterItem interface{}) bool {
		k := fn(iterItem)
		ret[k] += 1
		return true
	})
	return
}

// ---------------------------------------------------------------------------------------------------------------------

// SliceLimit 取切片前`PrefixNumLimit`个元素和后`SuffixNumLimit`个元素，通过`cb`回调给业务方
//
//	注意，内部会处理`PrefixNumLimit`或`SuffixNumLimit`过大的情况
//	`PrefixNumLimit`如果为-1，则没有限制，`SuffixNumLimit`同理
func SliceLimit(a interface{}, PrefixNumLimit int, SuffixNumLimit int, cb func(index int)) {
	v := reflect.ValueOf(a)
	if PrefixNumLimit == -1 && SuffixNumLimit == -1 {
	} else if PrefixNumLimit == -1 && SuffixNumLimit != -1 {
		if SuffixNumLimit < v.Len() {
			for i := v.Len() - SuffixNumLimit; i < v.Len(); i++ {
				cb(i)
			}
			return
		}
	} else if PrefixNumLimit != -1 && SuffixNumLimit == -1 {
		if PrefixNumLimit < v.Len() {
			for i := 0; i < PrefixNumLimit; i++ {
				cb(i)
			}
			return
		}
	} else {
		if PrefixNumLimit+SuffixNumLimit < v.Len() {
			for i := 0; i < PrefixNumLimit; i++ {
				cb(i)
			}
			for i := v.Len() - SuffixNumLimit; i < v.Len(); i++ {
				cb(i)
			}
			return
		}
	}

	for i := 0; i < v.Len(); i++ {
		cb(i)
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func Map2Strings(a interface{}, fn func(k, v interface{}) string) (ret []string) {
	IterateInterfaceAsMap(a, func(k, v interface{}) {
		ret = append(ret, fn(k, v))
	})
	return
}

// ---------------------------------------------------------------------------------------------------------------------

// IterateInterfaceAsSlice
//
// 遍历切片`a`，通过`onIterate`逐个回调元素
//
// @param onIterate:
//
//	@return keepIterate: 如果返回false，则停止变化
func IterateInterfaceAsSlice(a interface{}, onIterate func(iterItem interface{}) (keepIterate bool)) {
	// TODO(chef): fix haven't use keepIterate
	v := reflect.ValueOf(a)
	for i := 0; i < v.Len(); i++ {
		onIterate(v.Index(i).Interface())
	}
}

func IterateInterfaceAsMap(a interface{}, onIterate func(k, v interface{})) {
	v := reflect.ValueOf(a)
	for _, key := range v.MapKeys() {
		s := v.MapIndex(key)
		onIterate(key.Interface(), s.Interface())
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func SliceAllOf(a interface{}, fn func(originItem interface{}) bool) (ret bool) {
	ret = true
	IterateInterfaceAsSlice(a, func(iterItem interface{}) bool {
		if !fn(iterItem) {
			ret = false
			return false
		}
		return true
	})
	return
}

func SliceMinMax(a interface{}, less func(i, j int) bool) (min, max interface{}) {
	v := reflect.ValueOf(a)
	imin := 0
	imax := 0
	if v.Len() == 1 {
		return v.Index(imin).Interface(), v.Index(imax).Interface()
	}
	for i := 1; i < v.Len(); i++ {
		if less(i, imin) {
			imin = i
		}
		if less(imax, i) {
			imax = i
		}
	}
	return v.Index(imin).Interface(), v.Index(imax).Interface()
}

func SliceMax(a interface{}, less func(i, j int) bool) (max interface{}) {
	_, max = SliceMinMax(a, less)
	return
}
