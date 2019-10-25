// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazalog

import (
	"encoding/hex"
	"errors"
	originLog "log"
	"os"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
)

func TestLogger(t *testing.T) {
	l, err := New(func(option *Option) {
		option.Level = LevelInfo
		option.Filename = "/tmp/nazalogtest/aaa.log"
		option.IsToStdout = true
		option.IsRotateDaily = true
	})
	assert.Equal(t, nil, err)
	buf := []byte("1234567890987654321")
	l.Error(hex.Dump(buf))
	l.Debugf("l test msg by Debug%s", "f")
	l.Infof("l test msg by Info%s", "f")
	l.Warnf("l test msg by Warn%s", "f")
	l.Errorf("l test msg by Error%s", "f")
	l.Debug("l test msg by Debug")
	l.Info("l test msg by Info")
	l.Warn("l test msg by Warn")
	l.Error("l test msg by Error")
	l.Outputf(LevelInfo, 3, "l test msg by Output%s", "f")
	l.Output(LevelInfo, 3, "l test msg by Output")
	l.Out(LevelInfo, 3, "l test msg by Out")
}

func TestGlobal(t *testing.T) {
	buf := []byte("1234567890987654321")
	Error(hex.Dump(buf))
	Debugf("g test msg by Debug%s", "f")
	Infof("g test msg by Info%s", "f")
	Warnf("g test msg by Warn%s", "f")
	Errorf("g test msg by Error%s", "f")
	Debug("g test msg by Debug")
	Info("g test msg by Info")
	Warn("g test msg by Warn")
	Error("g test msg by Error")

	err := Init(func(option *Option) {
		option.Level = LevelInfo
		option.Filename = "/tmp/nazalogtest/bbb.log"
		option.IsToStdout = true

	})
	assert.Equal(t, nil, err)
	Debugf("gc test msg by Debug%s", "f")
	Infof("gc test msg by Info%s", "f")
	Warnf("gc test msg by Warn%s", "f")
	Errorf("gc test msg by Error%s", "f")
	Debug("gc test msg by Debug")
	Info("gc test msg by Info")
	Warn("gc test msg by Warn")
	Error("gc test msg by Error")
	Outputf(LevelInfo, 3, "gc test msg by Output%s", "f")
	Output(LevelInfo, 3, "gc test msg by Output")
	Out(LevelInfo, 3, "gc test msg by Out")
}

func TestNew(t *testing.T) {
	var (
		l   Logger
		err error
	)
	l, err = New(func(option *Option) {
		option.Level = LevelPanic + 1
	})
	assert.Equal(t, nil, l)
	assert.Equal(t, ErrLog, err)

	l, err = New(func(option *Option) {
		option.Filename = "/tmp"
	})
	assert.Equal(t, nil, l)
	assert.IsNotNil(t, err)

	l, err = New(func(option *Option) {
		option.Filename = "./log_test.go/111"
	})
	assert.Equal(t, nil, l)
	assert.IsNotNil(t, err)
}

func TestRotate(t *testing.T) {
	err := Init(func(option *Option) {
		option.Level = LevelInfo
		option.Filename = "/tmp/nazalogtest/ccc.log"
		option.IsToStdout = false
		option.IsRotateDaily = true

	})
	assert.Equal(t, nil, err)
	b := make([]byte, 1024)
	for i := 0; i < 2*1024; i++ {
		Info(b)
	}
	for i := 0; i < 2*1024; i++ {
		Infof("%+v", b)
	}
}

func withRecover(f func()) {
	defer func() {
		recover()
	}()
	f()
}

func TestPanic(t *testing.T) {
	withRecover(func() {
		Debug("ddd")
		Panic("aaa")
	})
	withRecover(func() {
		Panicf("%s", "bbb")
	})
	withRecover(func() {
		PanicIfErrorNotNil(errors.New("mock error"))
	})
	withRecover(func() {
		l, err := New()
		assert.Equal(t, nil, err)
		l.Panic("aaa")
	})
	withRecover(func() {
		l, err := New()
		assert.Equal(t, nil, err)
		l.Panicf("%s", "bbb")
	})
	withRecover(func() {
		l, err := New()
		assert.Equal(t, nil, err)
		l.PanicIfErrorNotNil(errors.New("mock error"))
	})
}

func BenchmarkStdout(b *testing.B) {
	b.ReportAllocs()

	err := Init(func(option *Option) {
		option.Level = LevelInfo
		option.Filename = "/dev/null"
	})
	assert.Equal(b, nil, err)
	for i := 0; i < b.N; i++ {
		Infof("hello %s %d", "world", i)
		Info("Info")
	}
}

func BenchmarkOriginLog(b *testing.B) {
	b.ReportAllocs()
	fp, err := os.Create("/dev/null")
	assert.Equal(b, nil, err)
	originLog.SetOutput(fp)
	originLog.SetFlags(originLog.Ldate | originLog.Ltime | originLog.Lmicroseconds | originLog.Lshortfile)
	for i := 0; i < b.N; i++ {
		originLog.Printf("hello %s %d\n", "world", i)
		originLog.Println("Info")
	}
}
