package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/emits-io/emits/configuration"
)

// Index struct
type Index struct {
	File []string `json:"file"`
}

func parseServe() (err error) {

	helpFlag := flag.Bool("h", false, "")
	flagSet := flag.NewFlagSet("serve", flag.ExitOnError)
	taskFlag := flagSet.String("task", "", "")
	groupFlag := flagSet.String("group", "", "")
	portFlag := flagSet.Int("port", 10892, "")
	flagSet.Usage = func() {
		usageServe()
	}
	flagSet.BoolVar(helpFlag, "h", false, "")
	flagSet.BoolVar(helpFlag, "help", false, "")
	flagSet.Parse(os.Args[2:])

	taskName := strings.ToLower(strings.Replace(*taskFlag, " ", "", -1))
	groupName := strings.ToLower(strings.Replace(*groupFlag, " ", "", -1))
	if len(taskName) == 0 && len(groupName) == 0 {
		usageServe()
		return fmt.Errorf(color("task or group argument is required\n", Red, false))
	}

	config, err := configuration.Open()
	if err != nil {
		return err
	}

	var tasks []string

	if len(groupName) > 0 {
		if !config.HasGroup(configuration.Group{Name: groupName}) {
			return fmt.Errorf(fmt.Sprintf("%s %s", color(groupName, Red, false), "is not a valid group"))
		}
		for _, t := range config.GetGroup(configuration.Group{Name: groupName}).Tasks {
			task, err := serve(config, t)
			if err == nil {
				tasks = append(tasks, fmt.Sprintf("%s%s%s%s%s", emits, string(os.PathSeparator), task, string(os.PathSeparator), emits+".json"))
			}
		}
	} else if len(taskName) > 0 {
		task, err := serve(config, taskName)
		if err != nil {
			return err
		}
		tasks = append(tasks, fmt.Sprintf("%s%s%s%s%s", emits, string(os.PathSeparator), task, string(os.PathSeparator), emits+".json"))
	}

	http.Handle("/", IndexHandler(Index{File: tasks}))
	fmt.Println("")
	fmt.Println(fmt.Sprintf("%s:%v", color("http://localhost", Cyan, true), color(fmt.Sprintf("%v", *portFlag), Cyan, true)))
	fmt.Println("")
	printExit("", true)
	http.ListenAndServe(fmt.Sprintf(":%v", *portFlag), nil)
	return nil
}

func serve(config configuration.File, name string) (taskName string, err error) {
	if !config.HasTask(configuration.Task{Name: name}) {
		return "", fmt.Errorf(fmt.Sprintf("%s %s", color(name, Red, false), "is not a valid task"))
	}
	task := config.GetTask(configuration.Task{Name: name})
	http.Handle(fmt.Sprintf("%s%s%s%s%s", string(os.PathSeparator), emits, string(os.PathSeparator), task.Name, string(os.PathSeparator)), AllowHandler())
	return name, nil
}

func printExit(line string, prefixSpace bool) {
	if prefixSpace {
		fmt.Println("")
	} else {
		fmt.Print("\r\033[1A\033[0K")
	}
	if len(line) > 0 {
		fmt.Print("\r\033[1A\033[0K")
		fmt.Println(line)
		fmt.Println("")
	}
	fmt.Println("Exit this utility to stop the server...")
}

// AllowHandler is a restrictive http handler that only serves .json files from a relative emits directory.
func AllowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		var p = fmt.Sprintf(".%s", filepath.Clean(r.URL.String()))
		if jsonExists(p) {
			http.ServeFile(w, r, p)
			return
		}
		http.Error(w, "", 500)
		return
	})
}

// IndexHandler func
func IndexHandler(index Index) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		if r.URL.Path == "/" {
			data, err := json.MarshalIndent(index, "", "\t")
			if err != nil {
				fmt.Fprintf(w, "index handler error")
				return
			}
			fmt.Fprintf(w, string(data))
			return
		}
		http.Error(w, "", 500)
	})
}

func jsonExists(name string) bool {
	if strings.HasSuffix(name, ".json") {
		file, err := os.Stat(name)
		if os.IsNotExist(err) {
			return false
		}
		return !file.IsDir()
	}
	return false
}

func usageServe() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(color("emits serve", Cyan, true), color("[arguments]", Magenta, true))
	fmt.Println("")
	fmt.Println("The arguments are:")
	fmt.Println("")
	fmt.Println(argument("task", "name of the configuration task", Magenta))
	fmt.Println(argument("group", "name of the configuration group", Magenta))
	fmt.Println(argument("port", "port to serve the task files", Magenta))
	fmt.Println("")
}
