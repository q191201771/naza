// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazastring

import (
	"bytes"
	"fmt"
	"unsafe"
)

type sliceT struct {
	array unsafe.Pointer
	len   int
	cap   int
}

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func SliceByteToStringTmp(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToSliceByteTmp(s string) []byte {
	str := (*stringStruct)(unsafe.Pointer(&s))
	ret := sliceT{array: unsafe.Pointer(str.str), len: str.len, cap: str.len}
	return *(*[]byte)(unsafe.Pointer(&ret))
}

// 我在写单元测试时经常遇到一个场景，
// 测试结果输出的字节切片很长，我在第一次编写时，人肉确认输出结果是正确的后，
// 我需要将这个字节切片的值作为期望值，硬编码进单元测试的代码中，供后续每次单元测试做验证。
//
// 有了这个函数，我可以在第一次编写时，调用该函数，将得到的结果拷贝至单元测试的代码中，
// 之后将调用该函数的代码删除。
//
func DumpSliceByte(b []byte) string {
	if len(b) == 0 {
		return "nil"
	}

	var bb bytes.Buffer
	bb.WriteString("[]byte{")
	for i := range b {
		if i != len(b)-1 {
			bb.WriteString(fmt.Sprintf("0x%02x, ", b[i]))

		} else {
			bb.WriteString(fmt.Sprintf("0x%02x", b[i]))
		}
	}
	bb.WriteString("}")
	return bb.String()
}

// @param m 对切片做子切片操作，取值范围大于0，如果超过了原切片大小，则返回原切片
func SubSliceSafety(b []byte, m int) []byte {
	if m <= len(b) {
		return b[:m]
	}
	return b
}
