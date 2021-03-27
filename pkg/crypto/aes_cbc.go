// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

var CommonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

// @param key 16字节 -> AES128
//            24字节 -> AES192
//            32字节 -> AES256
func EncryptAESWithCBC(in []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCBCEncrypter(block, iv)
	out := make([]byte, len(in))
	copy(out, in)
	encrypter.CryptBlocks(out, out)
	return out, nil
}

func DecryptAESWithCBC(in []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCBCDecrypter(block, iv)
	out := make([]byte, len(in))
	copy(out, in)
	decrypter.CryptBlocks(out, out)
	return out, nil
}
