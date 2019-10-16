package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/emits-io/emits/configuration"
	"bitbucket.org/emits-io/emits/data"
)

func parseRun() (err error) {
	helpFlag := flag.Bool("h", false, "")
	flagSet := flag.NewFlagSet("run", flag.ExitOnError)
	taskFlag := flagSet.String("task", "", "")
	outputFlag := flagSet.String("output", "", "")
	flagSet.Usage = func() {
		usageRun()
	}
	flagSet.BoolVar(helpFlag, "h", false, "")
	flagSet.BoolVar(helpFlag, "help", false, "")
	flagSet.Parse(os.Args[2:])

	name := strings.ToLower(strings.Replace(*taskFlag, " ", "", -1))
	if len(name) == 0 {
		usageRun()
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

	task = config.GetTask(task)

	matches, err := task.Files()
	if err != nil {
		return err
	}

	output := *outputFlag
	if len(output) == 0 {
		output = filepath.Join(emits, task.Name)
	}

	err = os.RemoveAll(output)
	if err != nil {

	}

	index := configuration.Index{}
	indexFilePath := filepath.Join(output, "emits.json")

	plural := "s"
	if len(matches) == 1 {
		plural = ""
	}

	fmt.Println(fmt.Sprintf("[\x1b[32;1m%s\x1b[0m] %v of %v file%s processed...", time.Now().Format(time.StampMicro), 0, len(matches), plural))

	for i, file := range matches {
		filePath := filepath.Join(file)
		err := data.Write(filePath, task, output)
		if err != nil {
			//fmt.Println(fmt.Sprintf("[\x1b[31;1m%s\x1b[0m] ✕ %s ➤ %s", time.Now().Format(time.StampMicro), filePath, err))
		} else {
			fmt.Print("\r\033[1A\033[0K")
			fmt.Println(fmt.Sprintf("[\x1b[32;1m%s\x1b[0m] %v of %v file%s processed...", time.Now().Format(time.StampMicro), i+1, len(matches), plural))
			index.Files = append(index.Files, filepath.Join(output, filePath+".json"))
		}
	}
	file, err := json.MarshalIndent(index, "", "\t")
	if err != nil {
		//fmt.Println(fmt.Sprintf("[\x1b[31;1m%s\x1b[0m] ✕ %s", time.Now().Format(time.StampMicro), indexFilePath))
	} else {
		err = ioutil.WriteFile(indexFilePath, file, 0644)
		if err != nil {
			//fmt.Println(fmt.Sprintf("[\x1b[31;1m%s\x1b[0m] ✕ %s", time.Now().Format(time.StampMicro), indexFilePath))
		} else {
			fmt.Println(fmt.Sprintf("[\x1b[32;1m%s\x1b[0m] complete", time.Now().Format(time.StampMicro)))
		}
	}
	return nil
}

func usageRun() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(color("emits run", Cyan, true), color("[arguments]", Magenta, true))
	fmt.Println("")
	fmt.Println("The arguments are:")
	fmt.Println("")
	fmt.Println(argument("task", "name of the configuration task", Magenta))
	fmt.Println(argument("output", "output path to emit data", Magenta))
	fmt.Println("")
}
