// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazahttp

import (
	"github.com/q191201771/naza/pkg/nazaerrors"
	"net/http"
	"strings"
)

type LineReader interface {
	ReadLine() (line []byte, isPrefix bool, err error)
}

// ReadHttpHeader
//
// @return firstLine: request的request line或response的status line
// @return headers:   request header fields的键值对
//
func ReadHttpHeader(r LineReader) (firstLine string, headers http.Header, err error) {
	headers = make(http.Header)

	readLineFn := func() (string, error) {
		var line string
		var bline []byte
		var isPrefix bool
		for {
			bline, isPrefix, err = r.ReadLine()
			if err != nil {
				return "", err
			}
			line += string(bline)
			if !isPrefix {
				break
			}
		}
		return line, nil
	}

	firstLine, err = readLineFn()
	if err != nil {
		err = nazaerrors.Wrap(err, firstLine)
		return
	}
	if len(firstLine) == 0 {
		err = nazaerrors.Wrap(ErrHttpHeader, firstLine)
		return
	}

	for {
		var l string
		l, err = readLineFn()
		if err != nil {
			err = nazaerrors.Wrap(err, l)
			return
		}
		if len(l) == 0 { // 读到一个空的 \r\n 表示http头全部读取完毕了
			break
		}

		pos := strings.Index(l, ":")
		if pos == -1 {
			err = nazaerrors.Wrap(ErrHttpHeader, l)
			return
		}
		headers.Add(strings.Trim(l[0:pos], " "), strings.Trim(l[pos+1:], " "))
	}
	return
}

// Request-Line = Method SP URI SP Version CRLF
func ParseHttpRequestLine(line string) (method string, uri string, version string, err error) {
	return parseFirstLine(line)
}

// Status-Line = Version SP Status-Code SP Reason CRLF
func ParseHttpStatusLine(line string) (version string, statusCode string, reason string, err error) {
	return parseFirstLine(line)
}

func parseFirstLine(line string) (item1, item2, item3 string, err error) {
	f := strings.Index(line, " ")
	if f == -1 {
		err = nazaerrors.Wrap(ErrFirstLine, line)
		return
	}
	s := strings.Index(line[f+1:], " ")
  if s == -1 {
    return line[0:f], line[f+1:], "", nil
  }
  if f+1+s+1 == len(line) {
    return line[0:f], line[f+1:f+1+s], "", nil
  }
	//if s == -1 || f+1+s+1 == len(line) {
	//	err = nazaerrors.Wrap(ErrFirstLine, line)
	//	return
	//}

	return line[0:f], line[f+1 : f+1+s], line[f+1+s+1:], nil
}
