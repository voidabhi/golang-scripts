
package main

import (
	"fmt"
	"io/ioutil"
)

// File to be read
const SAMPLE_FILE = "sample.txt"

func main(){
	// Content in byte[] from the file
	content,err:=ioutil.ReadFile(SAMPLE_FILE)
	
	if err!=nil	{
		fmt.Printf("Error Occurred")
	}
	
	// Converting byte[] to string for printing
	fmt.Printf(string(content))
}

