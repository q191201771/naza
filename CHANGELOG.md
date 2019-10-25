#### v0.5.0

- package filebatch:
    - 遍历读取文件发生错误时，不退出遍历，而是将错误在回调中返回给上层
- package connection:
    - bugfix，初始化 write chan 相关的信息是通过 write chan 的配置，而不是 write buf 的配置
- package slicebytepool:
    - 新增的包，一个 []byte 内存池
- package nazamd5:
    - 新增的包，md5 操作
- package consistenthash:
    - 新增的包，一致性hash
- package bufferpool:
    - 删除 bufferpool 包
- demo/myapp:
    - 用于演示 package bininfo 的使用
- `demo/add_blog_license`:
    - 修改 license 内容
- 其他：
    - 统一error变量的命名方式及内容格式，涉及到的 package：filebatch, connection, taskpool, nazalog

#### v0.4.3

- package bufferpool
    - 新增的包，bytes.Buffer 池
- package nazaatomic
    - 新增的包，对 sync.atomic 的再封装
- package taskpool
    - 新增的包，协程池
- test.sh
    - 做更多的 go tools 检查
- `demo/add_blog_license`
    - 修改许可证样式
    - 检查许可证是否存在时，只检查声明两个字
- `demo/add_go_license`
    - 用户名和邮箱由命令行参数传入

#### v0.4.2

- package filebatch:
    - 新增的包，用于文件批处理操作
- 新增 demo/add_go_license：给 Go 仓库的所有go源码文件添加MIT许可证
- 新增 demo/add_blog_license：给我自己博客的所有文章尾部添加声明

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
