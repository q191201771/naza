// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazasync

import (
	"bytes"
	"errors"
	"runtime"
	"strconv"
)

// NOTICE copy from https://github.com/golang/net/blob/master/http2/gotrack.go

var ErrObtainGoroutineID = errors.New("nazasync: obtain current goroutine id failed")

func CurGoroutineID() (int64, error) {
	var goroutineSpace = []byte("goroutine ")

	b := make([]byte, 128)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		return -1, ErrObtainGoroutineID
	}
	return strconv.ParseInt(string(b[:i]), 10, 64)
}
