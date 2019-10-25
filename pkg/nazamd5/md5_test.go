// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazamd5

import (
	"testing"

	"github.com/q191201771/naza/pkg/assert"
)

type md5Test struct {
	in  string
	out string
}

func TestMD5(t *testing.T) {
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", MD5(nil))
	golden := []md5Test{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"aaa", "47bce5c74f589f4867dbd57e9ca9f808"},
		{"AAA", "e1faffb3e614e6c2fba74296962386b7"},
		{"HELLO WORLD!", "b59bc37d6441d96785bda7ab2ae98f75"},
	}
	for _, g := range golden {
		assert.Equal(t, g.out, MD5([]byte(g.in)))
	}
}
