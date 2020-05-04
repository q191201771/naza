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
	"math/bits"
)

func main() {
	var a uint32
	a = 1
	fmt.Printf("%032b\n", a)
	fmt.Printf("%032b\n", bits.ReverseBytes32(a))
	a = a << 8
	fmt.Printf("%032b\n", a)
	fmt.Printf("%032b\n", bits.ReverseBytes32(a))

	var b uint64
	b = 1
	fmt.Printf("%064b\n", b)
	fmt.Printf("%064b\n", bits.ReverseBytes64(b))
	b = b << 8
	fmt.Printf("%064b\n", b)
	fmt.Printf("%064b\n", bits.ReverseBytes64(b))
	b = b << 8
	fmt.Printf("%064b\n", b)
	fmt.Printf("%064b\n", bits.ReverseBytes64(b))
}
