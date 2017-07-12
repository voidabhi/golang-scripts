package main

import "fmt"

type Command interface {
	Execute() string
}

type PingCommand struct{}
func (p *PingCommand) Execute() string {
	return "react pings"
}

type StatusCommand struct{}
func (p *StatusCommand) Execute() string {
	return "status command"
}

func execByName(name string) string {
	// Register commands
	commands := map[string]Command{
		"ping":   &PingCommand{},
		"status": &StatusCommand{},
	}

	if command := commands[name]; command == nil {
		return "No such command found, throw error?"
	} else {
		return command.Execute()
	}
}

func main() {

	fmt.Println(execByName("status"))
	fmt.Println(execByName("ping"))
	
	fmt.Println(execByName("unkown"))
}
