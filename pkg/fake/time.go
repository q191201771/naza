// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package fake

import "time"

var now = time.Now

func Time_Now() time.Time {
	return now()
}

func WithFakeTimeNow(n func() time.Time, fn func()) {
	now = n
	fn()
	now = time.Now
}
