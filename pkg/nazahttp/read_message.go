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
type HTTPReader interface {
	LineReader
	io.Reader
}

type HTTPMsgCtx struct {
	ReqMethodOrRespVersion string
	ReqURIOrRespStatusCode string
	ReqVersionOrRespReason string
	Headers                map[string]string
	Body                   []byte
}

type HTTPReqMsgCtx struct {
	Method  string
	URI     string
	Version string
	Headers map[string]string
	Body    []byte
}

type HTTPRespMsgCtx struct {
	Version    string
	StatusCode string
	Reason     string
	Headers    map[string]string
	Body       []byte
}

func ReadHTTPRequestMessage(r HTTPReader) (ctx HTTPReqMsgCtx, err error) {
	msgCtx, err := ReadHTTPMessage(r)
	if err != nil {
		return
	}
	ctx.Method = msgCtx.ReqMethodOrRespVersion
	ctx.URI = msgCtx.ReqURIOrRespStatusCode
	ctx.Version = msgCtx.ReqVersionOrRespReason
	ctx.Headers = msgCtx.Headers
	ctx.Body = msgCtx.Body
	return
}

func ReadHTTPResponseMessage(r HTTPReader) (ctx HTTPRespMsgCtx, err error) {
	msgCtx, err := ReadHTTPMessage(r)
	if err != nil {
		return
	}
	ctx.Version = msgCtx.ReqMethodOrRespVersion
	ctx.StatusCode = msgCtx.ReqURIOrRespStatusCode
	ctx.Reason = msgCtx.ReqVersionOrRespReason
	ctx.Headers = msgCtx.Headers
	ctx.Body = msgCtx.Body
	return
}

// 注意，如果HTTP Header中不包含`Content-Length`，则不会读取HTTP Body，并且err返回值为nil
func ReadHTTPMessage(r HTTPReader) (ctx HTTPMsgCtx, err error) {
	var requestLine string
	requestLine, ctx.Headers, err = ReadHTTPHeader(r)
	if err != nil {
		return ctx, err
	}
	ctx.ReqMethodOrRespVersion, ctx.ReqURIOrRespStatusCode, ctx.ReqVersionOrRespReason, err = ParseHTTPRequestLine(requestLine)
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
