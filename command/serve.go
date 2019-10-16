package command

import (
	"flag"
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/emits-io/emits/configuration"
)

func parseServe() (err error) {

	helpFlag := flag.Bool("h", false, "")
	flagSet := flag.NewFlagSet("serve", flag.ExitOnError)
	taskFlag := flagSet.String("task", "", "")
	portFlag := flagSet.Int("port", 10892, "")
	flagSet.Usage = func() {
		usageServe()
	}
	flagSet.BoolVar(helpFlag, "h", false, "")
	flagSet.BoolVar(helpFlag, "help", false, "")
	flagSet.Parse(os.Args[2:])

	name := strings.ToLower(strings.Replace(*taskFlag, " ", "", -1))
	if len(name) == 0 {
		usageServe()
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
	fileServe := http.Dir(emits + "/" + task.Name)
	fs := http.FileServer(fileServe)
	http.Handle("/", AllowHandler(http.StripPrefix("/"+emits+"/"+task.Name+"/", fs), task.Name, *portFlag))
	fmt.Println("")
	fmt.Println(fmt.Sprintf("%s:%v", color("http://localhost", Cyan, true), color(fmt.Sprintf("%v", *portFlag), Cyan, true)))
	fmt.Println("")
	printExit("", true)
	http.ListenAndServe(fmt.Sprintf(":%v", *portFlag), nil)
	return nil
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

func usageServe() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println(color("emits serve", Cyan, true), color("[arguments]", Magenta, true))
	fmt.Println("")
	fmt.Println("The arguments are:")
	fmt.Println("")
	fmt.Println(argument("task", "name of the configuration task", Magenta))
	fmt.Println(argument("port", "port to serve the task files", Magenta))
	fmt.Println("")
}

// AllowHandler is a restrictive http handler that only serves .json files from a relative emits directory.
func AllowHandler(next http.Handler, task string, port int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		if r.URL.Path != "/" {
			if strings.HasSuffix(r.URL.String(), ".json") {
				fp := filepath.Join("./", filepath.Clean(r.URL.Path))
				_, err := os.Stat(fp)
				if err != nil {
					if os.IsNotExist(err) {
						http.NotFound(w, r)
						printExit(fmt.Sprintf("[%s] http://%s%s", time.Now().Format(time.StampMicro), r.Host, r.URL.String()), false)
						fmt.Fprintf(w, "{\"error\":\"file not found\",\"file\":%q}", html.EscapeString(r.URL.Path))
						return
					}
				}
				printExit(fmt.Sprintf("[%s] http://%s%s", color(time.Now().Format(time.StampMicro), White, true), r.Host, r.URL.String()), false)
				next.ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
				printExit(fmt.Sprintf("[%s] http://%s%s", time.Now().Format(time.StampMicro), r.Host, r.URL.String()), false)
				fmt.Fprintf(w, "{\"error\":\"invalid file format\",\"file\":%q}", html.EscapeString(r.URL.Path))
			}
		} else {
			printExit(fmt.Sprintf("[%s] http://%s%s", color(time.Now().Format(time.StampMicro), White, true), r.Host, r.URL.String()), false)
			http.Redirect(w, r, "/"+emits+"/"+task+"/emits.json", http.StatusFound)
		}
	})
}
