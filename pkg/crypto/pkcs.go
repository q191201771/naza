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

var ErrPkcs = errors.New("naza.crypto: fxxk")

// @param blockSize 取值范围[0, 255]
//
//	如果是AES，见标准库中aes.BlockSize等于16
func EncryptPkcs7(in []byte, blockSize int) []byte {
	paddingLength := blockSize - len(in)%blockSize
	paddingBuf := bytes.Repeat([]byte{byte(paddingLength)}, paddingLength)
	return append(in, paddingBuf...)
}

func DecryptPkcs7(in []byte) ([]byte, error) {
	totalLength := len(in)
	if totalLength < 1 {
		return nil, ErrPkcs
	}
	paddingLength := int(in[totalLength-1])
	if totalLength < paddingLength {
		return nil, ErrPkcs
	}
	return in[:totalLength-int(paddingLength)], nil
}

func EncryptPkcs5(in []byte) []byte {
	return EncryptPkcs7(in, 8)
}

func DecryptPkcs5(in []byte) ([]byte, error) {
	return DecryptPkcs7(in)
}
