// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// package nazalog 日志库
package nazalog

import "errors"

// 这是一个以使用方便为主要目标的日志库，特性：
//
// * 带日志级别
// * 可选输出至控制台或文件，也可以同时输出
// * 日志文件支持按天翻转
// * 支持是否输出源码文件及行号
// * 业务日志起始位置固定，方便查看
// * 支持Assert，断言失败后的行为可配置
// * 支持全局日志对象，独立日志对象，日志对象都可以配置，相互间可以赋值
// * 支持设置前缀，并且前缀可叠加，使得可以按repo ，package，对象等维度添加不同的前缀
// * 支持标准库中的打印接口函数（但是没有适配非打印接口），方便替换标准库日志
// * 日志文件目录不存在则自动创建
//
// 目前性能和标准库log相当

var ErrLog = errors.New("naza.log:fxxk")

type Logger interface {
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{}) // 打印日志并退出程序
	Panicf(format string, v ...interface{})

	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Panic(v ...interface{})

	Out(level Level, calldepth int, s string)

	// 断言失败后的行为由配置项Option.AssertBehavior决定
	// 注意，expected和actual的类型必须相同，比如int(1)和int32(1)是不相等的
	Assert(expected interface{}, actual interface{})

	// flush to disk, typically
	Sync()

	// 添加前缀，新生成一个Logger对象，如果老Logger也有prefix，则老Logger依然打印老prefix，新Logger打印多个prefix
	WithPrefix(s string) Logger

	// 下面这些打印接口是为兼容标准库，让某些已使用标准库日志的代码替换到nazalog方便一些
	Output(calldepth int, s string) error
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatalln(v ...interface{})
	Panicln(v ...interface{})

	// 获取配置项，注意，作用是只读，非修改配置
	GetOption() Option
}

type Option struct {
	Level Level `json:"level"` // 日志级别，大于等于该级别的日志才会被输出

	// 文件输出和控制台输出可同时打开
	// 控制台输出主要用做开发时调试，打开后level字段使用彩色输出
	Filename   string `json:"filename"`     // 输出日志文件名，如果为空，则不写日志文件。可包含路径，路径不存在时，将自动创建
	IsToStdout bool   `json:"is_to_stdout"` // 是否以stdout输出到控制台

	IsRotateDaily bool `json:"is_rotate_daily"` // 日志按天翻转

	ShortFileFlag       bool `json:"short_file_flag"`        // 是否在每行日志尾部添加源码文件及行号的信息
	TimestampFlag       bool `json:"timestamp_flag"`         // 是否在每行日志首部添加时间戳的信息
	TimestampWithMSFlag bool `json:"timestamp_with_ms_flag"` // 时间戳是否精确到毫秒
	LevelFlag           bool `json:"level_flag"`             // 日志是否包含日志级别字段

	AssertBehavior AssertBehavior `json:"assert_behavior"` // 断言失败时的行为
}

// 没有配置的属性，将按如下配置
var defaultOption = Option{
	Level:               LevelDebug,
	Filename:            "",
	IsToStdout:          true,
	IsRotateDaily:       false,
	ShortFileFlag:       true,
	TimestampFlag:       true,
	TimestampWithMSFlag: true,
	LevelFlag:           true,
	AssertBehavior:      AssertError,
}

type Level uint8

const (
	LevelTrace Level = iota // 0
	LevelDebug              // 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
	LevelLogNothing
)

func (l Level) ReadableString() string {
	switch l {
	case LevelTrace:
		return "LevelTrace"
	case LevelDebug:
		return "LevelDebug"
	case LevelInfo:
		return "LevelInfo"
	case LevelWarn:
		return "LevelWarn"
	case LevelError:
		return "LevelError"
	case LevelFatal:
		return "LevelFatal"
	case LevelPanic:
		return "LevelPanic"
	case LevelLogNothing:
		return "LevelLogNothing"
	default:
		return "unknown"
	}
}

type AssertBehavior uint8

const (
	_           AssertBehavior = iota
	AssertError                // 1
	AssertFatal
	AssertPanic
)

func (a AssertBehavior) ReadableString() string {
	switch a {
	case AssertError:
		return "AssertError"
	case AssertFatal:
		return "AssertFatal"
	case AssertPanic:
		return "AssertPanic"
	default:
		return "unknown"
	}
}

type ModOption func(option *Option)

func New(modOptions ...ModOption) (Logger, error) {
	return newLogger(modOptions...)
}
