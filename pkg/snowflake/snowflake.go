// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package snowflake

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrInitial = errors.New("lal.snowflake: initial error")
	ErrGen     = errors.New("lal.snowflake: gen error")
)

type Option struct {
	DataCenterIDBits int   // 数据中心编号字段在所生成 ID 所占的位数，取值范围见 validate 函数
	WorkerIDBits     int   // 节点编号
	SequenceBits     int   // 递增序列
	Twepoch          int64 // 基准时间点
	AlwaysPositive   bool  // 是否只生成正数 ID，如果是，则时间戳所占位数会减少1位
}

var defaultOption = Option{
	DataCenterIDBits: 5,
	WorkerIDBits:     5,
	SequenceBits:     12,
	Twepoch:          int64(1288834974657), // 对应现实时间： 2010/11/4 9:42:54.657
	AlwaysPositive:   false,
}

type Node struct {
	dataCenterID int64
	workerID     int64
	option       Option

	seqMask           uint32
	workerIDShift     uint32
	dataCenterIDShift uint32
	timestampShift    uint32

	mu     sync.Mutex
	lastTs int64
	seq    uint32
}

type ModOption func(option *Option)

// dataCenterID 和 workerID 的取值范围取决于 DataCenterIDBits 和 WorkerIDBits
// 假设 DataCenterIDBits 为 5，则 dataCenterID 取值范围为 [0, 32]
func New(dataCenterID int, workerID int, modOptions ...ModOption) (*Node, error) {
	option := defaultOption
	for _, fn := range modOptions {
		fn(&option)
	}

	if err := validate(dataCenterID, workerID, option); err != nil {
		return nil, err
	}

	return &Node{
		dataCenterID:      int64(dataCenterID),
		workerID:          int64(workerID),
		option:            option,
		seqMask:           uint32(bitsToMax(option.SequenceBits)),
		workerIDShift:     uint32(option.SequenceBits),
		dataCenterIDShift: uint32(option.SequenceBits + option.WorkerIDBits),
		timestampShift:    uint32(option.SequenceBits + option.WorkerIDBits + option.DataCenterIDBits),
		lastTs:            -1,
	}, nil
}

func (n *Node) Gen(nowUnixMSec ...int64) (int64, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// 当前 Unix 时间戳可由外部传入
	var now int64
	if len(nowUnixMSec) == 0 {
		now = time.Now().UnixNano() / 1e6
	} else {
		now = nowUnixMSec[0]
	}

	// 时间戳回退，返回错误
	if now < n.lastTs {
		return -1, ErrGen
	}

	// 时间戳相同时，使用递增序号解决冲突
	if now == n.lastTs {
		n.seq = (n.seq + 1) & n.seqMask
		// 递增序号翻转为 0，表示该时间戳下的序号已经全部用完，阻塞等待系统时间增长
		if n.seq == 0 {
			for now <= n.lastTs {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		n.seq = 0
	}
	n.lastTs = now

	// 如果保证只返回正数，则生成的 ID 的最高位，也即时间戳的最高位保持为 0
	ts := now - n.option.Twepoch
	if n.option.AlwaysPositive {
		ts = clearBit(ts, 63-n.timestampShift)
	}
	ts <<= n.timestampShift

	// 用所有字段组合生成 ID 返回
	return ts | (n.dataCenterID << n.dataCenterIDShift) | (n.workerID << n.workerIDShift) | int64(n.seq), nil
}

func validate(dataCenterID int, workerID int, option Option) error {
	if option.DataCenterIDBits < 0 || option.DataCenterIDBits > 31 {
		return ErrInitial
	}
	if option.WorkerIDBits < 0 || option.WorkerIDBits > 31 {
		return ErrInitial
	}
	if option.SequenceBits < 0 || option.SequenceBits > 31 {
		return ErrInitial
	}

	if option.DataCenterIDBits+option.WorkerIDBits+option.SequenceBits >= 64 {
		return ErrInitial
	}

	if option.DataCenterIDBits > 0 {
		if dataCenterID > bitsToMax(option.DataCenterIDBits) {
			return ErrInitial
		}
	}
	if option.WorkerIDBits > 0 {
		if workerID > bitsToMax(option.WorkerIDBits) {
			return ErrInitial
		}
	}

	return nil
}

// 位的数量对应的最大值，该函数也可以叫做 bitsToMask
func bitsToMax(bits int) int {
	// -1 表示所有位都为 1
	return int(int32(-1) ^ (int32(-1) << uint32(bits)))
}

// 将 <num> 的第 <index> 设置为 0
func clearBit(num int64, index uint32) int64 {
	bit := int64(1 << index)
	mask := int64(-1) ^ bit
	return num & mask
}
