// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// package ic 将整型切片压缩成二进制字节切片
package ic

// 具体使用见 LFCompressor 和 OriginCompressor
type Compressor interface {
	// 将整型切片压缩成二进制字节切片
	Marshal(ids []uint32) (ret []byte)
	// 将二进制字节切片反序列化为整型切片
	// 反序列化后得到的整型切片，切片中整型的顺序和序列化之前保持不变
	Unmarshal(b []byte) (ids []uint32)
}
