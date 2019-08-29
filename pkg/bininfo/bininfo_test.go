package bininfo

import (
	"fmt"
	"testing"
)

func TestStringifySingleLine(t *testing.T) {
	fmt.Println(StringifySingleLine())
}

func TestStringifyMultiLine(t *testing.T) {
	fmt.Println(StringifyMultiLine())
}

func TestCorner(t *testing.T) {
	GitStatus = ""
	beauty()
}