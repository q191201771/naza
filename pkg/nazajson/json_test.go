// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazajson

import (
	"encoding/json"
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

func TestCollectNotExistFields(t *testing.T) {
	// 1. 测试自身的基础字段     完成
	// 1.2 测试数组            完成
	// 2. 测试自身的指针基础字段  完成
	// 3. 测试匿名结构体        完成
	// 4. 测试嵌套结构体        完成
	// 5. 测试没有tag          完成
	// 6. 测试小写不暴露的成员   完成

	type Sub struct {
		SubA int    `json:"suba"`
		SubB int    `json:"subb"`
		SubC string `json:"subc"`
		SubD string `json:"subd"`
	}
	type Anoy struct {
		AnoyA int `json:"anoya"`
		AnoyB int `json:"anoyb"`
	}
	type St struct {
		Anoy
		A     int     `json:"a"`
		B     string  `json:"b"`
		C     []bool  `json:"c"`
		D     []int   `json:"d"`
		E     *int    `json:"e"`
		F     *string `json:"f"`
		Sub   Sub     `json:"sub"`
		NoTag int
		low   int
	}

	b := []byte(`
{
  "anoya": 3,
  "a": 1,
  "c": [true, false],
  "e": 2,
  "sub": {
    "suba": 4,
    "subc": "c"
  }
}
`)
	var st St
	json.Unmarshal(b, &st)
	nazalog.Infof("%+v", st)

	//var st2 St
	//collect, err := CollectNotExistFields(b, st)
	//CollectNotExistFields(b, st2)
	//CollectNotExistFields(b, &st2)
	// 以上三种使用方式也都是正常的
	collect, err := CollectNotExistFields(b, &st)
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"anoyb", "b", "d", "f", "sub.subb", "sub.subd"}, collect)

	collect, err = CollectNotExistFields(b, st, "sub")
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"anoyb", "b", "d", "f"}, collect)

	collect, err = CollectNotExistFields(b, st, "s")
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"anoyb", "b", "d", "f"}, collect)
}
