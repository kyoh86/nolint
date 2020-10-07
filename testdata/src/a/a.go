package a

import (
	"fmt"
)

func main() {
	// Some bad format/argTypes
	fmt.Printf("%d")      // want "Printf format %d reads arg #1, but call has 0 args"
	fmt.Printf("%d")      // nolint
	fmt.Printf("%d", 3.6) // want "Printf format %d has arg 3.6 of wrong type float64"

	// With another comments
	fmt.Printf("%d") // nolint // foobar
	fmt.Printf("%d") // foobar // nolint

	// With nolint directive above the block.
	//nolint
	fmt.Printf("%d")

	// This statement has
	// a lot of comments
	// but in the end
	// we mark it with the directive
	//nolint
	fmt.Printf("%d")

	// This statement has
	// a lot of comments
	// and in the end
	// we MISS mark it with
	// the directive
	fmt.Printf("%d") // want "Printf format %d reads arg #1, but call has 0 args"

	// Multiline support
	fmt.Printf( // want "Printf format %d reads arg #1, but call has 0 args"
		"%d",
	)

	fmt.Printf( //nolint // multiline on diagnositcs line
		"%d",
	)

	//nolint // multiline on top line
	fmt.Printf(
		"%d",
	)

	// Directive not directly next to the node is not valid
	//nolint

	fmt.Printf( // want "Printf format %d reads arg #1, but call has 0 args"
		"%d",
	)
}
