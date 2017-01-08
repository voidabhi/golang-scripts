package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func CmdPiper() {
	fmt.Println("Cmd Piper")
	// set up the pipe
	pr, pw := io.Pipe()
	defer pw.Close()

	// tell the command to write to our pipe
	cmd := exec.Command("cat", "fruit.txt")
	cmd.Stdout = pw

	go func() {
		defer pr.Close()
		// copy the data written to the pipereader from the command to stdout
		if _, err := io.Copy(os.Stdout, pr); err != nil {
			log.Fatal(err)
		}
	}()

	// run the command, which will write all output to the pipewriter
	// which will end up in the pipereader
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
