<p align="center">
<br>
Go语言基础库
<br><br>
<a title="Release" target="_blank" href="https://github.com/q191201771/naza/tags"><img src="https://img.shields.io/github/tag/q191201771/naza.svg?label=release"></a>
<a title="CI" target="_blank" href="https://github.com/q191201771/naza/actions/workflows/ci.yml"><img src="https://github.com/q191201771/naza/actions/workflows/ci.yml/badge.svg"></a>
<a title="codecov" target="_blank" href="https://codecov.io/gh/q191201771/naza"><img src="https://codecov.io/gh/q191201771/naza/branch/master/graph/badge.svg?token=9wWcoiktZl"></a>
<a title="goreportcard" target="_blank" href="https://goreportcard.com/report/github.com/q191201771/naza"><img src="https://goreportcard.com/badge/github.com/q191201771/naza"></a>
<br>
<a title="codeline" target="_blank" href="https://github.com/q191201771/naza"><img src="https://sloc.xyz/github/q191201771/naza/?category=code"></a>
<a title="license" target="_blank" href="https://github.com/q191201771/naza/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
<a title="lastcommit" target="_blank" href="https://github.com/q191201771/naza/commits/master"><img src="https://img.shields.io/github/commit-activity/m/q191201771/naza.svg?style=flat-square"></a>
<a title="commitactivity" target="_blank" href="https://github.com/q191201771/naza/graphs/commit-activity"><img src="https://img.shields.io/github/last-commit/q191201771/naza.svg?style=flat-square"></a>
<br>
<a title="language" target="_blank" href="https://github.com/q191201771/naza"><img src="https://img.shields.io/github/languages/count/q191201771/naza.svg?style=flat-square"></a>
<a title="toplanguage" target="_blank" href="https://github.com/q191201771/naza"><img src="https://img.shields.io/github/languages/top/q191201771/naza.svg?style=flat-square"></a>
<a title="godoc" target="_blank" href="https://godoc.org/github.com/q191201771/naza"><img src="http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square"></a>
<br><br>
</p>

---

#### 工程目录说明

```
pkg/                     ...... 源码包
    |-- defertaskthread  ...... 执行延时任务
    |-- connection/      ...... 对net.Conn接口的二次封装
    |-- taskpool/        ...... 非阻塞协程池，协程数量可动态增长，可配置最大协程并发数量，可手动释放空闲的协程
    |-- nazajson/        ...... json操作
    |-- nazalog/         ...... 日志库

    |-- assert/          ...... 提供了单元测试时的断言功能，减少一些模板代码
    |-- bele/            ...... 大小端转换操作
    |-- bininfo/         ...... 将编译时源码的git版本信息（当前commit log的sha值和commit message），编译时间，Go版本，平台打入程序中
    |-- circularqueue    ...... 底层基于切片实现的固定容量大小的FIFO的环形队列
    |-- dataops/         ...... 数据处理
    |-- fake/            ...... 实现一些常用的接口，hook一些不方便测试的代码
    |-- filebatch/       ...... 文件批处理操作
    |-- filesystemlayer/ ...... 对文件操作的封装，可以使用内存作为磁盘使用
    |-- mock/            ...... 模拟一些标准库中的常用接口，方便测试
    |-- nazaatomic/      ...... 原子操作
    |-- nazabits/        ...... 位操作
    |-- nazabytes/       ...... 字节切片，内存块操作
    |-- nazacolor/       ...... 控制台打印颜色相关
    |-- nazaerrors/      ...... error相关
    |-- nazahttp/        ...... http操作
    |-- nazamd5/         ...... md5操作
    |-- nazanet/         ...... socket操作相关
    |-- nazareflect/     ...... 利用反射做的一些操作
    |-- nazastring/      ...... string和[]byte相关的操作
    |-- unique/          ...... 对象唯一ID
    |-- nazasync/        ...... 对sync的封装，比如定位sync.Mutex死锁
    |-- chartbar/        ...... ascii柱状图
    |-- bitrate/         ...... 计算带宽
    |-- ratelimit/       ...... 限流器，令牌桶，漏桶
    |-- lru/             ...... LRU缓存
    |-- consistenthash/  ...... 一致性哈希
    |-- crypto/          ...... 加解密操作
    |-- slicebytepool/   ...... []byte内存池
    |-- snowflake/       ...... 分布式唯一性ID生成器
playground/              ...... Go实验代码片段
demo/                    ...... 示例相关的代码
```

#### 依赖

无任何第三方依赖

#### 联系我

欢迎扫码加我微信，进行技术交流或扯淡。

<img src="https://pengrl.com/images/yoko_vx.jpeg" width="180" height="180" />

#### 项目名 naza 由来

本仓库主要用于存放我自己写的一些 Go 基础库代码。目前主要服务于我的另一个项目： [lal](https:////github.com/q191201771/lal)

naza 即哪吒（正确拼音为 nezha，我女儿发音读作 naza，少一个字母，挺好~），希望本仓库以后能像三头六臂，有多种武器的哪吒一样，为我提供一个趁手的工具箱。

