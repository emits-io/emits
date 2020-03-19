package main

import (
	"fmt"
	"os"
	"time"

	"github.com/emits-io/emits/colorize"
	"github.com/emits-io/emits/command"
)

func main() {
	if os.Geteuid() == 0 {
		fmt.Println(fmt.Sprintf("[%s] Runtime Error '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc("do not run emits with root privilages", colorize.Red, false)))
		os.Exit(1)
	}
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
