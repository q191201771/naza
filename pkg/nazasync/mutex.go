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
	"strings"
	"sync"
	"time"

	"github.com/q191201771/naza/pkg/unique"

	"github.com/q191201771/naza/pkg/nazalog"
)

// 用于debug锁方面的问题

var uniqueGen *unique.SingleGenerator

type Mutex struct {
	core          sync.Mutex
	startHoldTime time.Time

	genUniqueKeyOnce sync.Once
	uniqueKey        string
}

func (m *Mutex) Lock() {
	gid, _ := CurGoroutineId()

	m.genUniqueKeyOnce.Do(func() {
		m.uniqueKey = uniqueGen.GenUniqueKey()
	})

	b := time.Now()
	nazalog.Out(nazalog.LevelDebug, 3, fmt.Sprintf("[%s] > Lock(). gid=%d", m.uniqueKey, gid))
	globalMutexManager.beforeAcquireLock(m.uniqueKey, gid)
	m.core.Lock()
	globalMutexManager.afterAcquireLock(m.uniqueKey, gid)
	m.startHoldTime = time.Now()
	nazalog.Out(nazalog.LevelDebug, 3, fmt.Sprintf("[%s] < Lock(). gid=%d, acquire cost=%dms", m.uniqueKey, gid, time.Now().Sub(b).Milliseconds()))
}

func (m *Mutex) Unlock() {
	gid, _ := CurGoroutineId()

	// 运行时自己会检查对没有加锁的mutex进行Unlock调用的情况

	b := time.Now()
	nazalog.Out(nazalog.LevelDebug, 3, fmt.Sprintf("[%s] > Unlock(). gid=%d, hold cost=%dms", m.uniqueKey, gid, time.Now().Sub(m.startHoldTime).Milliseconds()))
	m.core.Unlock()
	globalMutexManager.afterUnlock(m.uniqueKey, gid)
	nazalog.Out(nazalog.LevelDebug, 3, fmt.Sprintf("[%s] < Unlock(). gid=%d, release cost=%dms", m.uniqueKey, gid, time.Now().Sub(b).Milliseconds()))
}

var globalMutexManager = NewMutexManager()

// 注意，key是由mutex唯一ID加上协程ID组合而成
type MutexManager struct {
	mu                   sync.Mutex
	waitAcquireContainer map[string]time.Time
	holdContainer        map[string]time.Time
}

func NewMutexManager() *MutexManager {
	m := &MutexManager{
		waitAcquireContainer: make(map[string]time.Time),
		holdContainer:        make(map[string]time.Time),
	}

	go m.printTmpDebug()

	return m
}

func (m *MutexManager) printTmpDebug() {
	var buf strings.Builder
	for {
		m.mu.Lock()
		now := time.Now()
		buf.Reset()
		buf.WriteString("long wait acquire:")
		for k, t := range m.waitAcquireContainer {
			duration := now.Sub(t).Milliseconds()
			if duration > 1000 {
				buf.WriteString(fmt.Sprintf(" (%s:%dms)", k, duration))
			}
		}
		nazalog.Out(nazalog.LevelDebug, 4, buf.String())
		buf.Reset()
		buf.WriteString("long hold:")
		for k, t := range m.holdContainer {
			duration := now.Sub(t).Milliseconds()
			if duration > 1000 {
				buf.WriteString(fmt.Sprintf(" (%s:%dms)", k, duration))
			}
		}
		nazalog.Out(nazalog.LevelDebug, 4, buf.String())
		m.mu.Unlock()

		time.Sleep(5 * time.Second)
	}
}

func (m *MutexManager) beforeAcquireLock(uk string, gid int64) {
	k := fmt.Sprintf("%s_%d", uk, gid)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exist := m.waitAcquireContainer[k]; exist {
		nazalog.Out(nazalog.LevelError, 4, fmt.Sprintf("[%s] one g try acquire lock twice(wait acquire already). gid=%d", uk, gid))
	}

	// 当前协程已持有锁，再次重入
	if _, exist := m.holdContainer[k]; exist {
		nazalog.Out(nazalog.LevelError, 4, fmt.Sprintf("[%s] one g try acquire lock twice(hold already). gid=%d", uk, gid))
	}

	m.waitAcquireContainer[k] = time.Now()
}

func (m *MutexManager) afterAcquireLock(uk string, gid int64) {
	k := fmt.Sprintf("%s_%d", uk, gid)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exist := m.waitAcquireContainer[k]; !exist {
		nazalog.Out(nazalog.LevelError, 4, fmt.Sprintf("[%s] acquired but not in wait acquire container. gid=%d", uk, gid))
	}
	if _, exist := m.holdContainer[k]; exist {
		nazalog.Out(nazalog.LevelError, 4, fmt.Sprintf("[%s] acquired but in  hold container already. gid=%d", uk, gid))
	}

	delete(m.waitAcquireContainer, k)
	m.holdContainer[k] = time.Now()
}

func (m *MutexManager) afterUnlock(uk string, gid int64) {
	k := fmt.Sprintf("%s_%d", uk, gid)
	m.mu.Lock()
	defer m.mu.Unlock()

	// 有可能是a协程Lock，b协程Unlock
	if _, exist := m.holdContainer[k]; !exist {
		nazalog.Out(nazalog.LevelError, 4, fmt.Sprintf("[%s] unlock but not in hold container. gid=%d", uk, gid))
	}

	if _, exist := m.waitAcquireContainer[k]; exist {
		nazalog.Out(nazalog.LevelError, 4, fmt.Sprintf("[%s] unlock but in wait acquire container already. gid=%d", uk, gid))
	}

	delete(m.holdContainer, k)
}

func init() {
	uniqueGen = unique.NewSingleGenerator("MUTEX")
}
