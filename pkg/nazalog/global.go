// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazalog

import (
	"fmt"

	"github.com/q191201771/naza/pkg/nazareflect"

	"github.com/q191201771/naza/pkg/fake"
)

var global Logger

func Tracef(format string, v ...interface{}) {
	global.Out(LevelTrace, 2, fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	global.Out(LevelDebug, 2, fmt.Sprintf(format, v...))
}

func Infof(format string, v ...interface{}) {
	global.Out(LevelInfo, 2, fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}) {
	global.Out(LevelWarn, 2, fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	global.Out(LevelError, 2, fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...interface{}) {
	global.Out(LevelFatal, 2, fmt.Sprintf(format, v...))
	fake.OS_Exit(1)
}

func Panicf(format string, v ...interface{}) {
	global.Out(LevelPanic, 2, fmt.Sprintf(format, v...))
	panic(fmt.Sprintf(format, v...))
}

func Trace(v ...interface{}) {
	global.Out(LevelTrace, 2, fmt.Sprint(v...))
}

func Debug(v ...interface{}) {
	global.Out(LevelDebug, 2, fmt.Sprint(v...))
}

func Info(v ...interface{}) {
	global.Out(LevelInfo, 2, fmt.Sprint(v...))
}

func Warn(v ...interface{}) {
	global.Out(LevelWarn, 2, fmt.Sprint(v...))
}

func Error(v ...interface{}) {
	global.Out(LevelError, 2, fmt.Sprint(v...))
}

func Fatal(v ...interface{}) {
	global.Out(LevelFatal, 2, fmt.Sprint(v...))
	fake.OS_Exit(1)
}

func Panic(v ...interface{}) {
	global.Out(LevelPanic, 2, fmt.Sprint(v...))
	panic(fmt.Sprint(v...))
}

func Output(calldepth int, s string) error {
	global.Out(LevelInfo, calldepth, s)
	return nil
}

func Print(v ...interface{}) {
	global.Out(LevelInfo, 2, fmt.Sprint(v...))
}

func Printf(format string, v ...interface{}) {
	global.Out(LevelInfo, 2, fmt.Sprintf(format, v...))
}
func Println(v ...interface{}) {
	global.Out(LevelInfo, 2, fmt.Sprint(v...))
}
func Fatalln(v ...interface{}) {
	global.Out(LevelInfo, 2, fmt.Sprint(v...))
	fake.OS_Exit(1)
}
func Panicln(v ...interface{}) {
	global.Out(LevelInfo, 2, fmt.Sprint(v...))
	panic(fmt.Sprint(v...))
}

func Assert(expected interface{}, actual interface{}) {
	if !nazareflect.Equal(expected, actual) {
		err := fmt.Sprintf("assert failed. excepted=%+v, but actual=%+v", expected, actual)
		switch global.GetOption().AssertBehavior {
		case AssertError:
			global.Out(LevelError, 2, err)
		case AssertFatal:
			global.Out(LevelFatal, 2, err)
			fake.OS_Exit(1)
		case AssertPanic:
			global.Out(LevelPanic, 2, err)
			panic(err)
		}
	}
}

func Out(level Level, calldepth int, s string) {
	global.Out(level, calldepth, s)
}

func Sync() {
	global.Sync()
}

func WithPrefix(s string) Logger {
	return global.WithPrefix(s)
}

func GetOption() Option {
	return global.GetOption()
}

// -------------------------------------------------------------------------------------------------

// 这些设置全局Logger的函数没有加锁，由调用方保证调用Init、SetGlobalLogger函数时不会并发调用全局Logger的其他方法
// 也即最好在程序启动时配置

func SetGlobalLogger(l Logger) {
	global = l
}

func GetGlobalLogger() Logger {
	return global
}

// 初始化全局Logger
func Init(modOptions ...ModOption) error {
	var err error
	global, err = newLogger(modOptions...)
	return err
}

func init() {
	_ = Init()
}
