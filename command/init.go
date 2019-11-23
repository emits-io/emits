package command

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/emits-io/emits/configuration"
)

func parseInit() (err error) {
	config, err := configuration.Open()
	if err != nil {
		return err
	}
	helpFlag := flag.Bool("h", false, "")
	flagSet := flag.NewFlagSet("init", flag.ExitOnError)
	nameFlag := flagSet.String("name", "", "")
	descriptionFlag := flagSet.String("description", "", "")
	includeFlag := flagSet.String("include", "", "")
	excludeFlag := flagSet.String("exclude", "", "")
	commentOpenFlag := flagSet.String("open", "", "")
	commentLineFlag := flagSet.String("line", "", "")
	commentCloseFlag := flagSet.String("close", "", "")
	commentInlineFlag := flagSet.String("inline", "", "")
	sourceFlag := flagSet.Bool("source", false, "")
	noPromptFlag := flagSet.Bool("no-prompt", false, "")
	flagSet.Usage = func() {
		usageInit()
	}
	flagSet.BoolVar(helpFlag, "h", false, "")
	flagSet.BoolVar(helpFlag, "help", false, "")
	flagSet.Parse(os.Args[2:])
	reader := bufio.NewReader(os.Stdin)
	if *noPromptFlag == false {
		fmt.Println("")
		fmt.Println(color("This utility will walk you through the most common fields to configure a task", White, true))
		fmt.Println("")
		fmt.Println("Use", color("emits help update", Cyan, true), "documentation on all available fields")
		fmt.Println("")
	}
	name := strings.ToLower(strings.Replace(*nameFlag, " ", "", -1))
	if len(name) == 0 {
		name = strings.ToLower(strings.Replace(readFlag(reader, "task name", true), " ", "", -1))
	}
	name = uniqueFlag(reader, "task name", "is already taken", true, name, config, isUniqueTaskName)
	description := strings.ToLower(strings.TrimSpace(*descriptionFlag))
	if len(description) == 0 && *noPromptFlag == false {
		description = strings.ToLower(strings.TrimSpace(readFlag(reader, "task description", false)))
	}
	var includePatterns []string
	include := strings.TrimSpace(*includeFlag)
	if len(include) == 0 {
		include = strings.TrimSpace(readFlag(reader, "file include patterns", true))
	}
	includePatterns = strings.Split(strings.ToLower(include), " ")
	var excludePatterns []string
	exclude := strings.TrimSpace(*excludeFlag)
	if len(exclude) == 0 && *noPromptFlag == false {
		exclude = strings.TrimSpace(readFlag(reader, "file exclude patterns", false))
	}
	excludePatterns = strings.Split(strings.ToLower(exclude), " ")
	commentOpen := strings.TrimSpace(*commentOpenFlag)
	commentLine := strings.TrimSpace(*commentLineFlag)
	commentClose := strings.TrimSpace(*commentCloseFlag)
	commentInline := strings.TrimSpace(*commentInlineFlag)
	if len(commentOpen) == 0 && *noPromptFlag == false || len(commentOpen) == 0 && len(commentClose) > 0 || len(commentOpen) == 0 && len(commentLine) > 0 {
		commentOpen = readFlag(reader, "block comment open", len(commentLine) > 0 || len(commentClose) > 0)
	}
	if len(commentLine) == 0 && *noPromptFlag == false {
		commentLine = readFlag(reader, "block comment line", false)
	}
	if len(commentClose) == 0 && *noPromptFlag == false || len(commentClose) == 0 && len(commentOpen) > 0 || len(commentClose) == 0 && len(commentLine) > 0 {
		commentClose = readFlag(reader, "block comment close", len(commentOpen) > 0 || len(commentLine) > 0)
	}
	if len(commentInline) == 0 && *noPromptFlag == false {
		commentInline = readFlag(reader, "inline comment", false)
	}
	task := configuration.Task{
		Name:        name,
		Description: description,
		Source:      *sourceFlag,
		Comment: configuration.Comment{
			Block: configuration.Block{
				Open:  commentOpen,
				Line:  commentLine,
				Close: commentClose,
			},
			Inline: commentInline,
		},
		File: configuration.Pattern{
			Include: includePatterns,
			Exclude: excludePatterns,
		},
	}

	task = task.Sanitize()

	fmt.Println("")
	if *noPromptFlag == false {
		fmt.Println("The following task will be added to the configuration file:")
	} else {
		fmt.Println("The following task has been added to the configuration file:")
	}
	fmt.Println("")

	preview, _ := json.MarshalIndent(task, "", "\t")
	fmt.Println(fmt.Sprintf("%v", string(preview)))
	fmt.Println("")
	if *noPromptFlag == false {
		fmt.Println("Type", color("yes", Red, true), "below to add this task...")
		fmt.Println("")
	}
	create := false
	if *noPromptFlag == false {
		input := strings.TrimSpace(readFlag(reader, "Is this OK?", true))
		fmt.Println("")
		if strings.ToLower(input) == "yes" {
			create = true
		}
	} else {
		create = true
	}
	if create {
		err = config.CreateTask(task)
		if err != nil {
			return err
		}
		err = config.Write()
		if err != nil {
			return err
		}
		fmt.Println("Use", color(fmt.Sprintf("emits run %s%s", color("--task ", Magenta, true), color(task.Name, Magenta, false)), Cyan, true), "to emit this task")
	} else {
		fmt.Println("Use", color("emits help", Cyan, true), "for more information")
	}
	fmt.Println("")
	return nil
}

func usageInit() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(color("emits init", Cyan, true), color("[arguments]", Magenta, true), color("[flags]", Green, true))
	fmt.Println("")
	fmt.Println("The arguments are:")
	fmt.Println("")
	fmt.Println(argument("name", "name of the configuration task", Magenta))
	fmt.Println(argument("description", "description of the task", Magenta))
	fmt.Println(argument("include", "space delemited match patterns", Magenta))
	fmt.Println(argument("exclude", "space delimited match patterns", Magenta))
	fmt.Println(argument("open", "comment block open characters", Magenta))
	fmt.Println(argument("line", "comment block line characters", Magenta))
	fmt.Println(argument("close", "comment block close characters", Magenta))
	fmt.Println(argument("inline", "comment inline characters", Magenta))
	fmt.Println("")
	fmt.Println("The flags are:")
	fmt.Println("")
	fmt.Println(argument("source", "allow source code to be emitted", Green))
	fmt.Println(argument("no-prompt", "do not promt for confirmation", Green))
	fmt.Println("")
}
