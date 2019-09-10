// package unique 对象唯一ID
package unique

import (
	"fmt"
	"sync"
)

var global Unique

func GenUniqueKey(prefix string) string {
	return global.GenUniqueKey(prefix)
}

type Unique struct {
	//id uint64

	m         sync.Mutex
	prefix2id map[string]uint64
}

func (u *Unique) GenUniqueKey(prefix string) string {
	//return fmt.Sprintf("%s%d", prefix, atomic.AddUint64(&u.id, 1))
	u.m.Lock()
	defer u.m.Unlock()
	id, ok := u.prefix2id[prefix]
	if ok {
		id++
	} else {
		id = 1
	}
	u.prefix2id[prefix] = id
	return fmt.Sprintf("%s%d", prefix, id)
}

func init() {
	global.prefix2id = make(map[string]uint64)
}
