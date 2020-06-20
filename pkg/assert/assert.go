// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// package assert 提供了单元测试时的断言功能，减少一些模板代码
package assert

import "github.com/q191201771/naza/pkg/nazareflect"

// 单元测试中的 *testing.T 和 *testing.B 都满足该接口
type TestingT interface {
	Errorf(format string, args ...interface{})
}

type tHelper interface {
	Helper()
}

func Equal(t TestingT, expected interface{}, actual interface{}, msg ...string) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !nazareflect.Equal(expected, actual) {
		t.Errorf("%s expected=%+v, actual=%+v", msg, expected, actual)
	}
	return
}

// 比如有时我们需要对 error 类型不等于 nil 做断言，但是我们并不关心 error 的具体值是什么
func IsNotNil(t TestingT, actual interface{}, msg ...string) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if nazareflect.IsNil(actual) {
		t.Errorf("%s expected not nil, but actual=%+v", msg, actual)
	}
	return
}
