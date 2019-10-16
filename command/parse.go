package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/emits-io/emits/configuration"
)

type uniqueOperator func(value string, config configuration.File) bool

// Color 8-bittype
type Color string

const (
	// emits constant
	emits = "emits"
	//Black Color
	Black Color = "\u001b[30"
	//Red Color
	Red Color = "\u001b[31"
	//Green Color
	Green Color = "\u001b[32"
	//Yellow Color
	Yellow Color = "\u001b[33"
	//Blue Color
	Blue Color = "\u001b[34"
	//Magenta Color
	Magenta Color = "\u001b[35"
	//Cyan Color
	Cyan Color = "\u001b[36"
	//White Color
	White Color = "\u001b[37"
)

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
	return fmt.Errorf("emits %s: unknown command", command)
}

func uniqueFlag(rd *bufio.Reader, name string, message string, required bool, value string, config configuration.File, unique uniqueOperator) (result string) {
	if !unique(value, config) {
		indicator := color("optinal", Yellow, true)
		if required {
			indicator = color("required", Red, true)
		}
		return uniqueFlag(rd, name, message, required, readFlag(rd, fmt.Sprintf("\r\033[1A\033[0K%s %s: `\u001b[33m%s\x1b[0m` \u001b[33m%s\x1b[0m\n%s", name, indicator, value, message, name), required), config, unique)
	}
	return value
}

func readFlag(rd *bufio.Reader, name string, required bool) string {
	indicator := color("optional", Yellow, true)
	if required {
		indicator = color("required", Red, true)
	}
	fmt.Printf("%s %s: ", name, indicator)
	input, _ := rd.ReadString('\n')
	input = strings.TrimSpace(input)
	if required && len(input) == 0 {
		return readFlag(rd, name, required)
	}
	return input
}

func argument(name string, description string, c Color) string {
	return fmt.Sprintf("%s%s %s %s", color("--", c, false), color(name, c, true), color("•", White, false), description)
}

func command(name string, description string, c Color) string {
	return fmt.Sprintf("%s %s %s", color(name, c, true), color("•", White, false), description)
}

func underline(text string) string {
	return fmt.Sprintf("\u001b[4m%s\u001b[0m", text)
}

func bold(text string) string {
	return fmt.Sprintf("\u001b[1m%s\u001b[0m", text)
}

func color(text string, color Color, bright bool) string {
	brightColor := ""
	if bright {
		brightColor = ";1"
	}
	return fmt.Sprintf("%s%sm%s\u001b[0m", color, brightColor, text)
}
