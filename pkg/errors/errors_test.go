package errors

import (
	"errors"
	"fmt"
	"testing"
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
