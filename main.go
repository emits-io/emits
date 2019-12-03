package main

import (
	"fmt"
	"os"

	"github.com/emits-io/emits/command"
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
