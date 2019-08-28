// package unique 对象唯一ID
package unique

import (
	"fmt"
	"sync/atomic"
)

var global Unique

func GenUniqueKey(prefix string) string {
	return global.GenUniqueKey(prefix)
}

type Unique struct {
	id uint64
}

func (u *Unique) GenUniqueKey(prefix string) string {
	return fmt.Sprintf("%s%d", prefix, atomic.AddUint64(&u.id, 1))
}
