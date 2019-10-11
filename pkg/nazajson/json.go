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
	"strings"
)

type JSON struct {
	//raw []byte
	m map[string]interface{}
}

func New(raw []byte) (JSON, error) {
	var j JSON
	err := j.Init(raw)
	return j, err
}

func (j *JSON) Init(raw []byte) error {
	return json.Unmarshal(raw, &j.m)
}

// 判断 json 中某个字段是否存在
// @param path 支持多级格式，用句号`.`分隔，比如 log.level
func (j *JSON) Exist(path string) bool {
	return exist(j.m, path)
}

func exist(m map[string]interface{}, path string) bool {
	ps := strings.Split(path, ".")

	if len(ps) > 1 {
		v, ok := m[ps[0]]
		if !ok {
			return false
		}
		mm, ok := v.(map[string]interface{})
		if !ok {
			return false
		}
		return exist(mm, ps[1])
	}

	_, ok := m[ps[0]]
	return ok
}
