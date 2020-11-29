// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazahttp

import "errors"

var (
	ErrHTTPHeader   = errors.New("nazahttp: fxxk")
	ErrParamMissing = errors.New("nazahttp: param missing")
)

const (
	HeaderFieldContentLength = "Content-Length"
	HeaderFieldContentType   = "application/json"
)
