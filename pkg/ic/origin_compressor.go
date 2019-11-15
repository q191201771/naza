package ic

import "encoding/binary"

type OriginCompressor struct {
	ZlibExt bool // 压缩之后，是否再用 zlib 进一步压缩
}

// 并不强制要求整型切片有序
func (oc *OriginCompressor) Marshal(ids []uint32) (ret []byte) {
	ret = make([]byte, len(ids)*4)
	for i, id := range ids {
		binary.LittleEndian.PutUint32(ret[i*4:], id)
	}
	if oc.ZlibExt {
		ret = zlibWrite(ret)
	}
	return
}

func (oc *OriginCompressor) Unmarshal(b []byte) (ids []uint32) {
	if oc.ZlibExt {
		b = zlibRead(b)
	}
	n := len(b) / 4
	for i := 0; i < n; i++ {
		id := binary.LittleEndian.Uint32(b[i*4:])
		ids = append(ids, id)
	}
	return
}
