package command

import (
	"fmt"
	"runtime"

	"bitbucket.org/emits-io/emits/version"
)

func parseVersion() (err error) {
	fmt.Println(fmt.Sprintf("emits version %v.%v.%v %s/%s", version.Major, version.Minor, version.Build, runtime.GOOS, runtime.GOARCH))
	return nil
}

func usageVersion() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(color("emits version", Cyan, true))
	fmt.Println("")
	fmt.Println("Prints the CLI version, operating system and architecture.")
	fmt.Println("")
}
