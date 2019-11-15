// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package ic

import (
	"encoding/binary"
)

type LFCompressor struct {
	FB      uint32 // 用几个字节的 bit 表示跟随的数据
	ZlibExt bool   // 压缩之后，是否再用 zlib 进一步压缩

	oc OriginCompressor // FB 为0时，退化成使用 OriginCompressor
}

// 传入的整型切片必须是从小到大有序排列
func (lfc *LFCompressor) Marshal(ids []uint32) (ret []byte) {
	if lfc.FB == 0 {
		ret = lfc.oc.Marshal(ids)
		if lfc.ZlibExt {
			ret = zlibWrite(ret)
		}
		return ret
	}

	lBuf := make([]byte, 4)
	fBuf := make([]byte, lfc.FB)

	maxDiff := 8 * lfc.FB

	var hasLeader bool
	var leader uint32
	var stage int
	for i := range ids {
		if !hasLeader {
			stage = 1
			leader = ids[i]
			hasLeader = true
			continue
		}

		diff := uint32(ids[i] - leader)

		if diff > maxDiff {
			binary.LittleEndian.PutUint32(lBuf, leader)
			ret = append(ret, lBuf...)
			ret = append(ret, fBuf...)

			resetBuf(fBuf)
			stage = 2
			leader = ids[i]
		} else {
			stage = 3
			fBuf[(diff-1)/8] = fBuf[(diff-1)/8] | (1 << byte((diff-1)%8))
		}
	}

	switch stage {
	case 1:
		binary.LittleEndian.PutUint32(lBuf, leader)
		ret = append(ret, lBuf...)
		dummy := make([]byte, lfc.FB)
		ret = append(ret, dummy...)
	case 2:
		binary.LittleEndian.PutUint32(lBuf, leader)
		ret = append(ret, lBuf...)
		dummy := make([]byte, lfc.FB)
		ret = append(ret, dummy...)
	case 3:
		binary.LittleEndian.PutUint32(lBuf, leader)
		ret = append(ret, lBuf...)
		ret = append(ret, fBuf...)
	}
	if lfc.ZlibExt {
		ret = zlibWrite(ret)
	}
	return
}

func (lfc *LFCompressor) Unmarshal(b []byte) (ids []uint32) {
	if lfc.ZlibExt {
		b = zlibRead(b)
	}
	if lfc.FB == 0 {
		return lfc.oc.Unmarshal(b)
	}

	isLeaderStage := true
	var item uint32
	var leader uint32
	var index uint32
	for {
		if isLeaderStage {
			leader = binary.LittleEndian.Uint32(b[index:])
			ids = append(ids, leader)
			isLeaderStage = false
			index += 4
		} else {
			for i := uint32(0); i < lfc.FB; i++ {
				for j := uint32(0); j < 8; j++ {
					if ((b[index+i] >> j) & 1) == 1 {
						item = leader + (i * 8) + j + 1
						ids = append(ids, item)
					}
				}
			}

			isLeaderStage = true
			index += lfc.FB
		}

		if int(index) == len(b) {
			break
		}
	}
	return
}
