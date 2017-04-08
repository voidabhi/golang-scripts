package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var host = flag.String("host", "", "host dns or ip address to gather metrics from")
var interval = flag.Int("interval", 60, "interval in seconds between collections")
var cpuprofile = flag.Bool("cpuprofile", false, "enable cpu profiling")
var memprofile = flag.Bool("memprofile", false, "enable memory profiling")
var blockprofile = flag.Bool("blockprofile", false, "enable goroutine block profiling")
var traceprofile = flag.String("traceprofile", "", "enable trace profiling for N seconds")

const (
	MemPath   = "/debug/pprof/heap"
	CPUPath   = "/debug/pprof/profile"
	BlockPath = "/debug/pprof/block"
	TracePath = "/debug/pprof/trace"
	OutputDIR = "./pprof"
)

func main() {
	returnCode := entryPoint(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	os.Exit(returnCode)
}

func entryPoint(cliArgs []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	flag.Parse()

	if *host == "" {
		fmt.Fprintln(stderr, "No host set")
		return 1
	}

	if dirExists, err := exists(OutputDIR); err != nil {
		fmt.Fprintln(stderr, "Directory not found:", err)
		return 1
	} else if !dirExists {
		if err := os.Mkdir(OutputDIR, 0755); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}

	collectionEndpoints := []*url.URL{}
	if !*cpuprofile && !*memprofile && !*blockprofile && *traceprofile == "" {
		for _, path := range []string{MemPath, BlockPath, CPUPath} {
			if url, err := url.Parse(*host + path); err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			} else {
				collectionEndpoints = append(collectionEndpoints, url)
			}
		}
	}

	if *memprofile {
		if url, err := url.Parse(*host + MemPath); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		} else {
			collectionEndpoints = append(collectionEndpoints, url)
		}
	}
	if *blockprofile {
		if url, err := url.Parse(*host + BlockPath); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		} else {
			collectionEndpoints = append(collectionEndpoints, url)
		}
	}
	if *cpuprofile {
		if url, err := url.Parse(*host + CPUPath); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		} else {
			collectionEndpoints = append(collectionEndpoints, url)
		}
	}
	if *traceprofile != "" {
		if url, err := url.Parse(*host + TracePath + "?seconds=" + *traceprofile); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		} else {
			collectionEndpoints = append(collectionEndpoints, url)
		}
	}

	var mu sync.Mutex
	fmt.Fprintln(stdout, "Beginning collection...")
	for {
		fmt.Fprintln(stdout, <-time.After(time.Second*time.Duration(*interval)))
		go func() {
			mu.Lock()
			defer mu.Unlock()
			for _, url := range collectionEndpoints {
				tokens := strings.Split(url.String(), "/")
				fmt.Fprintf(stdout, "\t%s\n", tokens[len(tokens)-1])
				fileName := fmt.Sprintf("pprof.%s.%s.%s.pb.gz", url.Host, tokens[len(tokens)-1], time.Now().Format(time.RFC3339))
				output, err := os.Create(OutputDIR + "/" + fileName)
				if err != nil {
					fmt.Fprintln(stderr, "Error while creating", fileName, "-", err)
					return
				}
				defer output.Close()

				response, err := http.Get(url.String())
				if err != nil {
					fmt.Fprintln(stderr, "Error while downloading", url.String(), "-", err)
					return
				}
				defer response.Body.Close()

				_, err = io.Copy(output, response.Body)
				if err != nil {
					fmt.Fprintln(stderr, "Error while downloading", url.String(), "-", err)
					return
				}
			}
		}()
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
