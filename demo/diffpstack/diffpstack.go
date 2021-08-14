// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/q191201771/naza/pkg/nazalog"
)

// 分析c++程序pstack的两次结果的差异

type PstackInfo struct {
	tis []ThreadInfo
	tim map[string]ThreadInfo
}

type ThreadInfo struct {
	Num int
	P   string
	Id  int

	RawLine       string
	RawStackLines string
}

func NewPstackInfo(filename string) PstackInfo {
	var tis []ThreadInfo

	contents, err := ioutil.ReadFile(filename)
	nazalog.Assert(nil, err)
	content := string(contents)
	lines := strings.Split(content, "\n")
	//nazalog.Debugf("len(lines)=%d", len(lines))

	var ti *ThreadInfo
	for _, line := range lines {
		if strings.HasPrefix(line, "Thread") {
			if ti != nil {
				tis = append(tis, *ti)
			}
			ti = &ThreadInfo{}

			//nazalog.Debugf("%s", line)
			ti.RawLine = line
			ti.Num, ti.P, ti.Id, err = parseThreadLine(line)
			nazalog.Assert(nil, err)
			continue
		}

		ti.RawStackLines += line + "\n"
	}
	if ti != nil {
		tis = append(tis, *ti)
	}
	nazalog.Debugf("len(tis)=%d", len(tis))

	tim := make(map[string]ThreadInfo)
	for _, ti := range tis {
		tim[ti.Uk()] = ti
	}

	return PstackInfo{
		tis: tis,
		tim: tim,
	}
}

func (pi *PstackInfo) Find(uk string) (ThreadInfo, bool) {
	ti, exist := pi.tim[uk]
	return ti, exist
}

func (ti *ThreadInfo) Uk() string {
	return fmt.Sprintf("%s_%d", ti.P, ti.Id)
}

func parseThreadLine(line string) (num int, p string, id int, err error) {
	p1 := strings.Index(line, "Thread")
	p2 := strings.Index(line, "(Thread")
	num, err = strconv.Atoi(line[p1+7 : p2-1])
	if err != nil {
		return
	}

	p3 := strings.Index(line, "(LWP")
	p = line[p2+8 : p3-1]

	p4 := strings.Index(line, "))")
	id, err = strconv.Atoi(line[p3+5 : p4-1])

	return
}

func main() {
	_ = nazalog.Init(func(option *nazalog.Option) {
		option.LevelFlag = false
		option.ShortFileFlag = false
		option.TimestampFlag = false
	})

	pi1 := NewPstackInfo("old.txt")
	pi2 := NewPstackInfo("new.txt")
	for _, ti2 := range pi2.tis {
		var pre, suf string
		suf = "\033[0m"
		ti1, exist := pi1.Find(ti2.Uk())
		if exist {
			if ti2.RawStackLines == ti1.RawStackLines {
				// 1, 2都有，但是堆栈没变 没色
				pre = ""
				suf = ""
				nazalog.Debugf("%s-------------------------------------------------------------------------%s", pre, suf)
				nazalog.Debugf("%s%s%s", pre, ti2.RawLine, suf)
				nazalog.Debugf("%s%s%s", pre, ti2.RawStackLines, suf)
			} else {
				// 1, 2都有，但是堆栈变化 红色
				// 注意，断站有变化也可能是函数参数变化了
				pre = "\033[22;31m"
				nazalog.Debugf("%s-------------------------------------------------------------------------%s", pre, suf)
				nazalog.Debugf("%s%s%s", pre, ti2.RawLine, suf)
				nazalog.Debugf("%s%s%s", pre, ti2.RawStackLines, suf)
			}
		} else {
			// 只在2有 绿色
			pre = "\033[22;36m"
			nazalog.Debugf("%s-------------------------------------------------------------------------%s", pre, suf)
			nazalog.Debugf("%s%s%s", pre, ti2.RawLine, suf)
			nazalog.Debugf("%s%s%s", pre, ti2.RawStackLines, suf)
		}
	}
}
