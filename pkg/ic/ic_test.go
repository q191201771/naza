// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package ic

import (
	"log"
	"testing"
)

//// 从文件中读取 uid 列表
//func obtainUIDList(filename string) (uids IDSlice) {
//	fp, err := os.Open(filename)
//	if err != nil {
//		panic(err)
//	}
//	buf, err := ioutil.ReadAll(fp)
//	if err != nil {
//		panic(err)
//	}
//	lines := bytes.Split(buf, []byte("\n"))
//	for _, line := range lines {
//		if len(line) == 0 {
//			continue
//		}
//		item, err := strconv.ParseUint(string(line), 10, 32)
//		if err != nil {
//			panic(err)
//		}
//		uids = append(uids, uint32(item))
//	}
//	return uids
//}

//var FILENAME = "uid.txt"

func marshalWrap(ids []uint32) (ret []byte) {
	log.Println("> sort.")
	Sort(ids)
	log.Println("< sort.")

	log.Println("> marshal.")
	//var oc OriginCompressor
	//ret = oc.Marshal(ids)

	var lfc LFCompressor
	lfc.FB = 4
	ret = lfc.Marshal(ids)
	log.Println("< marshal.")

	log.Println("> zlib. len:", len(ret))
	ret = zlibWrite(ret)
	log.Println("< zlib. len:", len(ret))
	return
}

func unmarshalWrap(b []byte) (ret []uint32) {
	b = zlibRead(b)

	//var oc OriginCompressor
	//ret = oc.Unmarshal(b)

	var lfc LFCompressor
	lfc.FB = 4
	ret = lfc.Unmarshal(b)
	return
}

func TestIC(t *testing.T) {
	log.SetFlags(log.Lmicroseconds)

	// 单元测试 case
	uidss := [][]uint32{
		{1, 2, 3, 18, 32, 100},
		{1, 2, 3, 18, 32},
		{1, 2, 3, 18},
		{1, 2, 3, 17},
		{1, 2, 3, 16},
		{1, 2, 3, 15, 16, 17, 18},
		{1, 2, 3, 15, 16, 17},
		{1, 2, 3, 15, 16},
		{1, 2, 3, 15},
		{1, 2, 3},
		{1, 2},
		{1},
	}

	var compressors []Compressor

	compressors = append(compressors, &OriginCompressor{})
	compressors = append(compressors, &OriginCompressor{ZlibExt: true})

	compressors = append(compressors, &LFCompressor{FB: 0})
	compressors = append(compressors, &LFCompressor{FB: 0, ZlibExt: true})
	compressors = append(compressors, &LFCompressor{FB: 2})
	compressors = append(compressors, &LFCompressor{FB: 4})
	compressors = append(compressors, &LFCompressor{FB: 4, ZlibExt: true})

	for _, c := range compressors {
		for _, uids := range uidss {
			log.Println("-----")
			log.Println("in uid len:", len(uids))

			Sort(uids)
			b := c.Marshal(uids)
			log.Println("len(b):", len(b))

			uids2 := c.Unmarshal(b)
			log.Println("out uid len:", len(uids2))

			// assert check
			if len(uids) != len(uids2) {
				panic(0)
			}
			for i := range uids {
				if uids[i] != uids2[i] {
					panic(0)
				}
			}
			log.Println("-----")
		}
	}
}
