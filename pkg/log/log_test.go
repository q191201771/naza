package log

import (
	"github.com/q191201771/nezha/pkg/assert"
	originLog "log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	c := Config{
		Level:       LevelInfo,
		Filename:    "/tmp/lallogtest/aaa.log",
		IsToStdout:  true,
		IsRotateDaily: true,
	}
	l, err := New(c)
	assert.Equal(t, nil, err)
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
	Debugf("g test msg by Debug%s", "f")
	Infof("g test msg by Info%s", "f")
	Warnf("g test msg by Warn%s", "f")
	Errorf("g test msg by Error%s", "f")
	Debug("g test msg by Debug")
	Info("g test msg by Info")
	Warn("g test msg by Warn")
	Error("g test msg by Error")

	c := Config{
		Level:       LevelInfo,
		Filename:    "/tmp/lallogtest/bbb.log",
		IsToStdout:  true,
	}
	err := Init(c)
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
	l, err = New(Config{Level: LevelFatal + 1})
	assert.Equal(t, nil, l)
	assert.Equal(t, LogErr, err)

	l, err = New(Config{Filename: "/tmp"})
	assert.Equal(t, nil, l)
	assert.IsNotNil(t, err)

	l, err = New(Config{Filename: "./log_test.go/111"})
	assert.Equal(t, nil, l)
	assert.IsNotNil(t, err)
}

func TestRotate(t *testing.T) {
	c := Config{
		Level:       LevelInfo,
		Filename:    "/tmp/lallogtest/ccc.log",
		IsToStdout:  false,
		IsRotateDaily: true,
	}
	err := Init(c)
	assert.Equal(t, nil, err)
	b := make([]byte, 1024)
	for i := 0; i < 2*1024; i++ {
		Info(b)
	}
	for i := 0; i < 2*1024; i++ {
		Infof("%+v", b)
	}
}

func BenchmarkStdout(b *testing.B) {
	b.ReportAllocs()
	c := Config{
		Level:    LevelInfo,
		//Filename: "/tmp/lallogtest/ddd.log",
		Filename:    "/dev/null",
		//IsToStdout:  true,
		ShortFileFlag:true,
	}
	err := Init(c)
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
