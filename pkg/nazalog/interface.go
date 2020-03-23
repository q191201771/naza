// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// package nazalog 日志库
package nazalog

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

// 1. 带日志级别
// 2. 可选输出至控制台或文件，也可以同时输出
// 3. 日志文件支持按天翻转
//
// 目前性能和标准库log相当

var ErrLog = errors.New("naza.log:fxxk")

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

	Assert(expected interface{}, actual interface{}) // 不相等时打印error级别日志
	FatalAssert(expected interface{}, actual interface{})
	PanicAssert(expected interface{}, actual interface{})

	Outputf(level Level, calldepth int, format string, v ...interface{})
	Output(level Level, calldepth int, v ...interface{})
	Out(level Level, calldepth int, s string)

	Sync()
}

type Option struct {
	Level Level `json:"level"` // 日志级别，大于等于该级别的日志才会被输出

	// 文件输出和控制台输出可同时打开
	// 控制台输出主要用做开发时调试，打开后level字段使用彩色输出
	Filename   string `json:"filename"`     // 输出日志文件名，如果为空，则不写日志文件。可包含路径，路径不存在时，将自动创建
	IsToStdout bool   `json:"is_to_stdout"` // 是否以stdout输出到控制台

	IsRotateDaily bool `json:"is_rotate_daily"` // 日志按天翻转

	ShortFileFlag bool `json:"short_file_flag"` // 是否在每行日志尾部添加源码文件及行号的信息
}

// 没有配置的属性，将按如下配置
var defaultOption = Option{
	Level:         LevelDebug,
	Filename:      "",
	IsToStdout:    true,
	IsRotateDaily: false,
	ShortFileFlag: true,
}

type Level uint8

const (
	_          Level = iota
	LevelDebug       // 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

type ModOption func(option *Option)

func New(modOptions ...ModOption) (Logger, error) {
	var err error

	l := new(logger)
	l.currRoundTime = time.Now()
	l.option = defaultOption

	for _, fn := range modOptions {
		fn(&l.option)
	}

	if err := validate(l.option); err != nil {
		return nil, err
	}
	if l.option.Filename != "" {
		l.dir = filepath.Dir(l.option.Filename)
		if err = os.MkdirAll(l.dir, 0777); err != nil {
			return nil, err
		}
		if l.fp, err = os.OpenFile(l.option.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
			return nil, err
		}
	}
	if l.option.IsToStdout {
		l.console = os.Stdout
	}

	return l, nil
}

func validate(option Option) error {
	if option.Level < LevelDebug || option.Level > LevelPanic {
		return ErrLog
	}
	return nil
}
