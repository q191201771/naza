#### v0.4.1

- package nazastring:
    - 新增 func SliceByteToStringTmp 和 func StringToSliceByteTmp，用于无拷贝的做string和[]byte的转换

#### v0.4.0

- package log:
    - rename -> package nazalog
    - mkdir 0777 and create file 0666, append if file exist
    - 配置使用 Option
    - 配置默认值修改： 打印至控制台开关默认打开，打印源码文件行号开关默认打开
- package connection:
    - erase func Printf
    - 配置使用 Option
- package nazajson:
    - 新增包，作为系统包 json 的补充
- 其它:
    - repo name nezha -> naza

#### v0.3.0

- package connection:
    - 可配置使用 channel 进行异步发送：Config 中增加 WChanSize。增加 Flush, Done, ModWriteChanSize 三个方法
- package log:
    - 增加 panic 相关的方法
- 其它:
    - test.sh 中添加 gofmt 检查

#### v0.2.0

- package log:
    - 去除了对标准库中log的依赖
    - 日志支持按天翻转
    - 增加 ShortFileFlag 可配置是否打印源码文件及行号的信息
    - 添加一个fatal日志级别，打印完后exit程序
    - 当同时打印至控制台和文件时，打印至文件中的level字段也带颜色属性
    - 增加 Out 接口
    - 日志不再支持按固定大小翻转 [不兼容]
    - 日志级别从0 -> 1开始 [不兼容]

#### v0.1.0

- 删除 /pkg/errors [不兼容]
- package log:
    - 增加 FatalIfErrorNotNil 接口函数，打印错误并退出程序
    - 日志内容中的级别字段右对齐
    - 日志内容中的源码文件名和行号放在整行日志的末尾
    - 增加一些 benchmark
- package assert: 打印正确的源码文件名和行号信息 [bugfix]
- package bele: 增加一些 benchmark
- package unique:
    - 不同的 key 使用不同的自增计数
    - 增加一些 benchmark
- package mockserver: 模拟一些服务端，用于快速测试其它代码
- package mockwriter: 模拟 Writer 接口，用于快速测试其它代码
- 删除 /demo/connstat
- test.sh 脚本只测试 /pkg 目录下的源码

#### v0.0.1

第一个版本
