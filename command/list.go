package command

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/emits-io/emits/colorize"
	"github.com/emits-io/emits/configuration"
)

func parseList() (err error) {

	helpFlag := flag.Bool("h", false, "")
	flagSet := flag.NewFlagSet("list", flag.ExitOnError)
	descriptionFlag := flagSet.Bool("description", false, "")
	flagSet.Usage = func() {
		usageList()
	}
	flagSet.BoolVar(helpFlag, "h", false, "")
	flagSet.BoolVar(helpFlag, "help", false, "")
	flagSet.Parse(os.Args[2:])

	config, err := configuration.Open(true)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("[%s] Runtime Error '%v'", colorize.Printc(time.Now().Format(time.Stamp), colorize.Red, false), colorize.Printc(err.Error(), colorize.Red, false)))
	}

	fmt.Println("")
	fmt.Println("The following tasks have been configured:")
	fmt.Println("")

	for _, task := range config.Tasks {
		description := ""
		if *descriptionFlag {
			description = fmt.Sprintf(" • %s", task.Description)
		}
		fmt.Println(fmt.Sprintf("%s%s", task.Name, description))
	}
	fmt.Println("")
	return nil
}

func usageList() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	//fmt.Println(Colorize("emits list", Cyan, true), Colorize("[flag]", Green, true))
	fmt.Println("")
	fmt.Println("The flag is:")
	fmt.Println("")
	fmt.Println(argument("description", "output description with task name", colorize.Green))
	fmt.Println("")
}
