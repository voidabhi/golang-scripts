package main

import (
	"compress/gzip"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

// That was easy!  Let's add another few features.  If -z is passed, we want any
// DestFile's to be gzipped.  If -md5 is passed, we want print the md5sum of the
// data that's been transfered instead of the data itself.
var Config struct {
	Silent   bool
	DestFile string
	Gzip     bool
	Md5      bool
}

func init() {
	flag.StringVar(&Config.DestFile, "o", "", "output file")
	flag.BoolVar(&Config.Silent, "s", false, "silent (do not output to stdout)")
	flag.BoolVar(&Config.Gzip, "z", false, "gzip file output")
	flag.BoolVar(&Config.Md5, "md5", false, "stdout md5sum instead of body")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Usage: go run 03-curl.go [options] <url>")
		os.Exit(-1)
	}
}

func main() {
	url := flag.Args()[0]
	r, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Our Md5 hash destination, which is an io.Writer that computes the
	// hash of whatever is written to it.
	hash := md5.New()
	var writers []io.Writer

	// if we aren't in Silent mode, we've got to output something
	if !Config.Silent {
		// If -md5 was passed, write to the hash instead of os.Stdout
		if Config.Md5 {
			writers = append(writers, hash)
		} else {
			writers = append(writers, os.Stdout)
		}
	}

	// if DestFile was provided, we've got to write a file
	if len(Config.DestFile) > 0 {
		// by declaring writer here as a WriteCloser, we're saying that we don't care
		// what the underlying implementation will be, all we require is something that
		// can Write and Close;  both os.File and the gzip.Writer are WriteClosers.
		var writer io.WriteCloser
		writer, err := os.Create(Config.DestFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		// If we're in Gzip mode, wrap the writer in gzip
		if Config.Gzip {
			writer = gzip.NewWriter(writer)
		}
		writers = append(writers, writer)
		defer writer.Close()
	}

	// MultiWriter(io.Writer...) returns a single writer which multiplexes its
	// writes across all of the writers we pass in.
	dest := io.MultiWriter(writers...)

	// write to dest the same way as before, copying from the Body
	io.Copy(dest, r.Body)
	if err = r.Body.Close(); err != nil {
		fmt.Println(err)
		return
	}

	// finally, if we were in Md5 output mode, lets output the checksum and url:
	if Config.Md5 {
		fmt.Printf("%x  %s\n", hash.Sum(nil), url)
	}
}
