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

	"github.com/emits-io/emits/colorize"
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
		return fmt.Errorf(fmt.Sprintf("[%s] Runtime Error '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc("task or group argument is required", colorize.Red, false)))
	}

	config, err := configuration.Open()
	if err != nil {
		return err
	}

	if len(groupName) > 0 {
		if !config.HasGroup(configuration.Group{Name: groupName}) {
			return fmt.Errorf(fmt.Sprintf("[%s] Runtime Error '%v' is not a valid group", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc(groupName, colorize.Red, false)))
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
		return fmt.Errorf(fmt.Sprintf("[%s] Runtime Error '%v' is not a valid task", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc(name, colorize.Red, false)))
	}
	task := config.GetTask(configuration.Task{Name: name})
	start := time.Now()
	fmt.Println(fmt.Sprintf("[%s] Starting Task '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Yellow, false), colorize.Printc(task.Name, colorize.Yellow, false)))
	cache := data.Cache{}
	for _, grammar := range task.Grammar {
		g, err := data.CacheGrammarFile(grammar)
		if err != nil {
			fmt.Println(fmt.Sprintf("[%s] Grammar Error '%v' could not be loaded", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc(grammar, colorize.Red, false)))
		} else {
			fmt.Println(fmt.Sprintf("[%s] Using Grammar '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Yellow, false), colorize.Printc(grammar, colorize.Yellow, false)))
			cache.GrammarFile = append(cache.GrammarFile, g)
		}
	}
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
	for _, file := range matches {
		filePath := filepath.Join(file)
		err := data.Write(filePath, task, cache, output)
		if err != nil {
			fmt.Println(fmt.Sprintf("[%s] Writing Error '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc("./"+filePath, colorize.Red, false)))
		} else {
			filename := filepath.Join(output, filePath+".json")
			index.Files = append(index.Files, filename)
			fmt.Println(fmt.Sprintf("[%s] Emitting File '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Yellow, false), colorize.Printc("./"+filename, colorize.Yellow, false)))
		}
	}
	file, err := json.MarshalIndent(index, "", "\t")
	if err != nil {
		fmt.Println(fmt.Sprintf("[%s] Writing Error '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc("./"+indexFilePath, colorize.Red, false)))
	} else {
		err = ioutil.WriteFile(indexFilePath, file, 0644)
		if err != nil {
			fmt.Println(fmt.Sprintf("[%s] Writing Error '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc("./"+indexFilePath, colorize.Red, false)))
		} else {
			t := time.Now()
			elapsed := t.Sub(start)
			fmt.Println(fmt.Sprintf("[%s] Finished Task '%v' after %v", colorize.Printc(time.Now().Format(time.Stamp), colorize.Yellow, false), colorize.Printc(task.Name, colorize.Yellow, false), elapsed))
		}
	}
	return nil
}

func usageRun() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(colorize.Printc("emits run", colorize.Cyan, true), colorize.Printc("[arguments]", colorize.Magenta, true))
	fmt.Println("")
	fmt.Println("The arguments are:")
	fmt.Println("")
	fmt.Println(argument("task", "name of the configuration task", colorize.Magenta))
	fmt.Println(argument("group", "name of the configuration group", colorize.Magenta))
	fmt.Println(argument("output", "output path to emit data", colorize.Magenta))
	fmt.Println("")
}
