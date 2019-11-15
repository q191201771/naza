package ic

import "encoding/binary"

type OriginCompressor struct {
}

func (oc *OriginCompressor) Marshal(ids IDSlice) (ret []byte) {
	ret = make([]byte, len(ids)*4)
	for i, id := range ids {
		binary.LittleEndian.PutUint32(ret[i*4:], id)
	}
	return
}

func (oc *OriginCompressor) Unmarshal(b []byte) (ids IDSlice) {
	n := len(b) / 4
	for i := 0; i < n; i++ {
		id := binary.LittleEndian.Uint32(b[i*4:])
		ids = append(ids, id)
	}
	return
}
