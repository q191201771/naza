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
	v := []Item{
		//{Name: "China", Num: 1},
		{Name: "中", Num: 22},
		{Name: "中国", Num: 333},
		{Name: "中国啊", Num: 4444},
	}
	for i := range v {
		fmt.Println(len(v[i].Name))
	}
	fmt.Println(WithItems(v, nil))
}
