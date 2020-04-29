// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package fake_test

import (
	"testing"
	"time"

	"github.com/q191201771/naza/pkg/assert"

	"github.com/q191201771/naza/pkg/fake"
)

func TestWithFakeTimeNow(t *testing.T) {
	fake.WithFakeTimeNow(func() time.Time {
		return time.Now().Add(time.Duration(2 * time.Hour))
	}, func() {
		n := fake.Time_Now()
		assert.Equal(t, true, n.Sub(time.Now()).Hours() > 1)
	})
}
