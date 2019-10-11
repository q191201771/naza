// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package bininfo

import (
	"fmt"
	"testing"
)

func TestStringifySingleLine(t *testing.T) {
	fmt.Println(StringifySingleLine())
}

func TestStringifyMultiLine(t *testing.T) {
	fmt.Println(StringifyMultiLine())
}

func TestCorner(t *testing.T) {
	GitStatus = ""
	beauty()
}
