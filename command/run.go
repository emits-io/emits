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

	"github.com/emits-io/emits/configuration"
	"github.com/emits-io/emits/data"
)

func parseRun() (err error) {
	helpFlag := flag.Bool("h", false, "")
	flagSet := flag.NewFlagSet("run", flag.ExitOnError)
	taskFlag := flagSet.String("task", "", "")
	groupFlag := flagSet.String("group", "", "")
	outputFlag := flagSet.String("output", "", "")
	flagSet.Usage = func() {
		usageRun()
	}
	flagSet.BoolVar(helpFlag, "h", false, "")
	flagSet.BoolVar(helpFlag, "help", false, "")
	flagSet.Parse(os.Args[2:])

	taskName := strings.ToLower(strings.Replace(*taskFlag, " ", "", -1))
	groupName := strings.ToLower(strings.Replace(*groupFlag, " ", "", -1))

	if len(taskName) == 0 && len(groupName) == 0 {
		usageRun()
		return fmt.Errorf(color("task or group argument is required\n", Red, false))
	}

	config, err := configuration.Open()
	if err != nil {
		return err
	}

	if len(groupName) > 0 {
		if !config.HasGroup(configuration.Group{Name: groupName}) {
			return fmt.Errorf(fmt.Sprintf("%s %s", color(groupName, Red, false), "is not a valid group"))
		}
		for _, t := range config.GetGroup(configuration.Group{Name: groupName}).Tasks {
			run(config, t, outputFlag)
		}
	} else if len(taskName) > 0 {
		err := run(config, taskName, outputFlag)
		if err != nil {
			return err
		}
	}
	return nil
}

func run(config configuration.File, name string, outputFlag *string) (err error) {
	if !config.HasTask(configuration.Task{Name: name}) {
		return fmt.Errorf(fmt.Sprintf("%s %s", color(name, Red, false), "is not a valid task"))
	}
	task := config.GetTask(configuration.Task{Name: name})

	//fmt.Println(fmt.Sprintf("[\x1b[32;1m%s\x1b[0m] %v", time.Now().Format(time.StampMicro), "initializing cache"))

	cache := data.Cache{}
	for _, grammar := range task.Grammar {
		g, err := data.CacheGrammarFile(grammar)
		if err != nil {
			fmt.Println(fmt.Sprintf("[\x1b[31;1m%s\x1b[0m] failed to load %v grammar", time.Now().Format(time.StampMicro), grammar))
		} else {
			fmt.Println(fmt.Sprintf("[\x1b[32;1m%s\x1b[0m] loaded %v grammar", time.Now().Format(time.StampMicro), grammar))
			cache.GrammarFile = append(cache.GrammarFile, g)
		}
	}

	fmt.Println(fmt.Sprintf("[\x1b[32;1m%s\x1b[0m] running %v task", time.Now().Format(time.StampMicro), task.Name))
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
		err := data.Write(filePath, task, cache, output)
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
	fmt.Println(argument("group", "name of the configuration group", Magenta))
	fmt.Println(argument("output", "output path to emit data", Magenta))
	fmt.Println("")
}
