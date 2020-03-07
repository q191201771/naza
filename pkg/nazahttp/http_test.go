package nazahttp_test

import (
	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazahttp"
	"testing"
)

func TestDownloadHttpFile(t *testing.T) {
	n, err := nazahttp.DownloadHttpFile("http://pengrl.com", "/tmp/index.html", 10000)
	assert.Equal(t, true, n > 0)
	assert.Equal(t, nil, err)

	n, err = nazahttp.DownloadHttpFile("http://127.0.0.1:12356", "/tmp/index.html", 10000)
	assert.IsNotNil(t, err)

	n, err = nazahttp.DownloadHttpFile("http://pengrl.com", "/notexist/index.html", 10000)
	assert.IsNotNil(t, err)
}
