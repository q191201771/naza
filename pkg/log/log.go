// package log 日志库
package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// 1. 带日志级别
// 2. 可选输出至控制台或文件，也可以同时输出
// 3. 日志文件支持按天翻转
//
// 目前性能和标准库log相当

var LogErr = errors.New("log:fxxk")

type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{}) // 打印日志并退出程序
	Panicf(format string, v ...interface{})

	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Panic(v ...interface{})

	FatalIfErrorNotNil(err error)
	PanicIfErrorNotNil(err error)

	Outputf(level Level, calldepth int, format string, v ...interface{})
	Output(level Level, calldepth int, v ...interface{})
	Out(level Level, calldepth int, s string)
}

type Config struct {
	Level Level `json:"level"` // 日志级别，大于等于该级别的日志才会被输出

	// 文件输出和控制台输出可同时打开
	// 控制台输出主要用做开发时调试，打开后level字段使用彩色输出
	Filename   string `json:"filename"`     // 输出日志文件名，如果为空，则不写日志文件。可包含路径，路径不存在时，将自动创建
	IsToStdout bool   `json:"is_to_stdout"` // 是否以stdout输出到控制台

	IsRotateDaily bool `json:"is_rotate_daily"` // 日志按天翻转

	ShortFileFlag bool `json:"short_file_flag"` // 是否在每行日志尾部添加源码文件及行号的信息
}

type Level uint8

const (
	_ Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func New(c Config) (Logger, error) {
	var (
		dir     string
		fp      *os.File
		console io.Writer
		err     error
	)
	if c.Level < LevelDebug || c.Level > LevelPanic {
		return nil, LogErr
	}
	if c.Filename != "" {
		dir = filepath.Dir(c.Filename)
		if err := os.MkdirAll(dir, 0644); err != nil {
			return nil, err
		}
		fp, err = os.Create(c.Filename)
		if err != nil {
			return nil, err
		}
	}
	if c.IsToStdout {
		console = os.Stdout
	}

	l := &logger{
		c:             c,
		dir:           dir,
		fp:            fp,
		console:       console,
		currRoundTime: time.Now(),
	}
	return l, nil
}

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
	c Config

	dir string

	m             sync.Mutex
	fp            *os.File
	console       io.Writer
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
	os.Exit(1)
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
	os.Exit(1)
}

func (l *logger) Panic(v ...interface{}) {
	l.Out(LevelPanic, 3, fmt.Sprint(v...))
	panic(fmt.Sprint(v...))
}

func (l *logger) FatalIfErrorNotNil(err error) {
	if err != nil {
		l.Out(LevelError, 3, fmt.Sprintf("fatal since error not nil. err=%+v", err))
		os.Exit(1)
	}
}

func (l *logger) PanicIfErrorNotNil(err error) {
	if err != nil {
		l.Out(LevelPanic, 3, fmt.Sprintf("panic since error not nil. err=%+v", err))
		panic(err)
	}
}

func (l *logger) Out(level Level, calldepth int, s string) {
	if l.c.Level > level {
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
	if l.c.ShortFileFlag {
		writeShortFile(&l.buf, calldepth)
	}
	if l.buf.Len() == 0 || l.buf.Bytes()[l.buf.Len()-1] != '\n' {
		l.buf.WriteByte('\n')
	}

	// 输出至控制台
	if l.console != nil {
		_, _ = l.console.Write(l.buf.Bytes())
	}

	// 输出至日志文件
	if l.fp != nil {
		if now.Day() != l.currRoundTime.Day() {
			backupName := l.c.Filename + "." + l.currRoundTime.Format("20060102")
			if err := os.Rename(l.c.Filename, backupName); err == nil {
				_ = l.fp.Close()
				l.fp, _ = os.Create(l.c.Filename)
			}
			l.currRoundTime = now
		}
		_, _ = l.fp.Write(l.buf.Bytes())
	}
}

// @NOTICE 该函数拷贝自 Go 标准库
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
