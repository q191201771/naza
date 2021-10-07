package color

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	fmt.Printf("\033[22;31;42mHello\033[0m\n")

	fmt.Printf("%s\n", Wrap("Hello", FormatNonBold, FgRed, BgGreen))

	fmt.Printf("%s\n", WrapWithFgColor("Hello", FgRed))

	fmt.Printf("%s\n", WrapBlack("Hello"))
	fmt.Printf("%s\n", WrapRed("Hello"))
	fmt.Printf("%s\n", WrapGreen("Hello"))
	fmt.Printf("%s\n", WrapYellow("Hello"))
	fmt.Printf("%s\n", WrapBlue("Hello"))
	fmt.Printf("%s\n", WrapCyan("Hello"))
	fmt.Printf("%s\n", WrapWhite("Hello"))

	fmt.Printf("%s%s%s\n", SimplePrefixRed, "Hello", SimpleSuffix)
}
