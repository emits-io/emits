package command

import (
	"fmt"
	"time"

	"github.com/emits-io/emits/colorize"
)

func parseHelp(command string) (err error) {
	switch command {
	case "help":
		usageHelp()
		return nil
	case "init":
		usageInit()
		return nil
	case "list":
		usageList()
		return nil
	case "run":
		usageRun()
		return nil
	case "serve":
		usageServe()
		return nil
	case "delete":
		usageDelete()
		return nil
	case "update":
		usageUpdate()
		return nil
	case "version":
		usageVersion()
		return nil
	}
	return fmt.Errorf(fmt.Sprintf("[%s] Runtime Error '%v' is not a known help command", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc(command, colorize.Red, false)))
}

func usageHelp() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(colorize.Printc("emits", colorize.Cyan, false), colorize.Printc("<command>", colorize.Cyan, true), colorize.Printc("[arguments]", colorize.Magenta, false))
	fmt.Println("")
	fmt.Println("The commands are:")
	fmt.Println("")
	fmt.Println(command("run", "emit files for a configuration task", colorize.Cyan))
	fmt.Println(command("init", "initialize a configuration task", colorize.Cyan))
	fmt.Println(command("list", "list all configuration tasks", colorize.Cyan))
	fmt.Println(command("serve", "serve files for a configuration task", colorize.Cyan))
	fmt.Println(command("update", "update configration task fields", colorize.Cyan))
	fmt.Println(command("delete", "delete configuration task", colorize.Cyan))
	fmt.Println(command("version", "print command line interface version", colorize.Cyan))
	fmt.Println("")
	fmt.Println("Use", colorize.Printc("emits help", colorize.Cyan, false), colorize.Printc("<command>", colorize.Cyan, true), "for more information")
	fmt.Println("")
}
