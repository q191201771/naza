package errors

import (
	"fmt"
	"testing"
	"errors"
)

func TestPanicIfErrorOccur(t *testing.T) {
	PanicIfErrorOccur(nil)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	PanicIfErrorOccur(errors.New("fxxk."))
}
