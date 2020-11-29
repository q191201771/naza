// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazahttp

import (
	"io"
	"strconv"
)

// e.g. bufio.Reader
type RequestReader interface {
	LineReader
	io.Reader
}

type HTTPRequestCtx struct {
	Method  string
	URI     string
	Headers map[string]string
	Body    []byte
}

// 注意，如果HTTP Header中不包含`Content-Length`，则不会读取HTTP Body，并且err返回值为nil
func ReadHTTPRequest(r RequestReader) (ctx HTTPRequestCtx, err error) {
	var requestLine string
	requestLine, ctx.Headers, err = ReadHTTPHeader(r)
	if err != nil {
		return ctx, err
	}
	ctx.Method, ctx.URI, _, err = ParseHTTPRequestLine(requestLine)
	if err != nil {
		return ctx, err
	}

	contentLength, ok := ctx.Headers[HeaderFieldContentLength]
	if !ok {
		return ctx, nil
	}
	cl, err := strconv.Atoi(contentLength)
	if err != nil {
		return ctx, err
	}
	ctx.Body = make([]byte, cl)
	_, err = io.ReadFull(r, ctx.Body)

	return ctx, err
}
