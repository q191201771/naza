// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazajson

import (
	"testing"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazalog"
)

var raw = []byte(`
{
	"num": 1,
	"num2": 0,
	"flag": true,
	"flag2": false,
	"str": "aaa",
	"str2": "",
	"obj": {
		"onum": 2,
		"onum3": 0,
		"oflag": true,
		"oflag3": false,
		"ostr": "bbb",
		"ostr3": ""
	},
	"arr": [],
	"arr2": [1, 2]
}
`)

func TestExist(t *testing.T) {
	var exist bool

	j, err := New(raw)
	assert.Equal(t, nil, err)

	exist = j.Exist("num")
	assert.Equal(t, true, exist)
	exist = j.Exist("flag")
	assert.Equal(t, true, exist)
	exist = j.Exist("str")
	assert.Equal(t, true, exist)

	exist = j.Exist("num2")
	assert.Equal(t, true, exist)
	exist = j.Exist("flag2")
	assert.Equal(t, true, exist)
	exist = j.Exist("str2")
	assert.Equal(t, true, exist)

	exist = j.Exist("arr")
	assert.Equal(t, true, exist)

	exist = j.Exist("arr2")
	assert.Equal(t, true, exist)

	exist = j.Exist("obj")
	assert.Equal(t, true, exist)
	exist = j.Exist("obj.onum")
	assert.Equal(t, true, exist)
	exist = j.Exist("obj.oflag")
	assert.Equal(t, true, exist)
	exist = j.Exist("obj.ostr")
	assert.Equal(t, true, exist)

	exist = j.Exist("obj.onum3")
	assert.Equal(t, true, exist)
	exist = j.Exist("obj.oflag3")
	assert.Equal(t, true, exist)
	exist = j.Exist("obj.ostr3")
	assert.Equal(t, true, exist)

	exist = j.Exist("notexist")
	assert.Equal(t, false, exist)

	exist = j.Exist("notexist.notexist")
	assert.Equal(t, false, exist)

	exist = j.Exist("obj.notexist")
	assert.Equal(t, false, exist)

	exist = j.Exist("obj.notexist.notexist")
	assert.Equal(t, false, exist)

	exist = j.Exist("num.notexist")
	assert.Equal(t, false, exist)

	exist = j.Exist(".")
	assert.Equal(t, false, exist)
	exist = j.Exist("..")
	assert.Equal(t, false, exist)
}

func BenchmarkExist(b *testing.B) {
	var exist bool

	j, _ := New(raw)

	for i := 0; i < b.N; i++ {
		exist = j.Exist("num")
		exist = j.Exist("flag")
		exist = j.Exist("str")

		exist = j.Exist("obj")
		exist = j.Exist("obj.onum")
		exist = j.Exist("obj.oflag")
		exist = j.Exist("obj.ostr")

		exist = j.Exist("notexist")

		exist = j.Exist("obj.notexist")

		exist = j.Exist("obj.notexist.notexist")

		exist = j.Exist(".")
		exist = j.Exist("..")
	}
	nazalog.Debug(exist)
}
