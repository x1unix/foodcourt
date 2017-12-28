package cmd

import (
	"fmt"
	"errors"
	"../../shared/logger"
	"os"
)

// Application version (defined at build args)
var version string = "1.0.0"

var commands = make(map[string] Handler)
var commandDescriptions = make(map[string] string)


const defaultCmd = ""

type Handler func() (bool, error)

func HandleFunc(name string, description string, handler Handler) {
	commands[name] = handler
	commandDescriptions[name] = description
}

func Default(handler Handler) {
	HandleFunc(defaultCmd, "Prints this message", handler)
}

func Bootstrap() {

	// Register default task
	Default(func() (bool, error) {
		lines := []string{
			fmt.Sprintf("Command worker for FoodCourt, version %s\n", version),
			"usage: fc-worker [taskname] args...\n",
			"  taskname : string tells what task to execute \n",
			"List of available tasks:",
		}

		for _, val := range lines {
			fmt.Println(val)
		}

		for taskName := range commands {
			if taskName == defaultCmd {
				continue
			}

			fmt.Println(fmt.Sprintf("  %s - %s", taskName, commandDescriptions[taskName]))
		}

		return false, errors.New("no valid task specified")
	})
}

func Run() {
	taskName := defaultCmd

	if len(os.Args) > 1 {
		taskName = os.Args[1]
	}

	Call(taskName)
}

func Call(taskName string) {
	originalCmd := taskName
	log := logger.GetLogger()

	if commands[taskName] == nil {
		taskName = defaultCmd
	}

	_, err := commands[taskName]();

	if err != nil {
		log.Error(fmt.Sprintf("Command '%s' ran with error: %s", originalCmd, err.Error()))
		return
	}

	log.Info(fmt.Sprintf("Command '%s' ran successfuly", taskName))
}