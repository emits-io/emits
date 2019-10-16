package command

import (
	"fmt"
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
	return fmt.Errorf("emits help %s: unknown command", command)
}

func usageHelp() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(color("emits", Cyan, false), color("<command>", Cyan, true), color("[arguments]", Magenta, false))
	fmt.Println("")
	fmt.Println("The commands are:")
	fmt.Println("")
	fmt.Println(command("run", "emit files for a configuration task", Cyan))
	fmt.Println(command("init", "initialize a configuration task", Cyan))
	fmt.Println(command("list", "list all configuration tasks", Cyan))
	fmt.Println(command("serve", "serve files for a configuration task", Cyan))
	fmt.Println(command("update", "update configration task fields", Cyan))
	fmt.Println(command("delete", "delete configuration task", Cyan))
	fmt.Println(command("version", "print command line interface version", Cyan))
	fmt.Println("")
	fmt.Println("Use", color("emits help", Cyan, false), color("<command>", Cyan, true), "for more information")
	fmt.Println("")
}
