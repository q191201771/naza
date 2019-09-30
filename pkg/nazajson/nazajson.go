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
