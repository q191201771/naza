// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// Package assert 提供了单元测试时的断言功能，减少一些模板代码
package assert

import (
	"bytes"
	"reflect"
)

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
	if !equal(expected, actual) {
		t.Errorf("%s expected=%+v, actual=%+v", msg, expected, actual)
	}
	return
}

func IsNotNil(t TestingT, actual interface{}, msg ...string) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if isNil(actual) {
		t.Errorf("%s expected not nil, but actual=%+v", msg, actual)
	}
	return
}

func isNil(actual interface{}) bool {
	if actual == nil {
		return true
	}
	v := reflect.ValueOf(actual)
	k := v.Kind()
	if k == reflect.Chan || k == reflect.Map || k == reflect.Ptr || k == reflect.Interface || k == reflect.Slice {
		return v.IsNil()
	}
	return false
}

func equal(expected, actual interface{}) bool {
	if expected == nil {
		return isNil(actual)
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	//if exp == nil || act == nil {
	//	return exp == nil && act == nil
	//}
	return bytes.Equal(exp, act)
}
