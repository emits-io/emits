package command

import (
	"fmt"
	"runtime"

	"github.com/emits-io/emits/colorize"
	"github.com/emits-io/emits/version"
)

func parseVersion() (err error) {
	fmt.Println(fmt.Sprintf("emits version %v.%v.%v %s/%s", version.Major, version.Minor, version.Build, runtime.GOOS, runtime.GOARCH))
	return nil
}

func usageVersion() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(colorize.Printc("emits version", colorize.Cyan, true))
	fmt.Println("")
	fmt.Println("Prints the CLI version, operating system and architecture.")
	fmt.Println("")
}
