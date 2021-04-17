// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazasync

import (
	"fmt"
	"sync"

	"github.com/q191201771/naza/pkg/unique"

	"github.com/q191201771/naza/pkg/nazalog"
)

// 用于debug锁方面的问题

const unlockedGid = 0

var uniqueGen *unique.SingleGenerator

type Mutex struct {
	core sync.Mutex

	mu        sync.Mutex
	uniqueKey string
	gid       int64
}

func (m *Mutex) Lock() {
	m.mu.Lock()
	if m.uniqueKey == "" {
		m.uniqueKey = uniqueGen.GenUniqueKey()
	}
	gid, _ := CurGoroutineID()
	if gid == m.gid {
		nazalog.Out(nazalog.LevelError, 3, fmt.Sprintf("[%s] recursive lock. gid=%d", m.uniqueKey, gid))
	}
	m.gid = gid
	m.mu.Unlock()

	nazalog.Out(nazalog.LevelDebug, 3, fmt.Sprintf("[%s] > Lock(). gid=%d", m.uniqueKey, gid))
	m.core.Lock()
	nazalog.Out(nazalog.LevelDebug, 3, fmt.Sprintf("[%s] < Lock(). gid=%d", m.uniqueKey, gid))
}

func (m *Mutex) Unlock() {
	m.mu.Lock()
	gid, _ := CurGoroutineID()
	if gid != m.gid {
		if m.gid == unlockedGid {
			nazalog.Out(nazalog.LevelError, 3,
				fmt.Sprintf("[%s] unlock of unlocked mutex. lock gid=%d, unlock gid=%d", m.uniqueKey, m.gid, gid))
		} else {
			nazalog.Out(nazalog.LevelError, 3,
				fmt.Sprintf("[%s] mismatched unlock. lock gid=%d, unlock gid=%d", m.uniqueKey, m.gid, gid))
		}
	}
	m.gid = unlockedGid
	m.mu.Unlock()

	nazalog.Out(nazalog.LevelDebug, 3, fmt.Sprintf("[%s] > Unlock(). gid=%d", m.uniqueKey, gid))
	m.core.Unlock()
	nazalog.Out(nazalog.LevelDebug, 3, fmt.Sprintf("[%s] < Unlock(). gid=%d", m.uniqueKey, gid))
}

func init() {
	uniqueGen = unique.NewSingleGenerator("MUTEX")
}
