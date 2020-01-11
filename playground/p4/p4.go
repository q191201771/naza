// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"fmt"
	//"sort"
)

// 把 sort.Search 拷过来，加些单步日志
func SearchWithLog(n int, f func(int) bool) int {
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1)
		fmt.Printf("i=%d, j=%d, h=%d", i, j, h)
		if !f(h) {
			fmt.Println(" not match.")
			i = h + 1
		} else {
			fmt.Println("match.")
			j = h
		}
	}
	return i
}

func main() {
	const key = 3
	//const key = 0
	//const key = 2
	//const key = 9
	arr := []int{1, 3, 5, 7}
	index := SearchWithLog(len(arr), func(i int) bool {
		//return arr[i] == 3
		return arr[i] >= key
	})
	fmt.Printf("index=%d\n", index)
	if index < len(arr) && arr[index] == key {
		fmt.Println("exist")
	} else {
		fmt.Println("not exist")
	}

	//var tmp uint
	//tmp = uint(2147483647 + 2147483647)
	//fmt.Println(tmp)
}
