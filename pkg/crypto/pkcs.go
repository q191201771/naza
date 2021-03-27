// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package crypto

import (
	"bytes"
	"errors"
)

var ErrPKCS = errors.New("naza.crypto: fxxk")

// @param blockSize 取值范围[0, 255]
//                  如果是AES，见标准库中aes.BlockSize等于16
func EncryptPKCS7(in []byte, blockSize int) []byte {
	paddingLength := blockSize - len(in)%blockSize
	paddingBuf := bytes.Repeat([]byte{byte(paddingLength)}, paddingLength)
	return append(in, paddingBuf...)
}

func DecryptPKCS7(in []byte) ([]byte, error) {
	totalLength := len(in)
	if totalLength < 1 {
		return nil, ErrPKCS
	}
	paddingLength := int(in[totalLength-1])
	if totalLength < paddingLength {
		return nil, ErrPKCS
	}
	return in[:totalLength-int(paddingLength)], nil
}

func EncryptPKCS5(in []byte) []byte {
	return EncryptPKCS7(in, 8)
}

func DecryptPKCS5(in []byte) ([]byte, error) {
	return DecryptPKCS7(in)
}
