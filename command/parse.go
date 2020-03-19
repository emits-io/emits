package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/emits-io/emits/colorize"
	"github.com/emits-io/emits/configuration"
)

const (
	// emits constant
	emits = "emits"
)

type uniqueOperator func(value string, config configuration.File) bool

var isUniqueTaskName uniqueOperator = func(value string, config configuration.File) bool {
	return !config.HasTask(configuration.Task{Name: value})
}

// Parse a command
func Parse(command string) (err error) {
	switch command {
	case "help":
		subcommand := "help"
		if len(os.Args) > 2 {
			subcommand = os.Args[2]
		}
		return parseHelp(subcommand)
	case "init":
		return parseInit()
	case "list":
		return parseList()
	case "run":
		return parseRun()
	case "serve":
		return parseServe()
	case "delete":
		return parseDelete()
	case "update":
		return parseUpdate()
	case "version":
		return parseVersion()
	}
	return fmt.Errorf(fmt.Sprintf("[%s] Runtime Error '%v' is not a known command", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc(command, colorize.Red, false)))
}

func uniqueFlag(rd *bufio.Reader, name string, message string, required bool, value string, config configuration.File, unique uniqueOperator) (result string) {
	if !unique(value, config) {
		indicator := colorize.Printc("optinal", colorize.Yellow, true)
		if required {
			indicator = colorize.Printc("required", colorize.Red, true)
		}
		return uniqueFlag(rd, name, message, required, readFlag(rd, fmt.Sprintf("\r\033[1A\033[0K%s %s: `\u001b[33m%s\x1b[0m` \u001b[33m%s\x1b[0m\n%s", name, indicator, value, message, name), required), config, unique)
	}
	return value
}

func readFlag(rd *bufio.Reader, name string, required bool) string {
	indicator := colorize.Printc("optional", colorize.Yellow, true)
	if required {
		indicator = colorize.Printc("required", colorize.Red, true)
	}
	fmt.Printf("%s %s: ", name, indicator)
	input, _ := rd.ReadString('\n')
	input = strings.TrimSpace(input)
	if required && len(input) == 0 {
		return readFlag(rd, name, required)
	}
	return input
}

func argument(name string, description string, c colorize.Color) string {
	return fmt.Sprintf("%s%s %s %s", colorize.Printc("--", c, false), colorize.Printc(name, c, true), colorize.Printc("•", colorize.White, false), description)
}

func command(name string, description string, c colorize.Color) string {
	return fmt.Sprintf("%s %s %s", colorize.Printc(name, c, true), colorize.Printc("•", colorize.White, false), description)
}
