package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

type Worker struct {
	Command string
	Args    string
	Output  chan string
}

func (cmd *Worker) Run() {
	out, err := exec.Command(cmd.Command, cmd.Args).Output()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Output <- string(out)
}

func Collect(c chan string) {
	for {
		msg := <-c
		fmt.Printf("The command result is %s\n", msg)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var read string
	fmt.Println("When you're ready press ENTER to spawn goroutine")
	fmt.Scanln(&read)
	c := make(chan string)

	phpService := &Worker{Command: "php", Args: "slowService.php", Output: c}
	pythonService := &Worker{Command: "python", Args: "mediumService.py", Output: c}

	go phpService.Run()
	fmt.Println("Doing something...")
	go pythonService.Run()
	go Collect(c)

	fmt.Scanln(&read)
}
