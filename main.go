/*
TODO for beta release:
x 00. Remove implicit seprator; now requires explicit .. line.
x 01. Make Inline comments optional; do not parse if empty.
x 02. Make Block Line comment optional; do not parse if block struct is empty.
x 03. Allow for single line block comments?
x 04. Strict formatting rules on keyword indent override.
x 05. Delete emits folder before writing files.
x 06. Optional emits config include/exclude filters.
x 07. Prettify command line; create helpers
x 08. Update task command; command to update all config options.
x 09. Delete task command
x 10. Add `flags`.
11. Tests

*/

package main

import (
	"fmt"
	"os"

	"bitbucket.org/emits-io/emits/command"
)

func main() {
	arg := "help"
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	if err := command.Parse(arg); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
