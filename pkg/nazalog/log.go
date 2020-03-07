// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazalog

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/q191201771/naza/pkg/fake"
)

const (
	levelDebugString = "DEBUG "
	levelInfoString  = " INFO "
	levelWarnString  = " WARN "
	levelErrorString = "ERROR "
	levelFatalString = "FATAL "
	levelPanicString = "PANIC "

	levelDebugColorString = "\033[22;37mDEBUG\033[0m "
	levelInfoColorString  = "\033[22;36m INFO\033[0m "
	levelWarnColorString  = "\033[22;33m WARN\033[0m "
	levelErrorColorString = "\033[22;31mERROR\033[0m "
	levelFatalColorString = "\033[22;31mFATAL\033[0m " // 颜色和 error 级别一样
	levelPanicColorString = "\033[22;31mPANIC\033[0m " // 颜色和 error 级别一样
)

var (
	levelToString = map[Level]string{
		LevelDebug: levelDebugString,
		LevelInfo:  levelInfoString,
		LevelWarn:  levelWarnString,
		LevelError: levelErrorString,
		LevelFatal: levelFatalString,
		LevelPanic: levelPanicString,
	}
	levelToColorString = map[Level]string{
		LevelDebug: levelDebugColorString,
		LevelInfo:  levelInfoColorString,
		LevelWarn:  levelWarnColorString,
		LevelError: levelErrorColorString,
		LevelFatal: levelFatalColorString,
		LevelPanic: levelPanicColorString,
	}
)

type logger struct {
	option Option

	dir string

	m             sync.Mutex
	fp            *os.File
	console       *os.File
	buf           bytes.Buffer
	currRoundTime time.Time
}

func (l *logger) Outputf(level Level, calldepth int, format string, v ...interface{}) {
	l.Out(level, 3, fmt.Sprintf(format, v...))
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.Out(LevelDebug, 3, fmt.Sprintf(format, v...))
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Out(LevelInfo, 3, fmt.Sprintf(format, v...))
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Out(LevelWarn, 3, fmt.Sprintf(format, v...))
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Out(LevelError, 3, fmt.Sprintf(format, v...))
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Out(LevelFatal, 3, fmt.Sprintf(format, v...))
	fake.Exit(1)
}

func (l *logger) Panicf(format string, v ...interface{}) {
	l.Out(LevelPanic, 3, fmt.Sprintf(format, v...))
	panic(fmt.Sprintf(format, v...))
}

func (l *logger) Output(level Level, calldepth int, v ...interface{}) {
	l.Out(level, 3, fmt.Sprint(v...))
}

func (l *logger) Debug(v ...interface{}) {
	l.Out(LevelDebug, 3, fmt.Sprint(v...))
}

func (l *logger) Info(v ...interface{}) {
	l.Out(LevelInfo, 3, fmt.Sprint(v...))
}

func (l *logger) Warn(v ...interface{}) {
	l.Out(LevelWarn, 3, fmt.Sprint(v...))
}

func (l *logger) Error(v ...interface{}) {
	l.Out(LevelError, 3, fmt.Sprint(v...))
}

func (l *logger) Fatal(v ...interface{}) {
	l.Out(LevelFatal, 3, fmt.Sprint(v...))
	fake.Exit(1)
}

func (l *logger) Panic(v ...interface{}) {
	l.Out(LevelPanic, 3, fmt.Sprint(v...))
	panic(fmt.Sprint(v...))
}

func (l *logger) FatalIfErrorNotNil(err error) {
	if err != nil {
		l.Out(LevelError, 3, fmt.Sprintf("fatal since error not nil. err=%+v", err))
		fake.Exit(1)
	}
}

func (l *logger) PanicIfErrorNotNil(err error) {
	if err != nil {
		l.Out(LevelPanic, 3, fmt.Sprintf("panic since error not nil. err=%+v", err))
		panic(err)
	}
}

func (l *logger) Out(level Level, calldepth int, s string) {
	if l.option.Level > level {
		return
	}

	l.m.Lock()
	defer l.m.Unlock()

	now := time.Now()

	// 格式化日志内容
	l.buf.Reset()
	writeTime(&l.buf, now)
	if l.console != nil {
		l.buf.WriteString(levelToColorString[level])
	} else {
		l.buf.WriteString(levelToString[level])
	}
	l.buf.WriteString(s)
	if l.option.ShortFileFlag {
		writeShortFile(&l.buf, calldepth)
	}
	if l.buf.Len() == 0 || l.buf.Bytes()[l.buf.Len()-1] != '\n' {
		l.buf.WriteByte('\n')
	}

	// 输出至控制台
	if l.console != nil {
		_, _ = l.console.Write(l.buf.Bytes())
		if level == LevelFatal || level == LevelPanic {
			_ = l.console.Sync()
		}
	}

	// 输出至日志文件
	if l.fp != nil {
		if l.option.IsRotateDaily && now.Day() != l.currRoundTime.Day() {
			backupName := l.option.Filename + "." + l.currRoundTime.Format("20060102")
			if err := os.Rename(l.option.Filename, backupName); err == nil {
				_ = l.fp.Close()
				l.fp, _ = os.Create(l.option.Filename)
			}
			l.currRoundTime = now
		}
		_, _ = l.fp.Write(l.buf.Bytes())
		if level == LevelFatal || level == LevelPanic {
			_ = l.fp.Sync()
		}
	}
}

func (l *logger) Sync() {
	l.m.Lock()
	defer l.m.Unlock()

	if l.console != nil {
		_ = l.console.Sync()
	}
	if l.fp != nil {
		_ = l.fp.Sync()
	}
}

func writeTime(buf *bytes.Buffer, t time.Time) {
	year, month, day := t.Date()
	itoa(buf, year, 4)
	buf.WriteByte('/')
	itoa(buf, int(month), 2)
	buf.WriteByte('/')
	itoa(buf, day, 2)
	buf.WriteByte(' ')

	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	buf.WriteByte(':')
	itoa(buf, min, 2)
	buf.WriteByte(':')
	itoa(buf, sec, 2)
	buf.WriteByte('.')
	itoa(buf, t.Nanosecond()/1e3, 6)
	buf.WriteByte(' ')
}

func writeShortFile(buf *bytes.Buffer, calldepth int) {
	buf.Write([]byte{' ', '-', ' '})

	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	buf.WriteString(file)
	buf.WriteByte(':')
	itoa(buf, line, -1)
}

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// @NOTICE 该函数拷贝自 Go 标准库 /src/log/log.go: func itoa(buf *[]byte, i int, wid int)
// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *bytes.Buffer, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	buf.Write(b[bp:])
}
