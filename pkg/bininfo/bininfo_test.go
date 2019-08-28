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
