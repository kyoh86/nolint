package a

import (
	"fmt"
)

func main() {
	// Some bad format/argTypes
	fmt.Printf("%d")      // want "Printf format %d reads arg #1, but call has 0 args"
	fmt.Printf("%d")      // nolint
	fmt.Printf("%d", 3.6) // want "Printf format %d has arg 3.6 of wrong type float64"
}
