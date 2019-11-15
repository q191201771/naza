package ic

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"sort"
)

func Sort(ids []uint32) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
}

func resetBuf(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
	return b
}

func zlibWrite(in []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, _ = w.Write(in)
	_ = w.Close()
	return b.Bytes()
}

func zlibRead(in []byte) (ret []byte) {
	b := bytes.NewReader(in)
	r, _ := zlib.NewReader(b)
	ret, _ = ioutil.ReadAll(r)
	return
}
