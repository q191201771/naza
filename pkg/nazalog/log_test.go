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
	"fmt"
	originLog "log"
	"os"
	"sync"
	"testing"
	"time"

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
	l.Tracef("l test msg by Trace%s", "f")
	l.Debugf("l test msg by Debug%s", "f")
	l.Infof("l test msg by Info%s", "f")
	l.Warnf("l test msg by Warn%s", "f")
	l.Errorf("l test msg by Error%s", "f")
	l.Trace("l test msg by Trace")
	l.Debug("l test msg by Debug")
	l.Info("l test msg by Info")
	l.Warn("l test msg by Warn")
	l.Error("l test msg by Error")
	l.Output(2, "l test msg by Output")
	l.Out(nazalog.LevelInfo, 1, "l test msg by Out")
	l.Print("l test msg by Print")
	l.Printf("l test msg by Print%s", "f")
	l.Println("l test msg by Print")
}

func TestGlobal(t *testing.T) {
	buf := []byte("1234567890987654321")
	nazalog.Error(hex.Dump(buf))
	nazalog.Tracef("g test msg by Trace%s", "f")
	nazalog.Debugf("g test msg by Debug%s", "f")
	nazalog.Infof("g test msg by Info%s", "f")
	nazalog.Warnf("g test msg by Warn%s", "f")
	nazalog.Errorf("g test msg by Error%s", "f")
	nazalog.Trace("g test msg by Trace")
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
	nazalog.Tracef("gc test msg by Trace%s", "f")
	nazalog.Debugf("gc test msg by Debug%s", "f")
	nazalog.Infof("gc test msg by Info%s", "f")
	nazalog.Warnf("gc test msg by Warn%s", "f")
	nazalog.Errorf("gc test msg by Error%s", "f")
	nazalog.Trace("gc test msg by Trace")
	nazalog.Debug("gc test msg by Debug")
	nazalog.Info("gc test msg by Info")
	nazalog.Warn("gc test msg by Warn")
	nazalog.Error("gc test msg by Error")
	nazalog.Output(2, "gc test msg by Output")
	nazalog.Out(nazalog.LevelInfo, 2, "gc test msg by Out")
	nazalog.Print("gc test msg by Print")
	nazalog.Printf("gc test msg by Print%s", "f")
	nazalog.Println("gc test msg by Print")
	nazalog.Sync()
}

func TestNew(t *testing.T) {
	var (
		l   nazalog.Logger
		err error
	)
	l, err = nazalog.New(func(option *nazalog.Option) {
		option.Level = nazalog.LevelLogNothing + 1
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
	nazalog.Info("aaa")
	fake.WithFakeTimeNow(func() time.Time {
		return time.Now().Add(48 * time.Hour)
	}, func() {
		nazalog.Info("bbb")
	})
}

func TestPanic(t *testing.T) {
	fake.WithRecover(func() {
		nazalog.Panic("aaa")
	})
	fake.WithRecover(func() {
		nazalog.Panicf("%s", "bbb")
	})
	fake.WithRecover(func() {
		nazalog.Panicln("aaa")
	})
	fake.WithRecover(func() {
		l, err := nazalog.New()
		assert.Equal(t, nil, err)
		l.Panic("aaa")
	})
	fake.WithRecover(func() {
		l, err := nazalog.New()
		assert.Equal(t, nil, err)
		l.Panicf("%s", "bbb")
	})
	fake.WithRecover(func() {
		l, err := nazalog.New()
		assert.Equal(t, nil, err)
		l.Panicln("aaa")
	})
}

func TestFatal(t *testing.T) {
	var er fake.ExitResult

	er = fake.WithFakeOSExit(func() {
		nazalog.Fatal("Fatal")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	er = fake.WithFakeOSExit(func() {
		nazalog.Fatalf("Fatalf%s", ".")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	er = fake.WithFakeOSExit(func() {
		nazalog.Fatalln("Fatalln")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	logger, err := nazalog.New(func(option *nazalog.Option) {
		option.Level = nazalog.LevelInfo
	})
	assert.IsNotNil(t, logger)
	assert.Equal(t, nil, err)
	er = fake.WithFakeOSExit(func() {
		logger.Fatal("Fatal")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	er = fake.WithFakeOSExit(func() {
		logger.Fatalf("Fatalf%s", ".")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	er = fake.WithFakeOSExit(func() {
		logger.Fatalln("Fatalln")
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
	err := fake.WithFakeOSExit(func() {
		nazalog.Assert(nil, 1)
	})
	assert.Equal(t, true, err.HasExit)
	assert.Equal(t, 1, err.ExitCode)

	_ = nazalog.Init(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertPanic
	})
	fake.WithRecover(func() {
		nazalog.Assert([]byte{}, "aaa")
	})

	l, _ := nazalog.New()
	l.Assert(nil, 1)

	l, _ = nazalog.New(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertFatal
	})
	err = fake.WithFakeOSExit(func() {
		l.Assert(nil, 1)
	})
	assert.Equal(t, true, err.HasExit)
	assert.Equal(t, 1, err.ExitCode)

	l, _ = nazalog.New(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertPanic
	})
	fake.WithRecover(func() {
		l.Assert([]byte{}, "aaa")
	})
}

func TestLogger_WithPrefix(t *testing.T) {
	im := 4
	jm := 4
	var wg sync.WaitGroup
	wg.Add(im * jm)
	nazalog.Debug(">")
	for i := 0; i != im; i++ {
		go func(ii int) {
			for j := 0; j != jm; j++ {
				s := fmt.Sprintf("%d", ii)
				l := nazalog.WithPrefix("log_test")
				l.Info(j)
				ll := l.WithPrefix("TestLogger_WithPrefix")
				ll.Info(j)
				lll := ll.WithPrefix(s)
				lll.Info(j)
				wg.Done()
			}
		}(i)
	}
	nazalog.Debug("<")
	wg.Wait()
}

func TestTimestamp(t *testing.T) {
	l, _ := nazalog.New(func(option *nazalog.Option) {
		option.TimestampFlag = false
	})
	l.Debug("without timestamp.")
	l.Info("without timestamp.")
	l, _ = nazalog.New(func(option *nazalog.Option) {
		option.TimestampWithMSFlag = false
	})
	l.Debug("without timestamp.")
	l.Info("timestamp without ms.")
}

func TestFieldFlag(t *testing.T) {
	l, _ := nazalog.New(func(option *nazalog.Option) {
	})
	l.Debug("1")
	l.Info("1")
	l, _ = nazalog.New(func(option *nazalog.Option) {
		option.ShortFileFlag = false
	})
	l.Debug("2")
	l.Info("2")
	l, _ = nazalog.New(func(option *nazalog.Option) {
		option.ShortFileFlag = false
		option.LevelFlag = false
	})
	l.Debug("3")
	l.Info("3")
	l, _ = nazalog.New(func(option *nazalog.Option) {
		option.ShortFileFlag = false
		option.LevelFlag = false
		option.TimestampFlag = false
	})
	l.Debug("4")
	l.Info("4")
}

func TestReadableString(t *testing.T) {
	assert.Equal(t, "LevelTrace", nazalog.LevelTrace.ReadableString())
	assert.Equal(t, "LevelDebug", nazalog.LevelDebug.ReadableString())
	assert.Equal(t, "LevelInfo", nazalog.LevelInfo.ReadableString())
	assert.Equal(t, "LevelWarn", nazalog.LevelWarn.ReadableString())
	assert.Equal(t, "LevelError", nazalog.LevelError.ReadableString())
	assert.Equal(t, "LevelFatal", nazalog.LevelFatal.ReadableString())
	assert.Equal(t, "LevelPanic", nazalog.LevelPanic.ReadableString())
	assert.Equal(t, "LevelLogNothing", nazalog.LevelLogNothing.ReadableString())
	assert.Equal(t, "unknown", nazalog.Level(100).ReadableString())

	assert.Equal(t, "AssertError", nazalog.AssertError.ReadableString())
	assert.Equal(t, "AssertFatal", nazalog.AssertFatal.ReadableString())
	assert.Equal(t, "AssertPanic", nazalog.AssertPanic.ReadableString())
	assert.Equal(t, "unknown", nazalog.AssertBehavior(100).ReadableString())
}

func TestLevel(t *testing.T) {
	var l nazalog.Logger
	l, _ = nazalog.New(func(option *nazalog.Option) {
		option.Level = nazalog.LevelTrace
	})
	l.Trace("log by trace")
	l.Debug("log by debug")
	l.Info("log by info")
	l.Warn("log by warn")
	l.Error("log by error")
}

func TestModGlobal(t *testing.T) {
	l, _ := nazalog.New(func(option *nazalog.Option) {
		option.Level = nazalog.LevelInfo
	})
	nazalog.SetGlobalLogger(l)
	nazalog.Debug("a")
	nazalog.Info("b")
	l = nazalog.GetGlobalLogger()
	l.Debug("c")
	l.Info("d")
}

func TestLogger_GetOption(t *testing.T) {
	l, _ := nazalog.New(func(option *nazalog.Option) {
		option.Level = nazalog.LevelDebug
	})
	o := l.GetOption()
	assert.Equal(t, nazalog.LevelDebug, o.Level)

	_ = nazalog.Init(func(option *nazalog.Option) {
		option.Level = nazalog.LevelDebug
	})
	o = nazalog.GetOption()
	assert.Equal(t, nazalog.LevelDebug, o.Level)
}

func BenchmarkNazaLog(b *testing.B) {
	b.ReportAllocs()

	err := nazalog.Init(func(option *nazalog.Option) {
		option.Level = nazalog.LevelInfo
		option.Filename = "/dev/null"
		option.IsToStdout = false
		option.IsRotateDaily = false
	})
	assert.Equal(b, nil, err)
	b.ResetTimer()
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		originLog.Printf("hello %s %d\n", "world", i)
		originLog.Println("Info")
	}
}
