// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazalog_test

import (
	"encoding/hex"
	"errors"
	originLog "log"
	"os"
	"testing"

	"github.com/q191201771/naza/pkg/nazalog"

	"github.com/q191201771/naza/pkg/fake"

	"github.com/q191201771/naza/pkg/assert"
)

func TestLogger(t *testing.T) {
	l, err := nazalog.New(func(option *nazalog.Option) {
		option.Level = nazalog.LevelInfo
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
	l.Outputf(nazalog.LevelInfo, 3, "l test msg by Output%s", "f")
	l.Output(nazalog.LevelInfo, 3, "l test msg by Output")
	l.Out(nazalog.LevelInfo, 2, "l test msg by Out")
}

func TestGlobal(t *testing.T) {
	buf := []byte("1234567890987654321")
	nazalog.Error(hex.Dump(buf))
	nazalog.Debugf("g test msg by Debug%s", "f")
	nazalog.Infof("g test msg by Info%s", "f")
	nazalog.Warnf("g test msg by Warn%s", "f")
	nazalog.Errorf("g test msg by Error%s", "f")
	nazalog.Debug("g test msg by Debug")
	nazalog.Info("g test msg by Info")
	nazalog.Warn("g test msg by Warn")
	nazalog.Error("g test msg by Error")

	err := nazalog.Init(func(option *nazalog.Option) {
		option.Level = nazalog.LevelInfo
		option.Filename = "/tmp/nazalogtest/bbb.log"
		option.IsToStdout = true

	})
	assert.Equal(t, nil, err)
	nazalog.Debugf("gc test msg by Debug%s", "f")
	nazalog.Infof("gc test msg by Info%s", "f")
	nazalog.Warnf("gc test msg by Warn%s", "f")
	nazalog.Errorf("gc test msg by Error%s", "f")
	nazalog.Debug("gc test msg by Debug")
	nazalog.Info("gc test msg by Info")
	nazalog.Warn("gc test msg by Warn")
	nazalog.Error("gc test msg by Error")
	nazalog.Outputf(nazalog.LevelInfo, 3, "gc test msg by Output%s", "f")
	nazalog.Output(nazalog.LevelInfo, 3, "gc test msg by Output")
	nazalog.Out(nazalog.LevelInfo, 3, "gc test msg by Out")
	nazalog.Sync()
}

func TestNew(t *testing.T) {
	var (
		l   nazalog.Logger
		err error
	)
	l, err = nazalog.New(func(option *nazalog.Option) {
		option.Level = nazalog.LevelPanic + 1
	})
	assert.Equal(t, nil, l)
	assert.Equal(t, nazalog.ErrLog, err)

	l, err = nazalog.New(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertPanic + 1
	})
	assert.Equal(t, nil, l)
	assert.Equal(t, nazalog.ErrLog, err)

	l, err = nazalog.New(func(option *nazalog.Option) {
		option.Filename = "/tmp"
	})
	assert.Equal(t, nil, l)
	assert.IsNotNil(t, err)

	l, err = nazalog.New(func(option *nazalog.Option) {
		option.Filename = "./log_test.go/111"
	})
	assert.Equal(t, nil, l)
	assert.IsNotNil(t, err)
}

func TestRotate(t *testing.T) {
	err := nazalog.Init(func(option *nazalog.Option) {
		option.Level = nazalog.LevelInfo
		option.Filename = "/tmp/nazalogtest/ccc.log"
		option.IsToStdout = false
		option.IsRotateDaily = true

	})
	assert.Equal(t, nil, err)
	b := make([]byte, 1024)
	for i := 0; i < 2*1024; i++ {
		nazalog.Info(b)
	}
	for i := 0; i < 2*1024; i++ {
		nazalog.Infof("%+v", b)
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
		nazalog.Debug("ddd")
		nazalog.Panic("aaa")
	})
	withRecover(func() {
		nazalog.Panicf("%s", "bbb")
	})
	withRecover(func() {
		nazalog.PanicIfErrorNotNil(errors.New("mock error"))
	})
	withRecover(func() {
		l, err := nazalog.New()
		assert.Equal(t, nil, err)
		l.Panic("aaa")
	})
	withRecover(func() {
		l, err := nazalog.New()
		assert.Equal(t, nil, err)
		l.Panicf("%s", "bbb")
	})
	withRecover(func() {
		l, err := nazalog.New()
		assert.Equal(t, nil, err)
		l.PanicIfErrorNotNil(errors.New("mock error"))
	})
}

func TestFatal(t *testing.T) {
	var er fake.ExitResult
	er = fake.WithFakeExit(func() {
		nazalog.FatalIfErrorNotNil(errors.New("fxxk"))
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)
	er = fake.WithFakeExit(func() {
		nazalog.Fatal("Fatal")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)
	er = fake.WithFakeExit(func() {
		nazalog.Fatalf("Fatalf%s", ".")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	logger, _ := nazalog.New()
	er = fake.WithFakeExit(func() {
		logger.FatalIfErrorNotNil(errors.New("fxxk"))
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)
	er = fake.WithFakeExit(func() {
		logger.Fatal("Fatal")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)
	er = fake.WithFakeExit(func() {
		logger.Fatalf("Fatalf%s", ".")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)
}

func TestAssert(t *testing.T) {
	// 成功
	nazalog.Assert(nil, nil)
	nazalog.Assert(nil, nil)
	nazalog.Assert(nil, nil)
	nazalog.Assert(1, 1)
	nazalog.Assert("aaa", "aaa")
	var ch chan struct{}
	nazalog.Assert(nil, ch)
	var m map[string]string
	nazalog.Assert(nil, m)
	var p *int
	nazalog.Assert(nil, p)
	var i interface{}
	nazalog.Assert(nil, i)
	var b []byte
	nazalog.Assert(nil, b)

	nazalog.Assert([]byte{}, []byte{})
	nazalog.Assert([]byte{0, 1, 2}, []byte{0, 1, 2})

	// 失败
	_ = nazalog.Init(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertError
	})
	nazalog.Assert(nil, 1)

	_ = nazalog.Init(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertFatal
	})
	err := fake.WithFakeExit(func() {
		nazalog.Assert(nil, 1)
	})
	assert.Equal(t, true, err.HasExit)
	assert.Equal(t, 1, err.ExitCode)

	_ = nazalog.Init(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertPanic
	})
	withRecover(func() {
		nazalog.Assert([]byte{}, "aaa")
	})

	l, _ := nazalog.New()
	l.Assert(nil, 1)

	l, _ = nazalog.New(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertFatal
	})
	err = fake.WithFakeExit(func() {
		l.Assert(nil, 1)
	})
	assert.Equal(t, true, err.HasExit)
	assert.Equal(t, 1, err.ExitCode)

	l, _ = nazalog.New(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertPanic
	})
	withRecover(func() {
		l.Assert([]byte{}, "aaa")
	})
}

func BenchmarkStdout(b *testing.B) {
	b.ReportAllocs()

	err := nazalog.Init(func(option *nazalog.Option) {
		option.Level = nazalog.LevelInfo
		option.Filename = "/dev/null"
	})
	assert.Equal(b, nil, err)
	for i := 0; i < b.N; i++ {
		nazalog.Infof("hello %s %d", "world", i)
		nazalog.Info("Info")
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
