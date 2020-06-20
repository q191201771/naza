// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazareflect

import (
	"errors"
	"testing"
)

func BenchmarkEqualInteger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EqualInteger("a", 1)
	}
}

func TestIsNil(t *testing.T) {
	sure(t, IsNil(nil))
	sure(t, !IsNil(1))
}

func TestEqual(t *testing.T) {
	sure(t, Equal(nil, nil))
	sure(t, Equal(1, 1))
	sure(t, Equal("aaa", "aaa"))

	var ch chan struct{}
	sure(t, Equal(nil, ch))
	var m map[string]string
	sure(t, Equal(nil, m))
	var p *int
	sure(t, Equal(nil, p))
	var i interface{}
	sure(t, Equal(nil, i))
	var b []byte
	sure(t, Equal(nil, b))

	sure(t, Equal([]byte{}, []byte{}))
	sure(t, Equal([]byte{0, 1, 2}, []byte{0, 1, 2}))

	sure(t, !Equal(nil, 1))
	sure(t, !Equal([]byte{}, "aaa"))
	sure(t, !Equal(nil, errors.New("mock error")))
}

func TestEqualInteger(t *testing.T) {
	sure(t, EqualInteger(0, 0))
	sure(t, EqualInteger(1, uint(1)))
	sure(t, EqualInteger(uint32(1), int16(1)))
	sure(t, EqualInteger(uint(1), uint8(1)))

	sure(t, !EqualInteger(1, 0))
	sure(t, !EqualInteger(0, "aaa"))
	sure(t, !EqualInteger(-1, uint(0)))
	sure(t, !EqualInteger(uint16(0), int32(-1)))
}

// 因为naza assert package引用了naza value package，如果这里再使用assert，就造成package循环引用了
// 所以这里写个简单的帮助测试的函数
func sure(t *testing.T, actual bool) {
	t.Helper()
	if !actual {
		t.Error()
	}
}
