// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package defertaskthread_test

import (
	"testing"
	"time"

	"github.com/q191201771/naza/pkg/nazalog"

	"github.com/q191201771/naza/pkg/defertaskthread"
)

func TestDeferTaskThread(t *testing.T) {
	d := defertaskthread.NewDeferTaskThread()
	for i := 0; i < 300; i += 50 {
		d.Go(i, func(param ...interface{}) {
			ii := param[0].(int)
			nazalog.Debugf("running %d", ii)
		}, i)
	}
	time.Sleep(300 * time.Millisecond)
}
