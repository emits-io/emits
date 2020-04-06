package command

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/emits-io/emits/colorize"
	"github.com/emits-io/emits/configuration"
)

func parseDelete() (err error) {
	helpFlag := flag.Bool("h", false, "")
	flagSet := flag.NewFlagSet("delete", flag.ExitOnError)
	taskFlag := flagSet.String("task", "", "")
	flagSet.Usage = func() {
		usageDelete()
	}
	flagSet.BoolVar(helpFlag, "h", false, "")
	flagSet.BoolVar(helpFlag, "help", false, "")
	flagSet.Parse(os.Args[2:])

	config, err := configuration.Open(true)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("[%s] Runtime Error '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc(err.Error(), colorize.Red, false)))
	}

	name := strings.ToLower(strings.Replace(*taskFlag, " ", "", -1))
	if len(name) == 0 {
		return fmt.Errorf(fmt.Sprintf("[%s] Runtime Error '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc("a task argument is required", colorize.Red, false)))
	}

	task := configuration.Task{Name: name}

	if !config.HasTask(task) {
		return fmt.Errorf(fmt.Sprintf("[%s] Runtime Error '%v' is not a valid task", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc(name, colorize.Red, false)))
	}

	ok := config.DeleteTask(task)
	if ok {
		err := config.Write()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf(fmt.Sprintf("[%s] Runtime Error '%v' could not be removed", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc(name, colorize.Red, false)))
	}
	fmt.Println(fmt.Sprintf("`%v` task configuration deleted", task.Name))
	return nil
}

func usageDelete() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(colorize.Printc("emits delete", colorize.Cyan, true), colorize.Printc("[argument]", colorize.Magenta, true))
	fmt.Println("")
	fmt.Println("The argument is:")
	fmt.Println("")
	fmt.Println(argument("task", "name of the configuration task", colorize.Magenta))
	fmt.Println("")
}
