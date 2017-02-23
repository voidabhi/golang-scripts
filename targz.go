package main

import (
	"archive/tar"
	"compress/gzip"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	inputPath  = "input/targz.go"
	outputPath = "output/targz.go.tar.gz"
)

func main() {
	var file *os.File
	var err error
	var writer *gzip.Writer
	var body []byte

	if file, err = os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if writer, err = gzip.NewWriterLevel(file, gzip.BestCompression); err != nil {
		log.Fatalln(err)
	}
	defer writer.Close()

	tw := tar.NewWriter(writer)
	defer tw.Close()

	if body, err = ioutil.ReadFile(inputPath); err != nil {
		log.Fatalln(err)
	}

	if body != nil {
		hdr := &tar.Header{
			Name: path.Base(inputPath),
			Mode: int64(0644),
			Size: int64(len(body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			println(err)
		}
		if _, err := tw.Write(body); err != nil {
			println(err)
		}
	}
}
