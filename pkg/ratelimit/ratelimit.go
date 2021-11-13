// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package ratelimit

// LeakBucket和TokenBucket的区别
//
// LeakBucket:  业务方从LeakBucket获取资源的时间间隔必须>=设置的时间间隔
// TokenBucket: 内部会持续生成资源并进行缓存，外部可以一次性获取
//

type RateLimiter interface {
	TryAquire() error
	WaitUntilAquire()
}
