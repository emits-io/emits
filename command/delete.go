package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/emits-io/emits/configuration"
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

	name := strings.ToLower(strings.Replace(*taskFlag, " ", "", -1))
	if len(name) == 0 {
		return fmt.Errorf("\x1b[31;1m%s\x1b[0m", "a task argument is required")
	}

	config, err := configuration.Open()
	if err != nil {
		return err
	}

	task := configuration.Task{Name: name}

	if !config.HasTask(task) {
		return fmt.Errorf("`\x1b[31;1m%s\x1b[0m` \x1b[31;1mis not a valid task\x1b[0m", name)
	}

	ok := config.DeleteTask(task)
	if ok {
		err := config.Write()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("task could not be removed")
	}
	fmt.Println(fmt.Sprintf("`%v` task configuration deleted", task.Name))
	return nil
}

func usageDelete() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(color("emits delete", Cyan, true), color("[argument]", Magenta, true))
	fmt.Println("")
	fmt.Println("The argument is:")
	fmt.Println("")
	fmt.Println(argument("task", "name of the configuration task", Magenta))
	fmt.Println("")
}
