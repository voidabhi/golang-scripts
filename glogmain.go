package main

import (
	"github.com/golang/glog"
	"os"
	"flag"
	"fmt"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n", )
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	// NOTE: This next line is key you have to call flag.Parse() for the command line 
	// options or "flags" that are defined in the glog module to be picked up.
	flag.Parse()
}

func main() {
	number_of_lines := 100000
	for i := 0; i < number_of_lines; i++ {
		glog.V(2).Infof("LINE: %d", i)
		message := fmt.Sprintf("TEST LINE: %d", i)
		glog.Error(message)
	}
	glog.Flush()
}
