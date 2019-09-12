
**v0.1.0**

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

**v0.0.1**

第一个版本
