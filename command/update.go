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

func parseUpdate() (err error) {
	helpFlag := flag.Bool("h", false, "")
	flagSet := flag.NewFlagSet("update", flag.ExitOnError)
	taskFlag := flagSet.String("task", "", "")
	noPromptFlag := flagSet.Bool("no-prompt", false, "")
	//
	nameFlag := flagSet.String("name", "", "")
	descriptionFlag := flagSet.String("description", "", "")
	fileIncludeFlag := flagSet.String("file-include", "", "")
	fileExcludeFlag := flagSet.String("file-exclude", "", "")
	keywordIncludeFlag := flagSet.String("keyword-include", "", "")
	keywordExcludeFlag := flagSet.String("keyword-exclude", "", "")
	configurationIncludeFlag := flagSet.String("configuration-include", "", "")
	configurationExcludeFlag := flagSet.String("configuration-exclude", "", "")
	commentBlockOpenFlag := flagSet.String("comment-block-open", "", "")
	commentBlockLineFlag := flagSet.String("comment-block-line", "", "")
	commentBlockCloseFlag := flagSet.String("comment-block-close", "", "")
	commentInlineFlag := flagSet.String("comment-inline", "", "")
	sourceFlag := flagSet.String("source", "", "")
	//
	flagSet.Usage = func() {
		usageUpdate()
	}
	flagSet.BoolVar(helpFlag, "h", false, "")
	flagSet.BoolVar(helpFlag, "help", false, "")
	flagSet.Parse(os.Args[2:])

	name := strings.ToLower(strings.Replace(*taskFlag, " ", "", -1))
	if len(name) == 0 {
		usageUpdate()
		return fmt.Errorf(color("task argument is required\n", Red, false))
	}

	config, err := configuration.Open()
	if err != nil {
		return err
	}

	task := configuration.Task{Name: name}

	if !config.HasTask(task) {
		return fmt.Errorf(fmt.Sprintf("%s %s", color(task.Name, Red, false), "is not a valid task"))
	}

	reader := bufio.NewReader(os.Stdin)

	task = config.GetTask(task)

	taskName := strings.ToLower(strings.Replace(*nameFlag, " ", "", -1))
	if len(taskName) > 0 {
		taskName = uniqueFlag(reader, "task name", "is already taken", true, taskName, config, isUniqueTaskName)
		task.Name = taskName
	}

	description := strings.ToLower(strings.TrimSpace(*descriptionFlag))
	if len(description) > 0 {
		task.Description = description
	}

	fileInclude := strings.ToLower(strings.TrimSpace(*fileIncludeFlag))
	if len(fileInclude) > 0 {
		task.File.Include = strings.Split(strings.ToLower(fileInclude), " ")
	}

	fileExclude := strings.ToLower(strings.TrimSpace(*fileExcludeFlag))
	if len(fileExclude) > 0 {
		task.File.Exclude = strings.Split(strings.ToLower(fileExclude), " ")
	}

	keywordInclude := strings.ToLower(strings.TrimSpace(*keywordIncludeFlag))
	if len(keywordInclude) > 0 {
		task.Keyword.Include = strings.Split(strings.ToLower(keywordInclude), " ")
	}

	keywordExclude := strings.ToLower(strings.TrimSpace(*keywordExcludeFlag))
	if len(keywordExclude) > 0 {
		task.Keyword.Exclude = strings.Split(strings.ToLower(keywordExclude), " ")
	}

	configurationInclude := strings.ToLower(strings.TrimSpace(*configurationIncludeFlag))
	if len(configurationInclude) > 0 {
		task.Configuration.Include = strings.Split(strings.ToLower(configurationInclude), " ")
	}

	configurationExclude := strings.ToLower(strings.TrimSpace(*configurationExcludeFlag))
	if len(configurationExclude) > 0 {
		task.Configuration.Exclude = strings.Split(strings.ToLower(configurationExclude), " ")
	}

	commentBlockOpen := strings.ToLower(strings.TrimSpace(*commentBlockOpenFlag))
	if len(commentBlockOpen) > 0 {
		task.Comment.Block.Open = commentBlockOpen
	}
	commentBlockLine := strings.ToLower(strings.TrimSpace(*commentBlockLineFlag))
	if len(commentBlockLine) > 0 {
		task.Comment.Block.Line = commentBlockLine
	}
	commentBlockClose := strings.ToLower(strings.TrimSpace(*commentBlockCloseFlag))
	if len(commentBlockClose) > 0 {
		task.Comment.Block.Close = commentBlockClose
	}
	commentInline := strings.ToLower(strings.TrimSpace(*commentInlineFlag))
	if len(commentInline) > 0 {
		task.Comment.Inline = commentInline
	}

	source := strings.ToLower(strings.TrimSpace(*sourceFlag))
	if len(source) > 0 && source == "true" || len(source) > 0 && source == "false" {
		task.Source = source == "true"
	}

	task = task.Sanitize()

	if *noPromptFlag == false {
		fmt.Println("Type", color("yes", Red, true), "below to add the following task to the configuration file:")
	} else {
		fmt.Println("The following task has been added to the configuration file:")
	}
	fmt.Println("")

	preview, _ := json.MarshalIndent(task, "", "\t")
	fmt.Println(fmt.Sprintf("%v", string(preview)))
	fmt.Println("")
	update := false
	if *noPromptFlag == false {
		input := strings.TrimSpace(readFlag(reader, "Is this OK?", true))
		fmt.Println("")
		if strings.ToLower(input) == "yes" {
			update = true
		}
	} else {
		update = true
	}
	if update {
		err := config.UpdateTask(task)
		if err != nil {
			return err
		}
		fmt.Println("Use", color(fmt.Sprintf("emits run %s", task.Name), Cyan, true), "to emit this task")
	} else {
		fmt.Println("Use", color("emits help", Cyan, true), "for more information")
	}
	fmt.Println("")

	return nil
}

func usageUpdate() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(color("emits update", Cyan, true), color("[arguments]", Magenta, true))
	fmt.Println("")
	fmt.Println("The arguments are:")
	fmt.Println("")
	fmt.Println(argument("task", "name of the configuration task", Magenta))
	fmt.Println("")
	fmt.Println(argument("name", "new name of the configuration task", Magenta))
	fmt.Println(argument("description", "new description", Magenta))
	fmt.Println(argument("file-include", "file include patterns", Magenta))
	fmt.Println(argument("file-exclude", "file exclude patterns", Magenta))
	fmt.Println(argument("keyword-include", "keyword includes", Magenta))
	fmt.Println(argument("keyword-exclude", "keyword excludes", Magenta))
	fmt.Println(argument("configuration-include", "configuration includes", Magenta))
	fmt.Println(argument("configuration-exclude", "configuration excludes", Magenta))
	fmt.Println(argument("comment-block-open", "comment block open", Magenta))
	fmt.Println(argument("comment-block-line", "comment block line", Magenta))
	fmt.Println(argument("comment-block-close", "comment block close", Magenta))
	fmt.Println(argument("comment-inline", "comment inline", Magenta))
	fmt.Println(argument("source", "allow source", Magenta))
	fmt.Println("")
}
